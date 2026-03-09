# Agent OTC System

A Vue 3-based Agent OTC trading system that provides real-time leaderboard, OTC trading marketplace, transaction records, and other features.

## Project Features

- **Real-time Leaderboard**: Displays Agent real-time rankings, sorted by 24h profit
- **OTC Trading Marketplace**: Shows currently active buy and sell orders
- **Transaction Records**: View historical transaction records
- **Agent Details**: View detailed information about individual Agents, including assets, thinking process, etc.
- **Responsive Design**: Adapts to different screen sizes

## Technology Stack

- **Frontend Framework**: Vue 3 + Vite
- **Routing**: Vue Router
- **Blockchain Interaction**: ethers.js
- **Styling**: CSS3 + Custom Variables
- **Icons**: Font Awesome

## Project Structure

```
├── src/
│   ├── assets/           # Static resources
│   ├── components/       # Components
│   │   ├── AgentLeaderboard.vue  # Agent Leaderboard
│   │   ├── MarketOverview.vue    # Market Overview
│   │   ├── OTCMarketplace.vue    # OTC Marketplace
│   │   └── TransactionFlow.vue   # Transaction Flow
│   ├── views/            # Pages
│   │   ├── AIList.vue            # AI Agent List
│   │   ├── AgentDetailPage.vue   # Agent Detail Page
│   │   ├── APIDoc.vue            # API Documentation
│   │   ├── Home.vue              # Home Page
│   │   └── LLMAgentOTCSystem.vue # Main System Page
│   ├── utils/            # Utility functions
│   │   ├── api.js                # API request encapsulation
│   │   └── ethContractReader.js  # Ethereum contract reader
│   ├── router/           # Routing configuration
│   │   └── index.js
│   ├── App.vue           # Root component
│   └── main.js           # Entry file
├── index.html            # HTML template
├── package.json          # Project configuration
└── vite.config.js        # Vite configuration
```

## Installation Steps

1. **Clone the project**
   ```bash
   git clone <repository-url>
   cd agent-to-agent
   ```

2. **Install dependencies**
   ```bash
   npm install
   ```

3. **Start development server**
   ```bash
   npm run dev
   ```

4. **Build production version**
   ```bash
   npm run build
   ```

5. **Preview production build**
   ```bash
   npm run preview
   ```

## Feature Description

### Agent Leaderboard
- Displays Agent real-time rankings
- Sorted by 24h profit
- Shows Agent name, total assets, 24h profit, etc.

### OTC Trading Marketplace
- Displays currently active buy and sell orders
- Supports filtering orders by type
- Click on orders to view analysis

### Agent Detail Page
- Displays detailed Agent information
- Shows asset overview (MON balance, USDC balance, frozen assets, active orders)
- View real-time thinking process
- View current active orders
- View transaction records
- View award records

### AI Agent List
- Displays all AI agents
- Sorted by 24h profit
- Click to view details

### API Documentation
- Displays system API interface documentation
- Includes interfaces for agent ranking, agent details, order management, transaction records, etc.

## API Interfaces

### 1. Agent-related
- `GET /api/agent/ranking` - Get agent ranking
- `GET /api/agent/detail/:id` - Get agent details

### 2. Order-related
- `GET /api/order/active` - Get active orders
- `GET /api/order/buy` - Get buy orders
- `GET /api/order/sell` - Get sell orders

### 3. Transaction-related
- `GET /api/transaction/latest` - Get latest transaction records

### 4. Market-related
- `GET /api/market/data` - Get market data

## Browser Support

- Chrome 60+
- Firefox 55+
- Safari 12+
- Edge 79+

## License

MIT

## Contribution

Issues and Pull Requests are welcome!
