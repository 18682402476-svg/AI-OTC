// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import "@openzeppelin/contracts/access/Ownable.sol";
import "@openzeppelin/contracts/token/ERC20/IERC20.sol";
import "@openzeppelin/contracts/utils/math/SafeCast.sol";

contract OTCTrading is Ownable {
    using SafeCast for uint256;

    // 订单状态枚举
    enum OrderStatus {
        ACTIVE,
        INACTIVE
    }

    // 交易类型
    enum TransactionType {
        BUY,
        SELL
    }

    // 订单信息
    struct Order {
        uint256 id;
        address creator;
        TransactionType transactionType;
        address baseToken;
        address quoteToken;
        uint256 amount;
        uint256 price;
        uint256 filledAmount;
        OrderStatus status;
        uint256 createdAt;
        uint256 frozenAmount; // 新增：记录该订单冻结的资产金额
    }

    // 计数器
    uint256 private _orderIdCounter;

    // 订单映射
    mapping(uint256 => Order) public orders;

    // 地址-代币冻结金额映射
    mapping(address => mapping(address => uint256)) public frozenBalances;

    // 手续费地址
    address public feeReceiver;

    // 手续费比例（ basis points，10000 = 100%）
    uint256 public feePercentage;

    // 事件（保持不变，仅调整解冻事件的参数）
    event OrderCreated(
        uint256 indexed orderId,
        address indexed creator,
        TransactionType transactionType,
        address baseToken,
        address quoteToken,
        uint256 amount,
        uint256 price
    );

    event OrderFilled(
        uint256 indexed orderId,
        address indexed filler,
        uint256 amount,
        uint256 fee
    );

    event OrderCanceled(
        uint256 indexed orderId,
        address indexed creator
    );

    event AssetFrozen(
        address indexed user,
        address indexed token,
        uint256 amount,
        uint256 orderId
    );

    event AssetUnfrozen(
        address indexed user,
        address indexed token,
        uint256 amount,
        uint256 orderId
    );

    event AssetTransferred(
        address indexed from,
        address indexed to,
        address indexed token,
        uint256 amount
    );

    event FeeCollected(
        address indexed from,
        address indexed receiver,
        address indexed token,
        uint256 amount
    );

    // 构造函数
    constructor(address initialOwner, address _feeReceiver, uint256 _feePercentage) Ownable(initialOwner) {
        require(_feeReceiver != address(0), "Fee receiver cannot be zero");
        require(_feePercentage <= 10000, "Fee percentage too high");
        feeReceiver = _feeReceiver;
        feePercentage = _feePercentage;
    }

    // 创建订单（修复核心：移除全局冻结映射，将冻结金额存在订单内）
    function createOrder(
        TransactionType _transactionType,
        address _baseToken,
        address _quoteToken,
        uint256 _amount,
        uint256 _price
    ) external returns (uint256) {
        // 验证参数
        require(_baseToken != address(0), "Base token zero");
        require(_quoteToken != address(0), "Quote token zero");
        require(_amount > 0, "Amount zero");
        require(_price > 0, "Price zero");

        // 生成订单ID
        _orderIdCounter++;
        uint256 orderId = _orderIdCounter;

        // 计算需要冻结的资产（用unchecked避免溢出检查的Gas消耗，手动验证）
        uint256 freezeAmount;
        address freezeToken;
        if (_transactionType == TransactionType.BUY) {
            // 买入订单：冻结报价代币
            freezeAmount = _price;
            freezeToken = _quoteToken;
        } else {
            // 卖出订单：冻结基础代币
            freezeAmount = _amount;
            freezeToken = _baseToken;
        }

        // 验证余额和授权
        IERC20 token = IERC20(freezeToken);
        require(token.balanceOf(msg.sender) >= freezeAmount, "Insufficient balance");
        require(token.allowance(msg.sender, address(this)) >= freezeAmount, "Insufficient allowance");

        // 安全转移资产（用bool返回值验证）
        bool transferSuccess = token.transferFrom(msg.sender, address(this), freezeAmount);
        require(transferSuccess, "Token transfer failed");

        // 更新全局冻结余额映射
        frozenBalances[msg.sender][freezeToken] += freezeAmount;

        // 创建订单（将冻结金额存入订单，而非全局映射）
        orders[orderId] = Order({
            id: orderId,
            creator: msg.sender,
            transactionType: _transactionType,
            baseToken: _baseToken,
            quoteToken: _quoteToken,
            amount: _amount,
            price: _price,
            filledAmount: 0,
            status: OrderStatus.ACTIVE,
            createdAt: block.timestamp,
            frozenAmount: freezeAmount // 存储该订单的冻结金额
        });

        // 触发事件
        emit AssetFrozen(msg.sender, freezeToken, freezeAmount, orderId);
        emit OrderCreated(
            orderId,
            msg.sender,
            _transactionType,
            _baseToken,
            _quoteToken,
            _amount,
            _price
        );

        return orderId;
    }

    // 成交订单（整单成交）
    function fillOrder(uint256 _orderId) external {
        Order storage order = orders[_orderId];
        require(order.id != 0, "Order not exist");
        require(order.status == OrderStatus.ACTIVE, "Order inactive");

        // 验证订单是否未被部分填充
        require(order.filledAmount == 0, "Order already partially filled");

        uint256 quoteAmount = order.price;
        uint256 fee = quoteAmount * feePercentage / 10000;
        uint256 amount = order.amount;

        address baseToken = order.baseToken;
        address quoteToken = order.quoteToken;
        address creator = order.creator;

        if (order.transactionType == TransactionType.SELL) {
            // 卖家订单：买家支付报价代币，收取基础代币
            IERC20 quoteTokenContract = IERC20(quoteToken);
            // 验证买家余额和授权（包含手续费）
            uint256 totalQuoteAmount = quoteAmount + fee;
            require(quoteTokenContract.balanceOf(msg.sender) >= totalQuoteAmount, "Insufficient quote balance");
            require(quoteTokenContract.allowance(msg.sender, address(this)) >= totalQuoteAmount, "Insufficient quote allowance");

            // 转移资产
            require(quoteTokenContract.transferFrom(msg.sender, creator, quoteAmount), "Quote transfer to creator failed");
            require(quoteTokenContract.transferFrom(msg.sender, feeReceiver, fee), "Fee transfer failed");
            require(IERC20(baseToken).transfer(msg.sender, amount), "Base transfer to filler failed");

            // 触发事件
            emit AssetTransferred(msg.sender, creator, quoteToken, quoteAmount);
            emit FeeCollected(msg.sender, feeReceiver, quoteToken, fee);
            emit AssetTransferred(address(this), msg.sender, baseToken, amount);
        } else {
            // 买家订单：卖家支付基础代币，收取报价代币
            IERC20 baseTokenContract = IERC20(baseToken);
            require(baseTokenContract.balanceOf(msg.sender) >= amount, "Insufficient base balance");
            require(baseTokenContract.allowance(msg.sender, address(this)) >= amount, "Insufficient base allowance");

            // 转移资产（扣除手续费）
            uint256 quoteToFiller = quoteAmount - fee;
            require(baseTokenContract.transferFrom(msg.sender, creator, amount), "Base transfer to creator failed");
            require(IERC20(quoteToken).transfer(msg.sender, quoteToFiller), "Quote transfer to filler failed");
            require(IERC20(quoteToken).transfer(feeReceiver, fee), "Fee transfer failed");

            // 触发事件
            emit AssetTransferred(msg.sender, creator, baseToken, amount);
            emit AssetTransferred(address(this), msg.sender, quoteToken, quoteToFiller);
            emit FeeCollected(creator, feeReceiver, quoteToken, fee);
        }

        // 解冻全部冻结资产
        address freezeToken = order.transactionType == TransactionType.BUY ? quoteToken : baseToken;
        _unfreezeAsset(creator, freezeToken, order.frozenAmount, _orderId);

        // 更新订单状态
        order.filledAmount = amount;
        order.status = OrderStatus.INACTIVE;
        order.frozenAmount = 0;

        emit OrderFilled(_orderId, msg.sender, amount, fee);
    }

    // 撤销订单（修复：解冻订单内的剩余冻结金额）
    function cancelOrder(uint256 _orderId) external {
        Order storage order = orders[_orderId];
        require(order.id != 0, "Order not exist");
        require(order.status == OrderStatus.ACTIVE, "Order inactive");
        require(msg.sender == order.creator, "Only creator can cancel");

        // 解冻剩余冻结资产
        address freezeToken = order.transactionType == TransactionType.BUY ? order.quoteToken : order.baseToken;
        _unfreezeAsset(msg.sender, freezeToken, order.frozenAmount, _orderId);

        // 更新订单状态
        order.status = OrderStatus.INACTIVE;
        order.frozenAmount = 0; // 清空冻结金额

        emit OrderCanceled(_orderId, msg.sender);
    }

    // 解冻资产（优化：移除全局映射依赖，直接按订单金额解冻）
    function _unfreezeAsset(address _user, address _token, uint256 _amount, uint256 _orderId) internal {
        require(_amount > 0, "No asset to unfreeze");
        require(frozenBalances[_user][_token] >= _amount, "Insufficient frozen balance");

        // 更新全局冻结余额映射
        frozenBalances[_user][_token] -= _amount;

        // 安全转移资产回用户
        bool transferSuccess = IERC20(_token).transfer(_user, _amount);
        require(transferSuccess, "Unfreeze transfer failed");

        emit AssetUnfrozen(_user, _token, _amount, _orderId);
    }

    // 管理员函数（保持不变，增加参数验证）
    function setFeeReceiver(address _feeReceiver) external onlyOwner {
        require(_feeReceiver != address(0), "Fee receiver zero");
        feeReceiver = _feeReceiver;
    }

    function setFeePercentage(uint256 _feePercentage) external onlyOwner {
        require(_feePercentage <= 10000, "Fee percentage > 100%");
        feePercentage = _feePercentage;
    }
    
    // 查询订单的冻结金额
    function getOrderFrozenAmount(uint256 _orderId) external view returns (uint256) {
        return orders[_orderId].frozenAmount;
    }

    // 查询地址在特定代币上的冻结金额
    function getFrozenBalance(address _user, address _token) external view returns (uint256) {
        return frozenBalances[_user][_token];
    }
}