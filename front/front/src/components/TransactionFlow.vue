<template>
  <section class="transaction-flow">
    <div class="section-header">
      <h2 class="section-title">
        <i class="fas fa-exchange-alt"></i>
        Latest Transactions
      </h2>
      <span id="transaction-count">{{ transactions.length }} Transactions</span>
    </div>
    
    <div class="transaction-list" id="transaction-list">
      <div 
        v-for="transaction in transactions" 
        :key="transaction.id"
        class="transaction-item"
      >
        <div style="display:flex;align-items:center;">
            <div class="transaction-icon">
              <i :class="transaction.type === 'sell' ? 'fas fa-arrow-up' : 'fas fa-arrow-down'"></i>
            </div>
            <div class="transaction-details">
              <div class="transaction-meta">
                {{ getTransactionDescription(transaction) }}
                <span><i class="fas fa-tag"></i> {{ transaction.tokenSymbol }}</span>
                <span><i class="fas fa-dollar-sign"></i> {{ (transaction.price || 0) }} USDC</span>
                <span><i class="fas fa-check-circle"></i> {{ transaction.status }}</span>
              </div>
              <div class="transaction-time">{{ transaction.time }}</div>
            </div>
          </div>
        
        <div class="transaction-amount" :class="transaction.type === 'sell' ? 'positive' : 'negative'">
          {{ transaction.type === 'sell' ? '+' : '-' }}{{ transaction.amount }} {{ transaction.asset }}
        </div>
      </div>
    </div>
  </section>
</template>

<script setup>
import { defineProps } from 'vue';

const props = defineProps({
  transactions: {
    type: Array,
    default: () => []
  }
});

const getTransactionDescription = (transaction) => {
  if (transaction.type === 'sell') {
    return `${transaction.agentName || transaction.fromAgent} sold ${transaction.amount} ${transaction.tokenSymbol}`;
  } else if (transaction.type === 'buy') {
    return `${transaction.agentName || transaction.fromAgent} bought ${transaction.amount} ${transaction.tokenSymbol}`;
  }
  return '';
};
</script>

<style scoped>
.transaction-flow {
  background-color: var(--bg-secondary);
  border-radius: 12px;
  padding: 25px;
  margin: 25px 0;
  box-shadow: 0 4px 6px rgba(0, 0, 0, 0.1);
}

.section-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
}

.section-title {
  font-size: 1.3rem;
  font-weight: 600;
  display: flex;
  align-items: center;
  gap: 10px;
}

.transaction-list {
  margin-top: 15px;
  max-height: 300px;
  overflow-y: auto;
}

.transaction-item {
  padding: 12px 15px;
  border-bottom: 1px solid var(--border-color);
  display: flex;
  justify-content:space-between;
  align-items: center;
  gap: 12px;
}
.transaction-item:hover {
  background-color: var(--bg-card);
  transform: translateY(-2px);
}
.transaction-item:last-child {
  border-bottom: none;
}

.transaction-icon {
  width: 36px;
  height: 36px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  background-color: var(--bg-card);
}

.transaction-details {
  margin-left: 12px;
  flex: 1;
}

.transaction-amount {
  font-weight: 600;
  color: var(--accent-green);
}

.transaction-time {
  text-align:left;
  font-size: 0.8rem;
  color: var(--text-secondary);
}

.transaction-meta {
  display: flex;
  gap: 15px;
  font-size: 0.8rem;
  color: var(--text-secondary);
  margin: 5px 0;
}

.transaction-meta i {
  margin-right: 3px;
}

.positive {
  color: var(--accent-green);
}

/* 滚动条样式 */
.transaction-list::-webkit-scrollbar {
  width: 6px;
}

.transaction-list::-webkit-scrollbar-track {
  background: var(--surface-tertiary);
  border-radius: 3px;
}

.transaction-list::-webkit-scrollbar-thumb {
  background: var(--text-tertiary);
  border-radius: 3px;
}

.transaction-list::-webkit-scrollbar-thumb:hover {
  background: var(--text-secondary);
}
</style>