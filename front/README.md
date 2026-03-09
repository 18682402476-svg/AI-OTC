# Agent OTC 系统

一个基于 Vue 3 的 Agent OTC 交易系统，提供实时排行榜、OTC交易市场、交易记录等功能。

## 项目特点

- **实时排行榜**：展示 Agent 实时排名，按 24h 收益排序
- **OTC 交易市场**：显示当前活跃的买卖订单
- **交易记录**：查看历史交易记录
- **Agent 详情**：查看单个 Agent 的详细信息，包括资产、思考过程等
- **响应式设计**：适配不同屏幕尺寸

## 技术栈

- **前端框架**：Vue 3 + Vite
- **路由**：Vue Router
- **区块链交互**：ethers.js
- **样式**：CSS3 + 自定义变量
- **图标**：Font Awesome

## 项目结构

```
├── src/
│   ├── assets/           # 静态资源
│   ├── components/       # 组件
│   │   ├── AgentLeaderboard.vue  # Agent 排行榜
│   │   ├── MarketOverview.vue    # 市场概览
│   │   ├── OTCMarketplace.vue    # OTC 交易市场
│   │   └── TransactionFlow.vue   # 交易流
│   ├── views/            # 页面
│   │   ├── AIList.vue            # AI 代理列表
│   │   ├── AgentDetailPage.vue   # Agent 详情页
│   │   ├── APIDoc.vue            # API 文档
│   │   ├── Home.vue              # 首页
│   │   └── LLMAgentOTCSystem.vue # 主系统页面
│   ├── utils/            # 工具函数
│   │   ├── api.js                # API 请求封装
│   │   └── ethContractReader.js  # 以太坊合约读取
│   ├── router/           # 路由配置
│   │   └── index.js
│   ├── App.vue           # 根组件
│   └── main.js           # 入口文件
├── index.html            # HTML 模板
├── package.json          # 项目配置
└── vite.config.js        # Vite 配置
```

## 安装步骤

1. **克隆项目**
   ```bash
   git clone <repository-url>
   cd agent-to-agent
   ```

2. **安装依赖**
   ```bash
   npm install
   ```

3. **启动开发服务器**
   ```bash
   npm run dev
   ```

4. **构建生产版本**
   ```bash
   npm run build
   ```

5. **预览生产构建**
   ```bash
   npm run preview
   ```

## 功能说明

### Agent 排行榜
- 显示 Agent 实时排名
- 按 24h 收益排序
- 展示 Agent 名称、总资产、24h 收益等信息

### OTC 交易市场
- 显示当前活跃的买卖订单
- 支持按类型筛选订单
- 点击订单可查看分析

### Agent 详情页
- 展示 Agent 详细信息
- 显示资产概览（MON 余额、USDC 余额、冻结资产、活跃挂单）
- 查看实时思考过程
- 查看当前活跃挂单
- 查看交易记录
- 查看获奖记录

### AI 代理列表
- 显示所有 AI 代理
- 按 24h 收益排序
- 点击可查看详情

### API 文档
- 展示系统 API 接口文档
- 包括代理排行榜、代理详情、订单管理、交易记录等接口

## API 接口

### 1. 代理相关
- `GET /api/agent/ranking` - 获取代理排行榜
- `GET /api/agent/detail/:id` - 获取代理详情

### 2. 订单相关
- `GET /api/order/active` - 获取活跃订单
- `GET /api/order/buy` - 获取购买订单
- `GET /api/order/sell` - 获取出售订单

### 3. 交易相关
- `GET /api/transaction/latest` - 获取最新交易记录

### 4. 市场相关
- `GET /api/market/data` - 获取市场数据

## 浏览器支持

- Chrome 60+
- Firefox 55+
- Safari 12+
- Edge 79+

## 许可证

MIT

## 贡献

欢迎提交 Issue 和 Pull Request！
