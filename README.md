# AgentOTC 项目说明

## 项目概述

AgentOTC 是一个基于区块链的OTC交易平台，支持代理注册、订单管理、资金管理等功能。本项目包含智能合约部分和后端服务部分，提供了完整的OTC交易解决方案。

## 项目结构

```
AgentOTC/
├── contract/           # 智能合约相关文件
│   ├── contracts/      # 合约源码
│   ├── scripts/        # 部署和测试脚本
│   └── test/           # 合约测试
├── mcp_server/         # MCP服务
├── planfrom/           # 后端应用
│   ├── src/main/java/com/agentotc/  # Java源码
│   │   ├── config/     # 配置类
│   │   ├── contract/   # 合约交互类
│   │   ├── controller/ # 控制器
│   │   ├── mapper/     # 数据访问层
│   │   ├── model/      # 数据模型
│   │   ├── service/    # 服务层
│   │   └── utils/      # 工具类
│   ├── src/main/resources/  # 资源文件
│   │   ├── application.yml  # 应用配置
│   │   └── init.sql    # 数据库初始化脚本
│   │   └── logback-spring.xml  # 日志配置文件
│   └── src/test/java/com/agentotc/  # 测试代码
│   └── API文档.md                # API文档
└── README.md           # 项目说明文档
```

## 合约部署

### 环境准备

1. 安装 Node.js 和 npm
2. 安装 Hardhat：`npm install --save-dev hardhat`
3. 安装依赖：在 `contract` 目录下运行 `npm install`

### 部署步骤

1. 配置网络：编辑 `hardhat.config.js` 文件，设置网络参数
2. 部署合约：
   - 部署 MON 代币：`npx hardhat run scripts/deployMON.js --network <network-name>`
   - 部署 USDC 代币：`npx hardhat run scripts/mintUSDC.js --network <network-name>`
   - 部署 OTC 交易合约：`npx hardhat run scripts/deployOTCTrading.js --network <network-name>`

3. 部署完成后，合约地址会保存在 `contract-addresses.json` 文件中

## 平台运行部署

### 环境准备

1. 安装 JDK 1.8 或更高版本
2. 安装 Maven
3. 安装 MySQL
4. 安装 Redis

### 配置修改

1. 编辑 `planfrom/src/main/resources/application.yml` 文件，修改以下配置：
   - 数据库连接信息
   - Redis 连接信息
   - 区块链网络配置
   - 合约地址配置

### 构建和运行

1. 在 `planfrom` 目录下执行 Maven 构建：`mvn clean package`
2. 运行应用：`java -jar target/AgentOTC-1.0-SNAPSHOT.jar`

## 数据库初始化

1. 启动 MySQL 服务
2. 创建数据库：`CREATE DATABASE agentotc DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;`
3. 执行初始化脚本：`mysql -u root -p agentotc < planfrom/src/main/resources/init.sql`

## 对接注册 Agent 接口

### 接口信息

- **接口地址**：`/agent/register`
- **请求方式**：POST
- **请求参数**：
  ```json
  {
    "agent_name": "代理名称", // 或 agentName
    "wallet_address": "代理钱包地址", // 或 profitAddress
    "callback_url": "回调URL" // 可选，或 callbackUrl
  }
  ```

- **响应示例**：
  ```json
  {
    "success": true,
    "data": {
      "agentId": 1,
      "agentName": "代理名称",
      "walletAddress": "代理钱包地址",
      "token": "访问令牌",
      "status": "状态",
      "profitAddress": "收益地址",
      "createdAt": "2026-02-27T16:00:00"
    }
  }
  ```

### 调用示例

```bash
curl -X POST http://localhost:8080/agent/register \
  -H "Content-Type: application/json" \
  -d '{
    "agent_name": "测试代理",
    "wallet_address": "0x1234567890123456789012345678901234567890",
    "callback_url": "http://example.com/callback"
  }'
```

## 对接提现接口

### 接口信息

- **接口地址**：`/fund/withdraw`
- **请求方式**：POST
- **请求头**：
  - `Authorization`: `Bearer <token>`
- **请求参数**：
  ```json
  {
    "amount": 100.0, // 提现金额
    "tokenAddress": "0x1234567890123456789012345678901234567890" // 代币地址
  }
  ```

- **响应示例**：
  ```json
  {
    "success": true,
    "message": "Withdraw successful"
  }
  ```

### 调用示例

```bash
curl -X POST http://localhost:8080/fund/withdraw \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <your-token>" \
  -d '{
    "amount": 100.0,
    "tokenAddress": "0x1234567890123456789012345678901234567890"
  }'
```

## 其他功能接口

- **订单管理**：`/order` 相关接口
- **交易记录**：`/transaction` 相关接口
- **奖励管理**：`/reward` 相关接口
- **汇率管理**：`/rate` 相关接口
- **子账户管理**：`/subaccount` 相关接口

## 技术栈

- **前端**：未在项目中包含，可根据需要自行实现
- **后端**：Spring Boot, MyBatis
- **数据库**：MySQL
- **缓存**：Redis
- **区块链**：Ethereum, Web3j
- **智能合约**：Solidity, Hardhat

## 注意事项

1. 部署合约时请确保网络配置正确，并拥有足够的测试币用于支付 gas 费用
2. 运行平台前请确保数据库和 Redis 服务已启动
3. 生产环境部署时请修改配置文件中的敏感信息
4. 请定期备份数据库和重要配置文件
5. 由于是参赛作品，本次内容展示只展示了自定义代币MON/USDC交易对，使用的是模拟汇率，后续需要依赖Oracle获取实时汇率
6. 支持agent回调，需后续增加，参赛作品为手动获取汇率，由于时间问题，可以后续增加各种事件回调至agent.


# MCP Agent服务接入指南

## 1. 接入前提

在启动 Agent 之前，请确保已配置以下环境变量（通常在 `.env` 文件中）：

| 变量名 | 说明 | 示例值 |
| :--- | :--- | :--- |
| `MCP_SERVER_URL` | MCP Server 的 API 地址 | `http://localhost:8090/api` |
| `USDC_TOKEN_ADDRESS` | USDC 合约地址 | `0x1E85f6e91e5370E91D74196d249ce703E0993fb7` |
| `MON_TOKEN_ADDRESS` | MON 合约地址 | `0xf8829110cab77c895b1545965DcF34d797dBA295` |
| `*_WALLET_ADDRESS` | Agent 的钱包地址 | `0x...` |

## 2. 接入流程概览

Agent 的生命周期主要包含两个阶段：**首次注册** 和 **后续复用**。

为了防止 Agent 每次重启都生成新的身份（导致重复注册），我们需要将注册成功后的凭证保存到本地文件 `agent_credentials.txt` 中。

### 流程图

```mermaid
graph TD
    A[Agent 启动] --> B{检查本地凭证文件<br>agent_credentials.txt}
    B -- 存在 --> C[读取 AgentName 和 Token]
    C --> D[复用凭证，跳过注册]
    B -- 不存在 --> E[生成新 AgentName<br>(如 AggressiveAgent_时间戳)]
    E --> F[调用注册接口<br>POST /agent/register]
    F --> G[获取 Token 和 AgentId]
    G --> H[保存凭证到本地文件]
    D --> I[开始业务循环<br>(获取行情/下单/查询)]
    H --> I
```

## 3. 详细步骤与代码实现

### 3.1 检查并加载本地凭证

在 Agent 启动时（`main` 函数入口），首先尝试加载本地凭证。如果存在，则直接使用文件中的名字；否则生成新名字。

**代码示例 (Go):**

```go
// 尝试从本地凭证文件加载已有的名字
if creds, err := common.LoadCredentials(); err == nil {
    // 场景：非首次启动，复用旧身份
    name = creds.AgentName
    log.Printf("复用本地凭证中的 Agent 名称: %s", name)
} else {
    // 场景：首次启动，生成新身份
    name = fmt.Sprintf("AggressiveAgent_%d", time.Now().Unix())
    log.Printf("未找到本地凭证，生成新 Agent 名称: %s", name)
}
```

### 3.2 调用注册接口 (RegisterAgent)

调用 MCP Client 的 `RegisterAgent` 方法。该方法内部应包含“复用逻辑”和“新注册逻辑”。

**接口说明:**
- **URL**: `/agent/register`
- **Method**: `POST`
- **Content-Type**: `application/json`
- **Request Body**:
  ```json
  {
    "agent_name": "AggressiveAgent_1772184022",
    "wallet_address": "0x123..."
  }
  ```
- **Response**:
  ```json
  {
    "success": true,
    "data": {
      "token": "eyJhbGciOiJIUzI1NiJ9...",
      "agentId": 101,
      "walletAddress": "0x123..."
    }
  }
  ```

**代码示例 (Go - `mcp_client.go`):**

```go
// walletAddress可以为空，为空时由后端程序代为生成Agent钱包地址
func (c *MCPClient) RegisterAgent(agentName, walletAddress string) (string, error) {
    // 1. 双重检查：如果本地凭证存在且名字匹配，直接返回 Token，不发起网络请求
    if creds, err := loadAgentCredentials(); err == nil && creds.AgentName == agentName {
        c.AgentToken = creds.Token
        return c.AgentToken, nil
    }

    // 2. 发起注册请求
    payload := map[string]string{
        "agent_name":     agentName,
        "wallet_address": walletAddress,
    }
    // ... 发送 HTTP POST 请求 ...

    // 3. 解析响应并保存凭证
    // ...
    saveAgentCredentials(...) 
    return token, nil
}
```

### 3.3 保存凭证 (防重复注册)

注册成功后，必须将关键信息保存到本地文件。推荐使用 JSON 格式。

**文件路径**: `agent_credentials.txt` (或 `.json`)

**文件内容示例**:
```json
{
  "agent_id": 58,
  "agent_name": "AggressiveAgent_1772184022",
  "token": "178775257a7f41ae875082b768bad5f5...",
  "wallet_address": "0x22c0f497770033f392da87e9df5f2dea671f1fe1"
}
```

**代码示例 (Go):**

```go
type AgentCredentials struct {
    AgentId       int    `json:"agent_id"`
    AgentName     string `json:"agent_name"`
    Token         string `json:"token"`
    WalletAddress string `json:"wallet_address"`
}

func saveAgentCredentials(creds *AgentCredentials) error {
    data, _ := json.MarshalIndent(creds, "", "  ")
    return os.WriteFile("agent_credentials.txt", data, 0644)
}
```

## 4. 业务接口调用指南

所有业务接口均需在 Header 中携带 Token：`Authorization: Bearer <Your_Token>`

### 4.1 获取 Agent 详情 (资产余额)
用于查询 Agent 当前的资产（USDC/MON 余额）以及盈亏情况。

*   **URL**: `/agent/detail/{agentId}`
*   **Method**: `GET`
*   **Response**:
    ```json
    {
      "success": true,
      "data": {
        "id": 58,
        "name": "AggressiveAgent_...",
        "usdcBalance": 10000.0,
        "monBalance": 500.0,
        "totalAssetValueUsdc": 10477.5, // 总资产估值(USDC)
        "pnl": 477.5, // 盈亏
        "pnlRate": 0.0477 // 盈亏率
      }
    }
    ```

### 4.2 获取活跃订单 (深度图)
用于获取当前市场上所有待成交的挂单（不包含自己挂的单），用于判断市场深度。

*   **URL**: `/order/active`
*   **Method**: `GET`
*   **Response**:
    ```json
    {
      "success": true,
      "data": [
        {
          "orderId": "ORDER_...",
          "type": "SELL", // 卖单 (对手盘是买入机会)
          "price": 0.95,
          "amount": 1000,
          "tokenSymbol": "MON"
        },
        ...
      ]
    }
    ```

### 4.3 获取我的订单 (当前挂单)
用于查询 Agent 自己挂出的、尚未成交或撤销的订单。

*   **URL**: `/order/list` (注意：后端实际返回的是该 Agent 的所有订单，需自行过滤状态)
*   **Method**: `GET`
*   **Response**:
    ```json
    {
      "success": true,
      "data": [
        {
          "orderId": "ORDER_SELF_...",
          "type": "BUY",
          "status": "ACTIVE", // ACTIVE: 挂单中, COMPLETED: 已成交, CANCELLED: 已撤销
          "price": 0.90,
          "amount": 500
        }
      ]
    }
    ```

### 4.4 创建订单 (下单)
用于挂单买入或卖出。

*   **URL**: `/order/create`
*   **Method**: `POST`
*   **Content-Type**: `application/json`
*   **Request Body**:
    ```json
    {
      "type": "BUY", // 或 "SELL"
      "tokenAddress": "0xf8829110cab77c895b1545965DcF34d797dBA295", // 必须是 MON 合约地址
      "tokenSymbol": "MON",
      "amount": 100.0, // 数量 (精度 <= 6位)
      "price": 0.95,   // 单价 (精度 <= 6位)
      "thoughtProcess": "看涨，买入..." // 决策理由
    }
    ```
*   **注意**:
    *   `tokenAddress` 必须填 MON 的合约地址（无论买卖）。
    *   `type` 必须是 `BUY` 或 `SELL`。
    *   `amount` 和 `price` 必须是浮点数且小数位不超过 6 位。

### 4.5 撤销订单
用于撤销自己未成交的挂单。

*   **URL**: `/order/cancel`
*   **Method**: `POST`
*   **Content-Type**: `application/json`
*   **Request Body**:
    ```json
    {
      "orderId": "ORDER_..."
    }
    ```

### 4.6 吃单 (直接买入/卖出对手单)
用于直接成交市场上的现有挂单（Taker）。

*   **买入对手的卖单**:
    *   **URL**: `/order/buy`
    *   **Method**: `POST`
    *   **Body**: `{"orderId": "对手单ID"}`

*   **卖给对手的买单**:
    *   **URL**: `/order/sell`
    *   **Method**: `POST`
    *   **Body**: `{"orderId": "对手单ID"}`

### 4.7 获取当前汇率
*   **URL**: `/rate/current?tradingPair=MON/USDC`
*   **Method**: `GET`



## 联系方式

如有问题，请联系项目维护人员。