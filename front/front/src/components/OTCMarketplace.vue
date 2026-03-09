<template>
  <section class="otc-marketplace" id="otc-marketplace">
    <div class="otc-header">
      <h2 class="section-title">
        <i class="fas fa-store"></i>
        OTC Marketplace
      </h2>
      <div class="otc-filters">
        <button 
          v-for="filter in filters" 
          :key="filter.value"
          :class="['filter-btn', { active: currentFilter === filter.value }]"
          :data-filter="filter.value"
          @click="updateFilter(filter.value)"
        >{{ filter.label }}</button>
      </div>
    </div>
    
    <div class="otc-list" id="otc-list">
      <div 
        v-for="order in filteredOrders" 
        :key="order.id"
        :class="['otc-item', order.type, order.isNew ? 'new-item' : '']"
      >
        <div class="otc-info">
          <div class="otc-title">
            {{ order.type === 'sell' ? 'Sell ' + order.ethAmount + ' MON' : 'Buy ' + order.ethAmount + ' MON' }}
          </div>
          <div class="otc-meta">
            <span><i class="far fa-clock"></i> {{ order.time }}</span>
            <span v-if="order.agentId"><i class="fas fa-user"></i> {{ order.agentId }}</span>
          </div>
        </div>
        <div class="otc-price">
          <div class="price-amount">{{ order.usdcAmount.toLocaleString() }} USDC</div>
          <div class="price-per-unit">Unit Price: {{ order.pricePerEth.toLocaleString() }} MON/USDC</div>
        </div>
        <div class="otc-actions">
          <button class="action-btn analyze" @click="showOrderAnalysis(order)">
            <i class="fas fa-chart-bar"></i> Analyze
          </button>
        </div>
      </div>
    </div>
    
    <!-- Empty state -->
    <div v-if="filteredOrders.length === 0" class="empty-state" id="empty-otc">
      <i class="fas fa-box-open"></i>
      <h3>No OTC Orders</h3>
      <p>Waiting for Agent to create first liquidity...</p>
    </div>
    
    <!-- 订单分析弹窗 -->
    <div v-if="showAnalysisModal" class="modal-overlay" @click="closeAnalysisModal">
      <div class="analysis-modal" @click.stop>
        <div class="modal-header">
          <h3>Order Analysis</h3>
          <button class="close-btn" @click="closeAnalysisModal">
            <i class="fas fa-times"></i>
          </button>
        </div>
        <div class="modal-body">
          <div class="analysis-info">
            <div class="analysis-item">
              <span class="label">Order ID:</span>
              <span class="value">{{ selectedOrder?.id }}</span>
            </div>
            <div class="analysis-item" v-if="selectedOrder?.agentId">
              <span class="label">Creator:</span>
              <span class="value">
                <i class="fas fa-user"></i> {{ selectedOrder?.agentId }}
              </span>
            </div>
            <div class="analysis-item">
              <span class="label">Type:</span>
              <span class="value">{{ selectedOrder?.type === 'buy' ? 'Buy MON' : 'Sell MON' }}</span>
            </div>
            <div class="analysis-item">
              <span class="label">Amount:</span>
              <span class="value">{{ selectedOrder?.ethAmount }} MON</span>
            </div>
            <div class="analysis-item">
              <span class="label">Price:</span>
              <span class="value">{{ selectedOrder?.pricePerEth }} MON/USDC</span>
            </div>
          </div>
        </div>
        <div class="modal-footer">
          <button class="btn-primary" @click="closeAnalysisModal">Close</button>
        </div>
      </div>
    </div>
  </section>
</template>

<script setup>
import { defineProps, defineEmits, computed, ref } from 'vue';

const props = defineProps({
  orders: {
    type: Array,
    required: true
  },
  currentFilter: {
    type: String,
    default: 'all'
  },
  filters: {
    type: Array,
    default: () => [
      { value: 'all', label: 'All Orders' },
      { value: 'buy', label: 'Buy MON' },
      { value: 'sell', label: 'Sell MON' }
    ]
  }
});

const emit = defineEmits(['filter-change']);

// 状态管理
const showAnalysisModal = ref(false);
const selectedOrder = ref(null);

const filteredOrders = computed(() => {
  if (props.currentFilter === 'all') {
    return props.orders;
  }
  return props.orders.filter(order => order.type === props.currentFilter);
});


const updateFilter = (filterValue) => {
  emit('filter-change', filterValue);
};

// 显示订单分析弹窗
const showOrderAnalysis = (order) => {
  selectedOrder.value = order;
  showAnalysisModal.value = true;
};

// 关闭订单分析弹窗
const closeAnalysisModal = () => {
  showAnalysisModal.value = false;
  selectedOrder.value = null;
};
</script>

<style scoped>
.otc-marketplace {
  background-color: var(--bg-secondary);
  border-radius: 12px;
  padding: 25px;
  box-shadow: 0 4px 6px rgba(0, 0, 0, 0.1);
  height: 500px;
  width: 100%;
  margin-top: 20px;
}

.otc-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 10px;
}

.otc-filters {
  display: flex;
  gap: 10px;
}

.filter-btn {
  padding: 8px 16px;
  background-color: var(--bg-card);
  border: none;
  border-radius: 6px;
  color: var(--text-secondary);
  cursor: pointer;
  font-size: 0.9rem;
}

.filter-btn.active {
  background-color: var(--accent-blue);
  color: white;
}

.otc-list {
  padding: 10px 0;
  display: flex;
  flex-direction: column;
  gap: 15px;
  max-height: 400px;
  overflow-y: auto;
}

.otc-item {
  background-color: rgba(0, 0, 0, 0.2);
  border-radius: 10px;
  padding: 20px;
  display: flex;
  justify-content: space-between;
  align-items: center;
  border-left: 5px solid var(--accent-green);
  transition: all 0.3s;
}

.otc-item.sell {
  border-left-color: var(--accent-red);
}

.otc-item:hover {
  transform: translateY(-3px);
  box-shadow: 0 6px 12px rgba(0, 0, 0, 0.2);
}

.otc-info {
  flex: 1;
}

.otc-title {
  font-size: 1.1rem;
  font-weight: 600;
  margin-bottom: 5px;
  display: flex;
  align-items: center;
  gap: 10px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.otc-description {
  color: var(--text-secondary);
  font-size: 0.9rem;
  margin-bottom: 10px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.otc-meta {
  display: flex;
  gap: 20px;
  font-size: 0.85rem;
  color: var(--text-secondary);
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.otc-price {
  text-align: right;
  min-width: 150px;
}

.price-amount {
  font-size: 1.5rem;
  font-weight: 700;
  margin-bottom: 5px;
}

.price-per-unit {
  font-size: 0.9rem;
  color: var(--text-secondary);
}

.otc-actions {
  display: flex;
  gap: 10px;
  margin-left: 20px;
}

.action-btn {
  padding: 8px 16px;
  border: none;
  border-radius: 6px;
  font-weight: 600;
  cursor: pointer;
  display: flex;
  align-items: center;
  gap: 5px;
  transition: all 0.2s;
}

.action-btn.buy {
  background-color: var(--accent-green);
  color: white;
}

.action-btn.analyze {
  background-color: var(--bg-card);
  color: var(--text-primary);
}

.action-btn:hover {
  opacity: 0.9;
  transform: scale(1.05);
}

.empty-state {
  text-align: center;
  padding: 40px 20px;
  color: var(--text-secondary);
}

.empty-state i {
  font-size: 3rem;
  margin-bottom: 15px;
  opacity: 0.5;
}

.new-item {
  position: relative;
}

.new-item::before {
  content: "NEW";
  position: absolute;
  top: 10px;
  right: 10px;
  background-color: var(--accent-yellow);
  color: var(--bg-primary);
  font-size: 0.7rem;
  font-weight: 700;
  padding: 3px 8px;
  border-radius: 4px;
}
/* 弹窗样式 */
.modal-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background-color: rgba(0, 0, 0, 0.7);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1000;
}

.analysis-modal {
  background-color: var(--bg-secondary);
  border-radius: 12px;
  box-shadow: 0 10px 25px rgba(0, 0, 0, 0.3);
  width: 90%;
  max-width: 500px;
  max-height: 80vh;
  overflow-y: auto;
  animation: modalFadeIn 0.3s ease;
}

.modal-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 20px;
  border-bottom: 1px solid var(--border-color);
}

.modal-header h3 {
  margin: 0;
  font-size: 1.2rem;
  font-weight: 600;
}

.close-btn {
  background: none;
  border: none;
  color: var(--text-secondary);
  font-size: 1.2rem;
  cursor: pointer;
  padding: 5px;
  border-radius: 4px;
  transition: all 0.2s;
}

.close-btn:hover {
  background-color: rgba(255, 255, 255, 0.1);
  color: var(--text-primary);
}

.modal-body {
  padding: 20px;
}

.analysis-info {
  display: flex;
  flex-direction: column;
  gap: 15px;
}

.analysis-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 10px 0;
  border-bottom: 1px solid rgba(255, 255, 255, 0.1);
}

.analysis-item:last-child {
  border-bottom: none;
}

.analysis-item .label {
  color: var(--text-secondary);
  font-size: 0.9rem;
}

.analysis-item .value {
  font-weight: 600;
  font-size: 0.95rem;
  display: flex;
  align-items: center;
  gap: 5px;
}

.analysis-item .value.negative {
  color: var(--accent-red);
}

.modal-footer {
  padding: 20px;
  border-top: 1px solid var(--border-color);
  display: flex;
  justify-content: flex-end;
  gap: 10px;
}

.btn-primary {
  padding: 10px 20px;
  background-color: var(--accent-blue);
  color: white;
  border: none;
  border-radius: 6px;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.2s;
}

.btn-primary:hover {
  transform: translateY(-2px);
}

@keyframes modalFadeIn {
  from {
    opacity: 0;
    transform: translateY(-20px) scale(0.95);
  }
  to {
    opacity: 1;
    transform: translateY(0) scale(1);
  }
}

/* 响应式调整 */
@media (max-width: 768px) {
  .analysis-modal {
    width: 95%;
    margin: 20px;
  }
  
  .analysis-item {
    flex-direction: column;
    align-items: flex-start;
    gap: 5px;
  }
  
  .analysis-item .value {
    margin-left: 0;
  }
}
</style>