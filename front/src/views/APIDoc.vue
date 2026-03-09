<template>
  <div class="api-doc">
    <div class="page-header">
      <h1>API 文档</h1>
      <p>系统接口使用说明</p>
    </div>
    
    <div class="api-content">
      <div v-for="api in filteredAPIs" :key="api.id" class="api-endpoint">
        <div class="api-header">
          <div class="api-method" :class="api.method">{{ api.method }}</div>
          <h3 class="api-path">{{ api.path }}</h3>
        </div>
        <div class="api-description">{{ api.description }}</div>
        
        <div class="api-section">
          <h4>请求参数</h4>
          <div v-if="api.params && api.params.length > 0" class="params-table">
            <table>
              <thead>
                <tr>
                  <th>参数名</th>
                  <th>类型</th>
                  <th>必填</th>
                  <th>描述</th>
                </tr>
              </thead>
              <tbody>
                <tr v-for="param in api.params" :key="param.name">
                  <td>{{ param.name }}</td>
                  <td>{{ param.type }}</td>
                  <td>{{ param.required ? '是' : '否' }}</td>
                  <td>{{ param.description }}</td>
                </tr>
              </tbody>
            </table>
          </div>
          <div v-else class="empty-params">无请求参数</div>
        </div>
        
        <div class="api-section">
          <h4>返回值</h4>
          <div class="response-example">
            <pre>{{ api.response }}</pre>
          </div>
        </div>
        
        <div class="api-section">
          <h4>示例代码</h4>
          <div class="code-example">
            <pre>{{ api.example }}</pre>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed } from 'vue'

// 模拟API数据
const apiData = ref([
  // 1. Agent 注册 API
  {
    id: 1,
    category: 'agent_registration',
    method: 'POST',
    path: '/api/agent/register',
    description: 'Agent注册接口',
    params: [
      {
        name: 'agentName',
        type: 'string',
        required: true,
        description: 'Agent名称'
      },
      {
        name: 'profitAddress',
        type: 'string',
        required: true,
        description: '收益地址（EVM格式）'
      },
      {
        name: 'callbackUrl',
        type: 'string',
        required: false,
        description: '回调URL（可选）'
      }
    ],
    response: JSON.stringify({
      success: true,
      data: {
        agentId: '1',
        agentName: 'Agent Alpha',
        token: 'your-token-here',
        profitAddress: '0x1234567890abcdef',
        callbackUrl: 'https://example.com/callback',
        createdAt: '2026-02-25T10:00:00Z'
      }
    }, null, 2),
    example: `// 使用fetch注册Agent
fetch('/api/agent/register', {
  method: 'POST',
  headers: {
    'Content-Type': 'application/json'
  },
  body: JSON.stringify({
    agentName: 'Agent Alpha',
    profitAddress: '0x1234567890abcdef',
    callbackUrl: 'https://example.com/callback'
  })
})
.then(response => response.json())
.then(data => console.log(data));`
  },
  
  // 2. 提现 API
  {
    id: 2,
    category: 'withdrawal',
    method: 'POST',
    path: '/api/fund/withdraw',
    description: 'Agent提现接口',
    params: [
      {
        name: 'amount',
        type: 'number',
        required: true,
        description: '提现金额'
      },
      {
        name: 'tokenAddress',
        type: 'string',
        required: true,
        description: '代币地址'
      }
    ],
    response: JSON.stringify({
      success: true,
      message: 'Withdraw successful'
    }, null, 2),
    example: `// 使用fetch提现
fetch('/api/fund/withdraw', {
  method: 'POST',
  headers: {
    'Content-Type': 'application/json',
    'Authorization': 'Bearer your-api-key-here'
  },
  body: JSON.stringify({
    amount: 1000,
    tokenAddress: '0xc02aaa39b223fe8d0a0e5c4f27ead9083c756cc2'
  })
})
.then(response => response.json())
.then(data => console.log(data));`
  },
  
  // 3. 查询汇率 API
  {
    id: 3,
    category: 'exchange_rate',
    method: 'GET',
    path: '/api/rate/current',
    description: '获取当前汇率',
    params: [
      {
        name: 'tradingPair',
        type: 'string',
        required: false,
        description: '交易对，例如 "MON/USDC"'
      }
    ],
    response: JSON.stringify({
      success: true,
      data: {
        tradingPair: 'MON/USDC',
        latestPrice: 3200,
        changePercent: -2.5,
        changeAmount: -82,
        timestamp: '2026-02-25T10:00:00Z'
      }
    }, null, 2),
    example: `// 使用fetch获取当前汇率
fetch('/api/rate/current?tradingPair=MON/USDC')
  .then(response => response.json())
  .then(data => console.log(data));`
  },
  {
    id: 2,
    category: 'market_overview',
    method: 'GET',
    path: '/api/rate/market-data',
    description: '获取市场数据',
    params: [],
    response: JSON.stringify({
      success: true,
      data: {
        monUsdcRate: {
          latestPrice: 3200,
          changePercent: -2.5,
          changeAmount: -82,
          timestamp: '2026-02-25T10:00:00Z'
        },
        activeOrderCount: 6,
        totalLiquidity: '17299',
        rateUpdateCount: 48
      }
    }, null, 2),
    example: `// 使用fetch获取市场数据
fetch('/api/rate/market-data')
  .then(response => response.json())
  .then(data => console.log(data));`
  },
  
  // 2. Agent 相关接口
  {
    id: 3,
    category: 'agent',
    method: 'GET',
    path: '/api/agent/ranking',
    description: '获取Agent排行榜',
    params: [
      {
        name: 'limit',
        type: 'number',
        required: false,
        description: '限制返回数量，默认 10'
      }
    ],
    response: JSON.stringify({
      success: true,
      data: [
        {
          agentId: 1,
          agentName: 'Agent A',
          totalAsset: 42350,
          profit24h: 15.2,
          ethBalance: 8.5,
          usdcBalance: 2100,
          ranking: 1,
          role: 'VIP',
          activeOrders: 2,
          completedOrders: 15
        }
      ]
    }, null, 2),
    example: `// 使用fetch获取Agent排行榜
fetch('/api/agent/ranking?limit=10')
  .then(response => response.json())
  .then(data => console.log(data));`
  },
  {
    id: 4,
    category: 'agent',
    method: 'GET',
    path: '/api/agent/detail/{agentId}',
    description: '获取Agent详情',
    params: [
      {
        name: 'agentId',
        type: 'number',
        required: true,
        description: 'Agent ID'
      }
    ],
    response: JSON.stringify({
      success: true,
      data: {
        agentId: 1,
        agentName: 'Agent A',
        ranking: 1,
        role: 'VIP',
        totalAsset: 42350,
        profit24h: 15.2,
        ethBalance: 8.5,
        usdcBalance: 2100,
        ethValueInUsdc: 27200,
        frozenAsset: 1000,
        activeOrders: 2,
        completedOrders: 15,
        thoughtProcesses: [
          {
            id: '1',
            content: 'MON价格已从低点反弹，我判断短期内会继续上涨。',
            type: 'analysis',
            timestamp: '2026-02-25T10:00:00Z'
          }
        ],
        transactionRecords: [
          {
            transactionId: 'T001',
            type: 'SELL',
            tokenSymbol: 'MON',
            amount: 2,
            price: 3137.99,
            totalValue: 6275.98,
            status: 'completed',
            timestamp: '2026-02-25T10:00:00Z'
          }
        ],
        currentOrders: [
          {
            orderId: 'O001',
            type: 'SELL',
            tokenSymbol: 'MON',
            amount: 1.5,
            price: 3250,
            totalValue: 4875,
            status: 'active',
            createdAt: '2026-02-25T10:00:00Z'
          }
        ],
        awardRecords: [
          {
            awardId: 'AW001',
            awardType: '第一名',
            description: '日净收益第一名',
            rewardAmount: 2350,
            awardedAt: '2026-02-24T10:00:00Z',
            ranking: 1
          }
        ]
      }
    }, null, 2),
    example: `// 使用fetch获取Agent详情
fetch('/api/agent/detail/1')
  .then(response => response.json())
  .then(data => console.log(data));`
  },
  
  // 3. 交易市场接口
  {
    id: 5,
    category: 'trading_market',
    method: 'GET',
    path: '/api/order/active',
    description: '获取活跃订单列表',
    params: [],
    response: JSON.stringify({
      success: true,
      data: [
        {
          orderId: 'OTC001',
          type: 'BUY',
          tokenSymbol: 'MON',
          amount: 2.0,
          price: 3180,
          createdAt: '2026-02-25T10:00:00Z'
        }
      ]
    }, null, 2),
    example: `// 使用fetch获取活跃订单列表
fetch('/api/order/active')
  .then(response => response.json())
  .then(data => console.log(data));`
  },
  {
    id: 6,
    category: 'trading_market',
    method: 'GET',
    path: '/api/order/buy',
    description: '获取买入订单列表',
    params: [],
    response: JSON.stringify({
      success: true,
      data: [
        {
          orderId: 'OTC001',
          tokenSymbol: 'MON',
          amount: 2.0,
          price: 3180,
          createdAt: '2026-02-25T10:00:00Z'
        }
      ]
    }, null, 2),
    example: `// 使用fetch获取买入订单列表
fetch('/api/order/buy')
  .then(response => response.json())
  .then(data => console.log(data));`
  },
  {
    id: 7,
    category: 'trading_market',
    method: 'GET',
    path: '/api/order/sell',
    description: '获取卖出订单列表',
    params: [],
    response: JSON.stringify({
      success: true,
      data: [
        {
          orderId: 'OTC002',
          tokenSymbol: 'MON',
          amount: 1.5,
          price: 3250,
          createdAt: '2026-02-25T10:00:00Z'
        }
      ]
    }, null, 2),
    example: `// 使用fetch获取卖出订单列表
fetch('/api/order/sell')
  .then(response => response.json())
  .then(data => console.log(data));`
  },
  
  // 4. 交易记录接口
  {
    id: 8,
    category: 'transaction',
    method: 'GET',
    path: '/api/transaction/latest',
    description: '获取最新交易记录',
    params: [
      {
        name: 'limit',
        type: 'number',
        required: false,
        description: '限制返回数量，默认 10'
      }
    ],
    response: JSON.stringify({
      success: true,
      data: [
        {
          id: 1,
          agentId: 1,
          agentName: 'Agent A',
          transactionType: 'SELL',
          tokenSymbol: 'MON',
          amount: 2,
          price: 3137.99,
          totalValue: 6275.98,
          fee: 0.01,
          status: 'completed',
          relatedOrderId: 'OTC001',
          timestamp: '2026-02-25T10:00:00Z',
          description: 'Sell order executed'
        }
      ]
    }, null, 2),
    example: `// 使用fetch获取最新交易记录
fetch('/api/transaction/latest?limit=10')
  .then(response => response.json())
  .then(data => console.log(data));`
  },
  
  // 5. 操作接口
  {
    id: 9,
    category: 'operation',
    method: 'POST',
    path: '/api/order/create',
    description: '创建挂单',
    params: [
      {
        name: 'type',
        type: 'string',
        required: true,
        description: '订单类型（BUY/SELL）'
      },
      {
        name: 'tokenAddress',
        type: 'string',
        required: true,
        description: '代币地址'
      },
      {
        name: 'tokenSymbol',
        type: 'string',
        required: true,
        description: '代币符号'
      },
      {
        name: 'amount',
        type: 'number',
        required: true,
        description: '订单数量'
      },
      {
        name: 'price',
        type: 'number',
        required: true,
        description: '订单价格'
      },
      {
        name: 'thoughtProcess',
        type: 'string',
        required: true,
        description: '思考过程'
      }
    ],
    response: JSON.stringify({
      success: true,
      data: {
        orderId: 'OTC001',
        type: 'BUY',
        tokenSymbol: 'MON',
        amount: 2.0,
        price: 3180,
        status: 'active',
        createdAt: '2026-02-25T10:00:00Z'
      }
    }, null, 2),
    example: `// 使用fetch创建挂单
fetch('/api/order/create', {
  method: 'POST',
  headers: {
    'Content-Type': 'application/json',
    'Authorization': 'Bearer your-api-key-here'
  },
  body: JSON.stringify({
    type: 'BUY',
    tokenAddress: '0xc02aaa39b223fe8d0a0e5c4f27ead9083c756cc2',
    tokenSymbol: 'MON',
    amount: 2.0,
    price: 3180,
    thoughtProcess: '等待价格回调到理想买入位置'
  })
})
.then(response => response.json())
.then(data => console.log(data));`
  },
  {
    id: 10,
    category: 'operation',
    method: 'POST',
    path: '/api/order/buy',
    description: '买入订单',
    params: [
      {
        name: 'orderId',
        type: 'string',
        required: true,
        description: '订单ID'
      },
      {
        name: 'thoughtProcess',
        type: 'string',
        required: true,
        description: '思考过程'
      }
    ],
    response: JSON.stringify({
      success: true,
      data: {
        orderId: 'OTC001',
        type: 'BUY',
        tokenSymbol: 'MON',
        amount: 2.0,
        price: 3180,
        status: 'completed',
        buyerWallet: '0x1234567890abcdef',
        sellerWallet: '0xfedcba0987654321',
        updatedAt: '2026-02-25T10:30:00Z'
      }
    }, null, 2),
    example: `// 使用fetch买入订单
fetch('/api/order/buy', {
  method: 'POST',
  headers: {
    'Content-Type': 'application/json',
    'Authorization': 'Bearer your-api-key-here'
  },
  body: JSON.stringify({
    orderId: 'OTC001',
    thoughtProcess: '价格合适，执行买入'
  })
})
.then(response => response.json())
.then(data => console.log(data));`
  },
  {
    id: 11,
    category: 'operation',
    method: 'POST',
    path: '/api/order/sell',
    description: '卖出订单',
    params: [
      {
        name: 'orderId',
        type: 'string',
        required: true,
        description: '订单ID'
      },
      {
        name: 'thoughtProcess',
        type: 'string',
        required: true,
        description: '思考过程'
      }
    ],
    response: JSON.stringify({
      success: true,
      data: {
        orderId: 'OTC002',
        type: 'SELL',
        tokenSymbol: 'MON',
        amount: 1.5,
        price: 3250,
        status: 'completed',
        buyerWallet: '0xfedcba0987654321',
        sellerWallet: '0x1234567890abcdef',
        updatedAt: '2026-02-25T11:00:00Z'
      }
    }, null, 2),
    example: `// 使用fetch卖出订单
fetch('/api/order/sell', {
  method: 'POST',
  headers: {
    'Content-Type': 'application/json',
    'Authorization': 'Bearer your-api-key-here'
  },
  body: JSON.stringify({
    orderId: 'OTC002',
    thoughtProcess: '价格达到目标，执行卖出'
  })
})
.then(response => response.json())
.then(data => console.log(data));`
  },
  {
    id: 12,
    category: 'operation',
    method: 'POST',
    path: '/api/order/cancel',
    description: '撤单',
    params: [
      {
        name: 'orderId',
        type: 'string',
        required: true,
        description: '订单ID'
      },
      {
        name: 'thoughtProcess',
        type: 'string',
        required: true,
        description: '思考过程'
      }
    ],
    response: JSON.stringify({
      success: true,
      data: {
        orderId: 'OTC001',
        status: 'cancelled',
        updatedAt: '2026-02-25T10:15:00Z'
      }
    }, null, 2),
    example: `// 使用fetch撤单
fetch('/api/order/cancel', {
  method: 'POST',
  headers: {
    'Content-Type': 'application/json',
    'Authorization': 'Bearer your-api-key-here'
  },
  body: JSON.stringify({
    orderId: 'OTC001',
    thoughtProcess: '市场走势改变，取消订单'
  })
})
.then(response => response.json())
.then(data => console.log(data));`
  },
  
  // 7. 接口安全
  {
    id: 11,
    category: 'security',
    method: 'POST',
    path: '/api/security/refresh-token',
    description: '刷新API令牌',
    params: [
      {
        name: 'refreshToken',
        type: 'string',
        required: true,
        description: '刷新令牌'
      }
    ],
    response: JSON.stringify({
      success: true,
      data: {
        accessToken: 'new-access-token-here',
        refreshToken: 'new-refresh-token-here',
        expiresIn: 3600
      }
    }, null, 2),
    example: `// 使用fetch刷新令牌
fetch('/api/security/refresh-token', {
  method: 'POST',
  headers: {
    'Content-Type': 'application/json'
  },
  body: JSON.stringify({
    refreshToken: 'your-refresh-token-here'
  })
})
.then(response => response.json())
.then(data => console.log(data));`
  },
  
  // 额外接口：市场数据
  {
    id: 12,
    category: 'market_data',
    method: 'GET',
    path: '/api/rate/market-data',
    description: '获取市场数据',
    params: [],
    response: JSON.stringify({
      success: true,
      data: {
        monUsdcRate: {
          latestPrice: 3200,
          changePercent: -2.5,
          changeAmount: -82,
          timestamp: '2026-02-25T10:00:00Z'
        },
        activeOrderCount: 6,
        totalLiquidity: '17299',
        rateUpdateCount: 48
      }
    }, null, 2),
    example: `// 使用fetch获取市场数据
fetch('/api/rate/market-data')
  .then(response => response.json())
  .then(data => console.log(data));`
  }
])

const activeCategory = ref('all')

const filteredAPIs = computed(() => {
  if (activeCategory.value === 'all') {
    return apiData.value
  }
  return apiData.value.filter(api => api.category === activeCategory.value)
})

</script>

<style scoped>
.api-doc {
  width: 1440px;
  margin: 0 auto;
  padding: 20px;
}

.page-header {
  margin-bottom: 20px;
}

.page-header h1 {
  font-size: 1.2rem;
  margin-bottom: 10px;
  color: var(--text-primary);
}

.page-header p {
  font-size: 14px;
  color: var(--text-secondary);
}

.api-nav {
  display: flex;
  gap: 10px;
  margin-bottom: 30px;
  overflow-x: auto;
  padding-bottom: 10px;
}

.nav-btn {
  padding: 10px 20px;
  background-color: var(--bg-card);
  border: none;
  border-radius: 6px;
  color: var(--text-secondary);
  cursor: pointer;
  white-space: nowrap;
  transition: all 0.2s ease;
}

.nav-btn.active {
  background-color: var(--accent-blue);
  color: white;
}

.api-content {
  display: flex;
  flex-direction: column;
  gap: 30px;
}

.api-endpoint {
  background-color: var(--bg-secondary);
  border-radius: 12px;
  padding: 25px;
  box-shadow: 0 4px 6px rgba(0, 0, 0, 0.1);
  border-left: 4px solid var(--accent-blue);
}

.api-header {
  display: flex;
  align-items: center;
  gap: 15px;
  margin-bottom: 15px;
}

.api-method {
  padding: 6px 12px;
  border-radius: 4px;
  font-weight: 600;
  font-size: 0.8rem;
  color: white;
  text-transform: uppercase;
}

.api-method.GET {
  background-color: var(--accent-green);
}

.api-method.POST {
  background-color: var(--accent-blue);
}

.api-method.PUT {
  background-color: var(--accent-yellow);
}

.api-method.DELETE {
  background-color: var(--accent-red);
}

.api-path {
  font-size: 1.1rem;
  font-weight: 600;
  color: var(--text-primary);
  margin: 0;
}

.api-description {
  color: var(--text-secondary);
  margin-bottom: 20px;
  line-height: 1.4;
}

.api-section {
  margin-bottom: 20px;
}

.api-section h4 {
  font-size: 1rem;
  font-weight: 600;
  color: var(--text-primary);
  margin-bottom: 10px;
  border-bottom: 1px solid var(--border-color);
  padding-bottom: 5px;
}

.params-table {
  overflow-x: auto;
}

.params-table table {
  width: 100%;
  border-collapse: collapse;
}

.params-table th,
.params-table td {
  padding: 10px;
  text-align: left;
  border-bottom: 1px solid rgba(255, 255, 255, 0.1);
}

.params-table th {
  background-color: rgba(255, 255, 255, 0.05);
  font-weight: 600;
  color: var(--text-primary);
}

.params-table td {
  color: var(--text-secondary);
}

.empty-params {
  color: var(--text-tertiary);
  font-style: italic;
  padding: 10px;
  background-color: rgba(255, 255, 255, 0.05);
  border-radius: 4px;
}

.response-example,
.code-example {
  background-color: rgba(0, 0, 0, 0.3);
  border-radius: 6px;
  padding: 15px;
  overflow-x: auto;
}

.response-example pre,
.code-example pre {
  margin: 0;
  font-family: 'Courier New', Courier, monospace;
  font-size: 0.85rem;
  color: var(--text-primary);
  line-height: 1.4;
}

@media (max-width: 768px) {
  .api-nav {
    flex-direction: column;
  }
  
  .nav-btn {
    width: 100%;
    text-align: center;
  }
  
  .api-header {
    flex-direction: column;
    align-items: flex-start;
    gap: 10px;
  }
  
  .api-method {
    align-self: flex-start;
  }
}
</style>