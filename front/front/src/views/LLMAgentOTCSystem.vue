<template>
  <div class="container">
    <!-- Market Overview -->
  <MarketOverview 
    :current-price="currentPrice"
    :price-change="priceChange"
    :price-change-amount="priceChangeAmount"
    :order-count="orderCount"
    :otc-liquidity="otcLiquidity"
    :current-tick="currentTick"
  />

  <!-- Main Content -->
  <div class="main-content">
    <!-- Left: Agent Leaderboard -->
    <AgentLeaderboard
      :agents="sortedAgents"
      @agent-click="showAgentDetail"
    />
    
    <!-- Right: OTC Marketplace -->
    <OTCMarketplace
      :orders="filteredOrders"
      :current-filter="currentFilter"
      :filters="filters"
      @filter-change="currentFilter = $event"
    />
  </div>
  
  <!-- Transaction Flow -->
  <TransactionFlow :transactions="transactions" />

  </div>
</template>

<script setup>
import { ref, computed, onMounted, onUnmounted } from 'vue';
import { useRouter } from 'vue-router';
import MarketOverview from '../components/MarketOverview.vue';
import AgentLeaderboard from '../components/AgentLeaderboard.vue';
import OTCMarketplace from '../components/OTCMarketplace.vue';
import TransactionFlow from '../components/TransactionFlow.vue';
import { marketApi, agentApi, tradingApi, transactionApi } from '../utils/api';

const router = useRouter();

// State variables
const agents = ref([]);
const otcOrders = ref([]);
const transactions = ref([]);
const currentTick = ref(0);
const currentPrice = ref(0);
const priceChange = ref(0);
const priceChangeAmount = ref(0);
const orderCount = ref(0);
const otcLiquidity = ref(0);
const currentFilter = ref("all");
const simulationInterval = ref(null);

// Load market data
const loadMarketData = async () => {
  try {
    // Get market data
    const marketData = await marketApi.getMarketData();
    if (marketData) {
      currentPrice.value = marketData.monUsdcRate?.latestPrice || 0;
      priceChange.value = marketData.monUsdcRate?.changePercent || 0;
      priceChangeAmount.value = marketData.monUsdcRate?.changeAmount ? 
        `${marketData.monUsdcRate.changeAmount >= 0 ? '+' : ''}${marketData.monUsdcRate.changeAmount}` : 
        0;
      orderCount.value = marketData.activeOrderCount || 0;
      otcLiquidity.value = marketData.totalLiquidity || 0;
      currentTick.value = marketData.rateUpdateCount || 0;
    }
  } catch (err) {
    console.error('Failed to load market data:', err);
  }
};

// Load agent data
const loadAgentData = async () => {
  try {
    const agentRanking = await agentApi.getAgentRanking(10);
    if (agentRanking && Array.isArray(agentRanking)) {
      agents.value = agentRanking.map((item, index) => ({
        id: item.agentId.toString(),
        name: item.agentName || `Agent ${item.agentId}`,
        avatarClass: ['fox', 'bear', 'rabbit', 'wolf'][index % 4],
        totalAssets: item.totalAsset,
        dailyChange: item.profit24h,
        ethBalance: item.ethBalance,
        usdcBalance: item.usdcBalance,
        activeOrders: item.activeOrders,
        rank: item.ranking,
        role: item.role,
        completedOrders: item.completedOrders
      }));
    }
  } catch (err) {
    console.error('Failed to load agent data:', err);
  }
};

// Load order data
const loadOrderData = async () => {
  try {
    const activeOrders = await tradingApi.getActiveOrders();
    if (activeOrders && Array.isArray(activeOrders)) {
      otcOrders.value = activeOrders.map((order, index) => ({
        id: order.orderId,
        agentId: order.agentId?.toString(),
        type: order.type?.toLowerCase(),
        ethAmount: order.amount,
        usdcAmount: order.amount * order.price,
        pricePerEth: order.price,
        tokenSymbol: order.tokenSymbol || 'MON',
        time: new Date(order.createdAt).toLocaleString(),
        status: order.status,
        isNew: index < 3
      }));
    }
  } catch (err) {
    console.error('Failed to load order data:', err);
  }
};

// Load transaction data
const loadTransactionData = async () => {
  try {
    const latestTransactions = await transactionApi.getLatestTransactions(10);
    if (latestTransactions && Array.isArray(latestTransactions)) {
      transactions.value = latestTransactions.map((tx, index) => ({
        id: tx.id,
        type: tx.transactionType?.toLowerCase() || 'sell',
        fromAgent: tx.agentName || `Agent ${tx.agentId || index + 1}`,
        amount: tx.amount,
        asset: tx.tokenSymbol || 'MON',
        time: new Date(tx.timestamp).toLocaleString(),
        agentId: tx.agentId,
        agentName: tx.agentName,
        transactionType: tx.transactionType,
        tokenSymbol: tx.tokenSymbol || 'MON',
        price: tx.price,
        totalValue: tx.totalValue,
        fee: tx.fee,
        status: tx.status,
        relatedOrderId: tx.relatedOrderId,
        timestamp: tx.timestamp,
        description: tx.description
      }));
    }
  } catch (err) {
    console.error('Failed to load transaction data:', err);
  }
};

// Load data when component mounts
onMounted(async () => {
  await loadMarketData();
  await loadOrderData();
  await loadTransactionData();
  await loadAgentData();
  simulationInterval.value = setInterval(() => {
    reloadData();
  }, 60000);
});

// Reload data
const reloadData = async () => {
  await loadMarketData();
  await loadOrderData();
  await loadTransactionData();
  await loadAgentData();
};

// Filter options
const filters = [
  { value: 'all', label: 'All Orders' },
  { value: 'buy', label: 'Buy MON' },
  { value: 'sell', label: 'Sell MON' }
];

// Computed properties
const sortedAgents = computed(() => {
  // Sort by 24h profit
  return [...agents.value]
});



const filteredOrders = computed(() => {
  if (currentFilter.value === 'all') {
    return otcOrders.value;
  }
  return otcOrders.value.filter(order => order.type === currentFilter.value);
});

const showAgentDetail = (agentId) => {
  router.push(`/agent/${agentId}`);
};

onUnmounted(() => {
  if (simulationInterval.value) {
    clearInterval(simulationInterval.value);
  }
});
</script>
<style scoped>
.container {
  max-width: 1440px;
  margin: 0 auto;
  padding: 2rem;
  position: relative;
}





/* 主体区域 - 两栏布局 */
.main-content {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  gap: 20px;
}

@keyframes spin {
  from { transform: rotate(0deg); }
  to { transform: rotate(360deg); }
}

/* 响应式设计 */
@media (max-width: 768px) {
  .main-content {
    flex-direction: column;
  }
}
</style>