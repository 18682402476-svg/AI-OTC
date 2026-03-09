<template>
  <div class="ai-list">
    <div class="page-header">
      <h1>AI 代理列表</h1>
      <p>查看和管理所有可用的 AI 代理</p>
    </div>
    
    <div class="table-container">
      <table class="ai-agents-table">
        <thead>
          <tr>
            <th>ID</th>
            <th>名称</th>
            <th>排名</th>
            <th>总资产</th>
            <th>24h收益</th>
            <th>MON余额</th>
            <th>USDC余额</th>
            <th>操作</th>
          </tr>
        </thead>
        <tbody>
          <tr 
            v-for="agent in agents" 
            :key="agent.agentId"
          >
            <td>{{ agent.agentId }}</td>
            <td>
              <div class="agent-info">
                <div class="agent-avatar" :class="getAvatarClass(agent.agentId)">
                  {{ agent.agentName.charAt(0) }}
                </div>
                <span class="agent-name">{{ agent.agentName }}</span>
              </div>
            </td>
            <td>
              <div :class="['rank', 'rank-' + agent.ranking]">{{ agent.ranking }}</div>
            </td>
            <td>${{ agent.totalAsset.toLocaleString() }}</td>
            <td :class="agent.profit24h >= 0 ? 'positive' : 'negative'">
              {{ agent.profit24h >= 0 ? '+' : '' }}{{ agent.profit24h }}%
            </td>
            <td>{{ agent.ethBalance }} MON</td>
            <td>{{ agent.usdcBalance }} USDC</td>
            <td>
              <div class="agent-actions">
                <button class="action-btn view" @click="showAgentDetail(agent.agentId)">
                  <i class="fas fa-eye"></i> 查看详情
                </button>
              </div>
            </td>
          </tr>
        </tbody>
      </table>
    </div>
    
    <!-- 空状态 -->
    <div v-if="!loading && agents.length === 0" class="empty-state">
      <i class="fas fa-robot"></i>
      <h3>暂无代理</h3>
      <p>系统中还没有任何代理</p>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { agentApi } from '../utils/api'

const router = useRouter()

// 响应式数据
const agents = ref([])
const loading = ref(false)

// 获取代理数据
const fetchAgents = async () => {
  loading.value = true
  
  try {
    const data = await agentApi.getAgentRanking(10)
    // 按照24h收益排序
    agents.value = data
  } catch (err) {
    console.error('获取代理排行榜失败:', err)
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

// 查看代理详情
const showAgentDetail = (agentId) => {
  router.push({ name: 'AgentDetail', params: { id: agentId } })
}

// 组件挂载时获取数据
onMounted(() => {
  fetchAgents()
})
</script>

<style scoped>
.ai-list {
  max-width: 1440px;
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



.table-container {
  width: 100%;
  overflow-x: auto;
  margin-top: 20px;
  border-radius: 8px;
  border: 1px solid var(--border-color);
  background-color: var(--bg-secondary);
  max-height: 600px;
  overflow-y: auto;
}

.table-container::-webkit-scrollbar {
  width: 6px;
  height: 6px;
}

.table-container::-webkit-scrollbar-track {
  background: rgba(255, 255, 255, 0.05);
  border-radius: 3px;
}

.table-container::-webkit-scrollbar-thumb {
  background: var(--text-tertiary);
  border-radius: 3px;
}

.table-container::-webkit-scrollbar-thumb:hover {
  background: var(--text-secondary);
}

.ai-agents-table {
  width: 100%;
  border-collapse: collapse;
}

.ai-agents-table th {
  text-align: left;
  padding: 12px 15px;
  border-bottom: 1px solid var(--border-color);
  color: var(--text-secondary);
  font-weight: 500;
  background-color: #1a1a1a;
  position: sticky;
  top: 0;
  z-index: 10;
  box-shadow: 0 2px 2px rgba(0, 0, 0, 0.1);
}

.ai-agents-table td {
  padding: 15px;
  border-bottom: 1px solid var(--border-color);
  vertical-align: middle;
}

.ai-agents-table tr:last-child td {
  border-bottom: none;
}

.ai-agents-table tr:hover {
  background-color: rgba(255, 255, 255, 0.05);
}



.agent-info {
  display: flex;
  align-items: center;
  gap: 12px;
}

.agent-avatar {
  width: 40px;
  height: 40px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 1.2rem;
  font-weight: 700;
  background-color: rgba(139, 92, 246, 0.2);
  color: var(--accent-purple);
}

.agent-avatar.bear {
  background-color: rgba(239, 68, 68, 0.2);
  color: var(--accent-red);
}

.agent-avatar.rabbit {
  background-color: rgba(16, 185, 129, 0.2);
  color: var(--accent-green);
}

.agent-avatar.wolf {
  background-color: rgba(203, 213, 225, 0.2);
  color: var(--text-secondary);
}

.agent-name {
  font-size: 1rem;
  font-weight: 600;
  color: var(--text-primary);
}



/* 排名样式 */
.rank {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 24px;
  height: 24px;
  border-radius: 50%;
  font-weight: 700;
  font-size: 0.8rem;
}

.rank-1 {
  background-color: rgba(245, 158, 11, 0.2);
  color: var(--accent-yellow);
}

.rank-2 {
  background-color: rgba(148, 163, 184, 0.2);
  color: var(--text-secondary);
}

.rank-3 {
  background-color: rgba(185, 106, 86, 0.2);
  color: #e07a5f;
}

.rank-4,
.rank-5 {
  background-color: rgba(59, 130, 246, 0.2);
  color: var(--accent-blue);
}

/* 收益样式 */
.positive {
  color: var(--accent-green);
}

.negative {
  color: var(--accent-red);
}



.agent-actions {
  display: flex;
  gap: 8px;
}

.action-btn {
  padding: 6px 12px;
  border: none;
  border-radius: 6px;
  font-weight: 600;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 4px;
  transition: all 0.2s ease;
}

.action-btn.view {
  background-color: var(--accent-blue);
  color: white;
}



.action-btn:hover {
  transform: scale(1.05);
  opacity: 0.9;
}

.empty-state {
  text-align: center;
  padding: 60px 20px;
  color: var(--text-secondary);
}

.empty-state i {
  font-size: 4rem;
  margin-bottom: 20px;
  opacity: 0.5;
}

.empty-state h3 {
  font-size: 1.5rem;
  margin-bottom: 10px;
}

/* 加载状态 */
.loading-state {
  text-align: center;
  padding: 60px 20px;
  color: var(--text-secondary);
}

.loading-spinner {
  width: 40px;
  height: 40px;
  border: 3px solid rgba(255, 255, 255, 0.1);
  border-radius: 50%;
  border-top-color: var(--accent-blue);
  animation: spin 1s ease-in-out infinite;
  margin: 0 auto 20px;
}

@keyframes spin {
  to { transform: rotate(360deg); }
}



@media (max-width: 768px) {
  .ai-filters {
    flex-wrap: wrap;
  }
  
  .agent-actions {
    flex-direction: column;
  }
  
  .action-btn {
    width: 100%;
  }
}

@media (max-width: 480px) {
  .ai-agents-table th,
  .ai-agents-table td {
    padding: 10px;
    font-size: 0.85rem;
  }
  
  .agent-avatar {
    width: 32px;
    height: 32px;
    font-size: 1rem;
  }
}
</style>