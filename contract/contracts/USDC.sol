// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import {ERC20} from '@openzeppelin/contracts/token/ERC20/ERC20.sol';
import {Ownable} from '@openzeppelin/contracts/access/Ownable.sol';
import {ReentrancyGuard} from '@openzeppelin/contracts/utils/ReentrancyGuard.sol';
import {IERC20} from '@openzeppelin/contracts/token/ERC20/IERC20.sol';

/**
 * @title USDC
 * @notice 标准ERC20代币
 * @dev 资金池使用独立的FundPool合约管理
 */
contract USDC is ERC20, Ownable, ReentrancyGuard {
    // Token消费事件
    event TokenConsumed(address indexed user, uint256 amount, bytes32 indexed reason);
    
    // 代币交易历史事件
    event TokenTransaction(address indexed user, address indexed from, address indexed to, uint256 amount, string action, uint256 timestamp);
    
    // 资金池合约地址
    address public immutable fundsPool;
    
    // 代币交易结构
    struct TokenTransactionRecord {
        uint256 id;
        address from;
        address to;
        uint256 amount;
        string action;
        string description;
        uint256 timestamp;
    }
    
    // 用户交易历史记录
    mapping(address => TokenTransactionRecord[]) public userTransactions;
    
    // 全局交易ID
    uint256 public nextTransactionId;
    
    /**
     * @dev 构造函数
     * @param name 代币名称
     * @param symbol 代币符号
     * @param initialSupply 初始供应量（18位小数）
     * @param fundPoolAddress 资金池合约地址
     */
    constructor(
        string memory name,
        string memory symbol,
        uint256 initialSupply,
        address fundPoolAddress
    ) ERC20(name, symbol) Ownable(msg.sender) {
        require(fundPoolAddress != address(0), "USDC: Fund pool address cannot be zero");
        // 设置资金池合约地址
        fundsPool = fundPoolAddress;
        // 铸造初始供应量到资金池
        _mint(fundsPool, initialSupply);
        // 初始化交易ID
        nextTransactionId = 1;
    }
    
    /**
     * @dev 记录代币交易（带描述）
     */
    function _recordTransaction(address user, address from, address to, uint256 amount, string memory action, string memory description) internal {
        userTransactions[user].push(TokenTransactionRecord({
            id: nextTransactionId,
            from: from,
            to: to,
            amount: amount,
            action: action,
            description: description,
            timestamp: block.timestamp
        }));
        
        emit TokenTransaction(user, from, to, amount, action, block.timestamp);
        nextTransactionId++;
    }
    
    /**
     * @dev 记录代币交易（兼容旧版本）
     */
    function _recordTransaction(address user, address from, address to, uint256 amount, string memory action) internal {
        _recordTransaction(user, from, to, amount, action, "");
    }
    
    /**
     * @dev 从资金池转移代币到用户（仅允许授权合约调用）
     * @param to 接收奖励的用户地址
     * @param amount 奖励金额
     */
    function reward(address to, uint256 amount) external nonReentrant {
        require(to != address(0), "USDC: Cannot reward to zero address");
        require(amount > 0, "USDC: Reward amount must be greater than 0");
        
        // 从资金池转移代币
        _transfer(fundsPool, to, amount);
        
        // 记录交易
        _recordTransaction(to, fundsPool, to, amount, "REWARD", "System reward");
    }
    
    /**
     * @dev 从资金池转移代币到用户（带描述的版本）
     * @param to 接收奖励的用户地址
     * @param amount 奖励金额
     * @param description 奖励描述
     */
    function reward(address to, uint256 amount, string calldata description) external nonReentrant {
        require(to != address(0), "USDC: Cannot reward to zero address");
        require(amount > 0, "USDC: Reward amount must be greater than 0");
        
        // 从资金池转移代币
        _transfer(fundsPool, to, amount);
        
        // 记录交易
        _recordTransaction(to, fundsPool, to, amount, "REWARD", description);
    }
        
    /**
     * @dev 消费代币（带描述的版本）
     * @param user 用户地址
     * @param amount 消费金额
     * @param reason 消费原因
     * @param description 消费描述
     */
    function consume(address user, uint256 amount, bytes32 reason, string calldata description) external nonReentrant {
        require(user != address(0), "USDC: User address cannot be zero");
        require(amount > 0, "USDC: Consumption amount must be greater than 0");
        require(balanceOf(user) >= amount, "USDC: Insufficient balance");
        
        // 将消费的代币转移到资金池
        _transfer(user, fundsPool, amount);
        
        // 记录消费事件
        emit TokenConsumed(user, amount, reason);
        
        // 记录交易
        _recordTransaction(user, user, fundsPool, amount, "CONSUME", description);
    }
    
    /**
     * @dev 铸造代币（仅允许管理员调用，且只能铸造到资金池）
     * @param amount 铸造金额
     */
    function mintToFundsPool(uint256 amount) external onlyOwner nonReentrant {
        require(amount > 0, "USDC: Mint amount must be greater than 0");
        
        // 铸造代币到资金池
        _mint(fundsPool, amount);
        
        // 记录交易
        _recordTransaction(fundsPool, address(0), fundsPool, amount, "MINT");
    }
    
    /**
     * @dev 从资金池提取资金（仅管理员）
     * @param to 转移目标地址
     * @param amount 转移金额
     */
    function transferFromFundsPool(address to, uint256 amount) external onlyOwner nonReentrant {
        require(to != address(0), "USDC: Cannot transfer to zero address");
        require(amount > 0, "USDC: Transfer amount must be greater than 0");
        
        _transfer(fundsPool, to, amount);
        
        // 记录交易
        _recordTransaction(to, fundsPool, to, amount, "WITHDRAW");
    }
    
    /**
     * @dev 获取用户的交易历史记录
     * @param user 用户地址
     * @param start 起始索引
     * @param end 结束索引
     * @return TokenTransactionRecord[] 用户的交易历史记录
     */
    function getUserTransactions(address user, uint256 start, uint256 end) external view returns (TokenTransactionRecord[] memory) {
        TokenTransactionRecord[] storage transactions = userTransactions[user];
        uint256 length = transactions.length;
        
        if (start >= length) {
            return new TokenTransactionRecord[](0);
        }
        
        uint256 actualEnd = end < length ? end : length - 1;
        uint256 actualLength = actualEnd - start + 1;
        
        TokenTransactionRecord[] memory result = new TokenTransactionRecord[](actualLength);
        for (uint256 i = 0; i < actualLength; i++) {
            result[i] = transactions[start + i];
        }
        
        return result;
    }
    
    /**
     * @dev 获取用户的交易记录总数
     * @param user 用户地址
     * @return uint256 用户的交易记录总数
     */
    function getUserTransactionCount(address user) external view returns (uint256) {
        return userTransactions[user].length;
    }
}
