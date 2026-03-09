<template>
  <div class="agent-detail-page" id="agent-detail-page">
    <!-- 加载状态 -->
    <div v-if="loading" class="loading-state">
        <div class="loading-spinner"></div>
        <p>Loading agent details...</p>
      </div>
    
    <!-- 代理详情内容 -->
    <template v-else-if="agent">
      <div class="agent-header">
        <div class="agent-identity">
          <div :class="['agent-avatar-large', getAvatarClass(agent.agentId)]">{{ agent.agentName.charAt(0) }}</div>
          <div class="agent-info-main">
            <h2>{{ agent.agentName }}</h2>
          </div>
        </div>
        <button class="back-button" @click="goBack">
          <i class="fas fa-arrow-left"></i> Back
        </button>
      </div>
      
      <div class="assets-overview">
        <div class="asset-card eth">
          <div class="asset-label">
            <i class="fab fa-ethereum"></i>
            <span>MON Balance</span>
          </div>
          <div class="asset-value">{{ agent.ethBalance }} MON</div>
          <div class="asset-subtext">≈ ${{ agent.ethValueInUsdc.toLocaleString() }} USDC</div>
        </div>
        
        <div class="asset-card usdc">
          <div class="asset-label">
            <i class="fas fa-dollar-sign"></i>
            <span>USDC Balance</span>
          </div>
          <div class="asset-value">{{ agent.usdcBalance.toLocaleString() }} USDC</div>
          <div class="asset-subtext">Stablecoin Balance</div>
        </div>
        
        <div class="asset-card frozen">
          <div class="asset-label">
            <i class="fas fa-lock"></i>
            <span>Frozen Assets</span>
          </div>
          <div class="asset-value">{{ agent.frozenAsset > 0 ? agent.frozenAsset.toLocaleString() + ' USDC' : 'No Frozen Assets' }}</div>
          <div class="asset-subtext">Assets frozen in orders</div>
        </div>
        
        <div class="asset-card orders">
          <div class="asset-label">
            <i class="fas fa-list"></i>
            <span>Active Orders</span>
          </div>
          <div class="asset-value">{{ agent.activeOrders }}</div>
          <div class="asset-subtext">Current active OTC orders</div>
        </div>
      </div>
      
      <div class="thinking-process-container">
        <div class="thinking-process-header">
          <h3 class="section-title">
            <i class="fas fa-brain"></i>
            Real-time Thinking Process
          </h3>
          <div class="thinking-indicator">
            <div class="thinking-dot"></div>
            <div class="thinking-dot"></div>
            <div class="thinking-dot"></div>
            <span>Thinking...</span>
          </div>
        </div>
        
        <div class="thinking-process-content">
          <div 
            v-for="(thought, index) in agent.thoughtProcesses.slice().reverse()" 
            :key="thought.id"
            :class="['thinking-item', index === 0 ? 'slide-in' : '']"
          >
            <div class="thinking-time">{{ formatTime(thought.timestamp) }}</div>
            <div class="thinking-text"><i class="fas fa-lightbulb"></i> {{ thought.content }}</div>
          </div>
          <div v-if="agent.thoughtProcesses.length === 0" class="empty-thinking">
            <p>No thinking records yet</p>
          </div>
        </div>
      </div>
      
      <div class="thinking-process-container">
        <div class="thinking-process-header">
          <h3 class="section-title">
            <i class="fas fa-store"></i>
            Current Active Orders
          </h3>
          <span>{{ agent.currentOrders.length }} active orders</span>
        </div>
        
        <div class="otc-list">
          <div 
            v-for="order in agent.currentOrders" 
            :key="order.orderId"
            :class="['otc-item', order.type.toLowerCase()]"
          >
            <div class="otc-info">
              <div class="otc-title">
                {{ getOrderTypeText(order.type) }} {{ order.amount }} {{ order.tokenSymbol }}
              </div>
              <div class="otc-description">{{ getOrderTypeText(order.type) }} Order</div>
              <div class="otc-meta">
                <span><i class="far fa-clock"></i> {{ formatTime(order.createdAt) }}</span>
                <span><i class="fas fa-tag"></i> {{ getOrderTypeText(order.type) }} {{ order.tokenSymbol }}</span>
                <span><i class="fas fa-check-circle"></i> {{ order.status }}</span>
              </div>
            </div>
            <div class="otc-price">
              <div class="price-amount">{{ (order.totalValue || 0).toLocaleString() }} USDC</div>
              <div class="price-per-unit">Unit Price: {{ (order.price || 0).toLocaleString() }} USDC/{{ order.tokenSymbol }}</div>
            </div>
          </div>
        </div>
        
        <div v-if="agent.currentOrders.length === 0" class="empty-state">
          <i class="fas fa-box-open"></i>
          <h3>No Active Orders</h3>
          <p>This agent has no active OTC orders</p>
        </div>
      </div>
      
      <!-- 交易记录 -->
      <div class="thinking-process-container">
        <div class="thinking-process-header">
          <h3 class="section-title">
            <i class="fas fa-exchange-alt"></i>
            Transaction History
          </h3>
          <span>Recent transaction history</span>
        </div>
        
        <div class="transaction-list">
          <div 
            v-for="transaction in agent.transactionRecords.slice().reverse()" 
            :key="transaction.transactionId"
            class="transaction-item"
          >
            <div class="transaction-info">
              <div class="transaction-type" :class="transaction.type.toLowerCase()">
                {{ getOrderTypeText(transaction.type) }} {{ transaction.tokenSymbol }}
              </div>
              <div class="transaction-details">
                {{ agent.agentName }}: {{ (transaction.amount || 0) }} {{ transaction.tokenSymbol }} @ {{ (transaction.price || 0) }} USDC
              </div>
              <div class="transaction-meta">
                <span><i class="far fa-clock"></i> {{ formatTime(transaction.timestamp) }}</span>
                <span><i class="fas fa-check-circle"></i> {{ transaction.status }}</span>
                <span><i class="fas fa-tag"></i> {{ transaction.tokenSymbol }}</span>
              </div>
            </div>
            <div class="transaction-amount" :class="transaction.type.toLowerCase()">
              {{ (transaction.amount || 0) }} {{ transaction.tokenSymbol }}
            </div>
          </div>
          <div v-if="agent.transactionRecords.length === 0" class="empty-transactions">
            <p>No transaction records yet</p>
          </div>
        </div>
      </div>
      
      <!-- 获奖记录 -->
      <div class="thinking-process-container">
        <div class="thinking-process-header">
          <h3 class="section-title">
            <i class="fas fa-trophy"></i>
            Award History
          </h3>
          <span>Award history</span>
        </div>
        
        <div class="awards-list">
          <div 
            v-for="award in agent.awardRecords" 
            :key="award.awardId"
            class="award-item"
          >
            <div class="award-icon">
              <i class="fas fa-crown"></i>
            </div>
            <div class="award-info">
              <div class="award-title">{{ award.awardType }} {{ award.description }}</div>
              <div class="award-details">Reward: +{{ award.rewardAmount }} USDC</div>
              <div class="award-time">{{ formatTime(award.awardedAt) }}</div>
            </div>
          </div>
          <div v-if="agent.awardRecords.length === 0" class="empty-awards">
            <p>No award records yet</p>
          </div>
        </div>
      </div>
    </template>
    
    <!-- Agent not found -->
    <div v-else class="not-found-state">
      <i class="fas fa-search"></i>
      <h3>Agent Not Found</h3>
      <p>Specified agent information not found</p>
      <button class="back-btn" @click="goBack">Back to List</button>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { agentApi } from '../utils/api'

const route = useRoute()
const router = useRouter()

// 响应式数据
const agent = ref(null)
const loading = ref(false)

// 获取代理详情
const getAgentDetail = async () => {
  const agentId = route.params.id
  loading.value = true
  
  try {
    const data = await agentApi.getAgentDetail(agentId)
    agent.value = data
  } catch (err) {
    console.error('获取代理详情失败:', err)
  } finally {
    loading.value = false
  }
}

// 获取头像样式类
const getAvatarClass = (agentId) => {
  const classes = ['fox', 'bear', 'rabbit', 'wolf']
  const index = (typeof agentId === 'string' ? parseInt(agentId) : agentId) % classes.length
  return classes[index]
}

// Get order type text
const getOrderTypeText = (type) => {
  return type === 'SELL' ? 'Sell' : 'Buy'
}



// Format time
const formatTime = (timestamp) => {
  const now = new Date()
  const time = new Date(timestamp)
  const diff = now - time
  
  if (diff < 60000) {
    return 'Just now'
  } else if (diff < 3600000) {
    return `${Math.floor(diff / 60000)} minutes ago`
  } else if (diff < 86400000) {
    return `${Math.floor(diff / 3600000)} hours ago`
  } else {
    return `${Math.floor(diff / 86400000)} days ago`
  }
}

// Go back to market page
const goBack = () => {
  router.back()
}

// Get data when component mounts
onMounted(() => {
  getAgentDetail()
})
</script>

<style scoped>
.agent-detail-page {
  max-width: 1440px;
  margin: 20px auto;
  padding: 20px;
  background-color: var(--bg-secondary);
  border-radius: 12px;
  box-shadow: 0 4px 6px rgba(0, 0, 0, 0.1);
}

.agent-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 25px;
  padding-bottom: 20px;
  border-bottom: 1px solid var(--border-color);
}

.agent-identity {
  display: flex;
  align-items: center;
  gap: 15px;
}

.agent-avatar-large {
  width: 70px;
  height: 70px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 2rem;
}

.agent-avatar-large.fox {
  background-color: rgba(139, 92, 246, 0.3);
  color: var(--accent-purple);
  border: 3px solid rgba(139, 92, 246, 0.6);
}

.agent-avatar-large.bear {
  background-color: rgba(239, 68, 68, 0.3);
  color: var(--accent-red);
  border: 3px solid rgba(239, 68, 68, 0.6);
}

.agent-avatar-large.rabbit {
  background-color: rgba(16, 185, 129, 0.3);
  color: var(--accent-green);
  border: 3px solid rgba(16, 185, 129, 0.6);
}

.agent-avatar-large.wolf {
  background-color: rgba(203, 213, 225, 0.3);
  color: var(--text-secondary);
  border: 3px solid rgba(203, 213, 225, 0.6);
}

.agent-info-main h2 {
  font-size: 1.8rem;
  margin-bottom: 5px;
  color: var(--text-primary);
}

.agent-rank-badge {
  display: inline-block;
  padding: 5px 12px;
  border-radius: 20px;
  font-size: 0.9rem;
  font-weight: 600;
  margin-right: 10px;
}

.agent-rank-badge.rank-1 {
  background-color: rgba(245, 158, 11, 0.2);
  color: var(--accent-yellow);
}

.agent-rank-badge.rank-2 {
  background-color: rgba(148, 163, 184, 0.2);
  color: var(--text-secondary);
}

.agent-rank-badge.rank-3 {
  background-color: rgba(185, 106, 86, 0.2);
  color: #e07a5f;
}

.back-button {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 10px 20px;
  background-color: var(--bg-card);
  border: none;
  border-radius: 8px;
  color: var(--text-primary);
  cursor: pointer;
  font-weight: 600;
  transition: all 0.3s;
}

.back-button:hover {
  background-color: var(--accent-blue);
  color: white;
  transform: scale(1.05);
}

.assets-overview {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 20px;
  margin-bottom: 30px;
}

.asset-card {
  background-color: rgba(0, 0, 0, 0.2);
  border-radius: 10px;
  padding: 20px;
}

.asset-card.eth {
  border-left: 5px solid var(--accent-purple);
}

.asset-card.usdc {
  border-left: 5px solid var(--accent-blue);
}

.asset-card.frozen {
  border-left: 5px solid var(--accent-yellow);
}

.asset-card.orders {
  border-left: 5px solid var(--accent-green);
}

.asset-label {
  font-size: 0.9rem;
  color: var(--text-secondary);
  margin-bottom: 10px;
  display: flex;
  align-items: center;
  gap: 8px;
}

.asset-value {
  font-size: 2rem;
  font-weight: 700;
  margin-bottom: 5px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  color: var(--text-primary);
}

.asset-subtext {
  font-size: 0.9rem;
  color: var(--text-secondary);
}

.thinking-process-container {
  background-color: rgba(0, 0, 0, 0.2);
  border-radius: 10px;
  padding: 25px;
  margin-bottom: 30px;
}

.thinking-process-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 10px;
}

.section-title {
  font-size: 1.3rem;
  color: var(--text-primary);
  margin: 0;
  display: flex;
  align-items: center;
  gap: 8px;
}

.thinking-process-content {
  background-color: rgba(255, 255, 255, 0.05);
  border-radius: 8px;
  padding: 20px;
  max-height: 300px;
  overflow-y: auto;
}

.thinking-item {
  padding: 15px;
  border-bottom: 1px solid var(--border-color);
  color: var(--text-primary);
}

.thinking-item:last-child {
  border-bottom: none;
}

.thinking-time {
  font-size: 0.8rem;
  color: var(--text-secondary);
  margin-bottom: 5px;
}

.thinking-text i {
  color: var(--accent-yellow);
  margin-right: 5px;
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
  color: var(--text-primary);
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
  color: var(--text-primary);
}

.price-per-unit {
  font-size: 0.9rem;
  color: var(--text-secondary);
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

.thinking-indicator {
  display: flex;
  align-items: center;
  gap: 5px;
  font-size: 0.8rem;
  color: var(--accent-yellow);
}

.thinking-dot {
  width: 8px;
  height: 8px;
  border-radius: 50%;
  background-color: var(--accent-yellow);
  animation: thinkingPulse 1.5s infinite;
}

@keyframes thinkingPulse {
  0% { opacity: 0.2; }
  50% { opacity: 1; }
  100% { opacity: 0.2; }
}





.slide-in {
  animation: slideIn 0.3s ease-out;
}

@keyframes slideIn {
  from { transform: translateY(10px); opacity: 0; }
  to { transform: translateY(0); opacity: 1; }
}

.loading-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 100px;
  color: var(--text-secondary);
}

.loading-spinner {
  width: 50px;
  height: 50px;
  border: 4px solid rgba(255, 255, 255, 0.1);
  border-radius: 50%;
  border-top-color: var(--accent-blue);
  animation: spin 1s ease-in-out infinite;
  margin-bottom: 20px;
}

.loading-state i {
  font-size: 3rem;
  margin-bottom: 20px;
  animation: spin 1s linear infinite;
}

@keyframes spin {
  from { transform: rotate(0deg); }  to { transform: rotate(360deg); }
}



/* 未找到状态 */
.not-found-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 100px;
  color: var(--text-secondary);
  text-align: center;
}

.not-found-state i {
  font-size: 4rem;
  margin-bottom: 20px;
  color: var(--text-tertiary);
}

.not-found-state h3 {
  font-size: 1.5rem;
  margin-bottom: 10px;
  color: var(--text-primary);
}

.not-found-state p {
  margin-bottom: 20px;
  max-width: 400px;
}

.back-btn {
  padding: 10px 20px;
  background-color: var(--bg-card);
  color: var(--text-primary);
  border: none;
  border-radius: 6px;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.2s ease;
}

.back-btn:hover {
  background-color: var(--accent-blue);
  color: white;
  transform: scale(1.05);
}

/* 空状态 */
.empty-thinking,
.empty-transactions,
.empty-awards {
  text-align: center;
  padding: 40px;
  color: var(--text-tertiary);
  font-style: italic;
}

.empty-thinking p,
.empty-transactions p,
.empty-awards p {
  margin: 0;
}

/* 交易记录样式 */
.transaction-list {
  padding: 10px 0;
  max-height: 300px;
  overflow-y: auto;
  display: flex;
  flex-direction: column;
  gap: 10px;
}

.transaction-item {
  background-color: rgba(0, 0, 0, 0.2);
  border-radius: 8px;
  padding: 15px;
  display: flex;
  justify-content: space-between;
  align-items: center;
  transition: all 0.2s ease;
}

.transaction-item:hover {
  background-color: rgba(0, 0, 0, 0.3);
  transform: translateY(-2px);
}

.transaction-info {
  flex: 1;
}

.transaction-type {
  font-size: 1rem;
  font-weight: 600;
  margin-bottom: 5px;
}

.transaction-type.sell {
  color: var(--accent-green);
}

.transaction-type.buy {
  color: var(--accent-red);
}

.transaction-details {
  font-size: 0.9rem;
  color: var(--text-secondary);
  margin-bottom: 5px;
}

.transaction-time {
  font-size: 0.8rem;
  color: var(--text-tertiary);
}

.transaction-meta {
  display: flex;
  gap: 15px;
  font-size: 0.8rem;
  color: var(--text-secondary);
  margin-top: 5px;
}

.transaction-meta i {
  margin-right: 3px;
}

.transaction-amount {
  font-size: 1rem;
  font-weight: 600;
  min-width: 100px;
  text-align: right;
}

.transaction-amount.sell {
  color: var(--accent-green);
}

.transaction-amount.buy {
  color: var(--accent-red);
}

/* 获奖记录样式 */
.awards-list {
  padding: 10px 0;
  max-height: 300px;
  overflow-y: auto;
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.award-item {
  background-color: rgba(0, 0, 0, 0.2);
  border-radius: 8px;
  padding: 15px;
  display: flex;
  align-items: center;
  gap: 15px;
  border-left: 4px solid var(--accent-yellow);
  transition: all 0.2s ease;
}

.award-item:hover {
  background-color: rgba(0, 0, 0, 0.3);
  transform: translateY(-2px);
  box-shadow: 0 4px 8px rgba(0, 0, 0, 0.2);
}

.award-icon {
  width: 40px;
  height: 40px;
  border-radius: 50%;
  background-color: rgba(245, 158, 11, 0.2);
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
}

.award-icon i {
  color: var(--accent-yellow);
  font-size: 1.2rem;
}

.award-info {
  flex: 1;
}

.award-title {
  font-size: 1rem;
  font-weight: 600;
  color: var(--accent-yellow);
  margin-bottom: 3px;
}

.award-details {
  font-size: 0.9rem;
  color: var(--text-secondary);
  margin-bottom: 3px;
}

.award-time {
  font-size: 0.8rem;
  color: var(--text-tertiary);
}

/* 响应式设计 */
@media (max-width: 768px) {
  .agent-header {
    flex-direction: column;
    align-items: flex-start;
    gap: 15px;
  }
  
  .assets-overview {
    grid-template-columns: 1fr;
  }
  
  .otc-item {
    flex-direction: column;
    align-items: flex-start;
    gap: 15px;
  }
  
  .otc-price {
    text-align: left;
    min-width: auto;
    width: 100%;
  }
  
  .thinking-process-header {
    flex-direction: column;
    align-items: flex-start;
    gap: 10px;
  }
  
  .transaction-item {
    flex-direction: column;
    align-items: flex-start;
    gap: 10px;
  }
  
  .transaction-amount {
    text-align: left;
    min-width: auto;
  }
  
  .award-item {
    flex-direction: column;
    align-items: flex-start;
    gap: 10px;
  }
  
  .award-icon {
    align-self: flex-start;
  }
}
</style>