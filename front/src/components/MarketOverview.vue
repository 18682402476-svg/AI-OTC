<template>
  <div class="market-overview">
    <div class="stat-card">
      <div class="stat-title">
        <i class="fas fa-chart-line"></i>
        <span>当前价格</span>
      </div>
      <div class="stat-value">{{ currentPrice }}</div>
      <div class="stat-change">
        <span>MON/USDC</span>
      </div>
    </div>
    
    <div class="stat-card">
      <div class="stat-title">
        <i class="fas fa-exchange-alt"></i>
        <span>24H变化</span>
      </div>
      <div class="stat-value">{{ priceChange }}%</div>
      <div class="stat-change" :class="{ negative: priceChange < 0, positive: priceChange >= 0 }">
        <i :class="priceChange < 0 ? 'fas fa-arrow-down' : 'fas fa-arrow-up'"></i>
        <span>{{ priceChangeAmount }} USDC</span>
      </div>
    </div>
    
    <div class="stat-card">
      <div class="stat-title">
        <i class="fas fa-layer-group"></i>
        <span>活跃OTC订单</span>
      </div>
      <div class="stat-value">{{ orderCount }}</div>
      <div class="stat-change">
        <span>可供交易</span>
      </div>
    </div>
    
    <div class="stat-card">
      <div class="stat-title">
        <i class="fas fa-money-bill-wave"></i>
        <span>OTC总流动金额</span>
      </div>
      <div class="stat-value">{{ otcLiquidity }}</div>
      <div class="stat-change">
        <span>USDC</span>
      </div>
    </div>
    
  </div>
</template>

<script setup>
import { defineProps } from 'vue';

const props = defineProps({
  currentPrice: {
    type: Number,
    required: true
  },
  priceChange: {
    type: Number,
    required: true
  },
  priceChangeAmount: {
    type: String,
    required: true
  },
  orderCount: {
    type: Number,
    required: true
  },
  otcLiquidity: {
    type: Number,
    required: true
  },
  currentTick: {
    type: Number,
    required: true
  }
});
</script>

<style scoped>
.market-overview {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(160px, 1fr));
  gap: 15px;
}

.stat-card {
  background-color: var(--bg-secondary);
  border-radius: 12px;
  padding: 20px;
  border-left: 4px solid var(--accent-blue);
  box-shadow: 0 4px 6px rgba(0, 0, 0, 0.1);
}

.stat-card:nth-child(2) {
  border-left-color: var(--accent-green);
}

.stat-card:nth-child(3) {
  border-left-color: var(--accent-purple);
}

.stat-card:nth-child(4) {
  border-left-color: var(--accent-red);
}

.stat-card:nth-child(5) {
  border-left-color: var(--accent-yellow);
}

.stat-title {
  font-size: 0.9rem;
  color: var(--text-secondary);
  margin-bottom: 8px;
  display: flex;
  align-items: center;
  gap: 8px;
}

.stat-value {
  font-size: 2rem;
  font-weight: 700;
  margin-bottom: 5px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.stat-change {
  font-size: 0.9rem;
  display: flex;
  align-items: center;
  gap: 5px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.positive {
  color: var(--accent-green);
}

.negative {
  color: var(--accent-red);
}
</style>