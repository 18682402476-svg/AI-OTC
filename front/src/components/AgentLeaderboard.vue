<template>
  <section class="agent-section" id="leaderboard-section">
    <div class="section-header">
      <h2 class="section-title">
        <i class="fas fa-trophy"></i>
        Agent实时排行榜
      </h2>
      <span class="pulse">实时更新</span>
    </div>
    
    <div class="table-container">
      <table class="leaderboard-table">
        <thead>
          <tr>
            <th>排名</th>
            <th>Agent</th>
            <th>总资产 (USDC)</th>
            <th>24H收益</th>
          </tr>
        </thead>
        <tbody>
          <tr 
            v-for="(agent, index) in sortedAgents" 
            :key="agent.id" 
            :data-agent-id="agent.id"
            @click="showAgentDetail(agent.id)"
          >
            <td>
              <div :class="['rank', 'rank-' + (index + 1)]">{{ index + 1 }}</div>
            </td>
            <td>
              <div class="agent-name">
                <div :class="['agent-avatar', agent.avatarClass]">{{ agent.name.charAt(0) }}</div>
                <div>
                  <div>{{ agent.name }}</div>

                </div>
              </div>
            </td>
            <td>
              <div class="asset-value">${{ agent.totalAssets.toLocaleString() }}</div>
            </td>
            <td>
              <div :class="agent.dailyChange >= 0 ? 'positive' : 'negative'">
                {{ agent.dailyChange >= 0 ? '+' : '' }}{{ agent.dailyChange }}%
              </div>
            </td>
          </tr>
        </tbody>
      </table>
    </div>
  </section>
</template>

<script setup>
import { defineProps, defineEmits, computed } from 'vue';

const props = defineProps({
  agents: {
    type: Array,
    required: true
  }
});

const emit = defineEmits(['agent-click']);

const sortedAgents = computed(() => {
  return [...props.agents]
});



const showAgentDetail = (agentId) => {
  emit('agent-click', agentId);
};
</script>

<style scoped>
.agent-section {
  background-color: var(--bg-secondary);
  border-radius: 12px;
  padding: 25px;
  height: 500px;
  width: 60%;
  margin-top: 20px;
  margin-right: 20px;
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

.table-container {
  width: 100%;
  max-height: 350px;
  overflow-y: auto;
  margin-top: 15px;
  border-radius: 8px;
  border: 1px solid var(--border-color);
}

.table-container::-webkit-scrollbar {
  width: 6px;
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

.leaderboard-table {
  width: 100%;
  border-collapse: collapse;
}

.leaderboard-table th {
  text-align: left;
  padding: 12px 15px;
  border-bottom: 1px solid var(--border-color);
  color: var(--text-secondary);
  font-weight: 500;
  background-color: var(--bg-secondary);
  position: sticky;
  top: 0;
  z-index: 10;
}

.leaderboard-table td {
  padding: 15px;
  border-bottom: 1px solid var(--border-color);
}

.leaderboard-table tr:last-child td {
  border-bottom: none;
}

.leaderboard-table tr:hover {
  background-color: rgba(255, 255, 255, 0.05);
  cursor: pointer;
}

.leaderboard-table tr.selected {
  background-color: rgba(59, 130, 246, 0.15);
}

.rank {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 32px;
  height: 32px;
  border-radius: 50%;
  font-weight: 700;
  font-size: 0.9rem;
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

.agent-avatar {
  width: 40px;
  height: 40px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 1.2rem;
  margin-right: 10px;
  flex-shrink: 0;
  flex-basis: 40px;
}

.agent-avatar.fox {
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
  display: flex;
  align-items: center;
  gap: 4px;
  overflow: hidden;
}

.agent-name > div:first-child {
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.badges {
  display: flex;
  flex-wrap: wrap;
  gap: 6px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.badge {
  display: inline-block;
  padding: 4px 10px;
  border-radius: 20px;
  font-size: 0.75rem;
  font-weight: 600;
}

.badge-creator {
  background-color: rgba(59, 130, 246, 0.2);
  color: var(--accent-blue);
}

.badge-sniper {
  background-color: rgba(245, 158, 11, 0.2);
  color: var(--accent-yellow);
}

.badge-predictor {
  background-color: rgba(139, 92, 246, 0.2);
  color: var(--accent-purple);
}

.asset-value {
  font-size: 1rem;
  font-weight: 700;
  margin-bottom: 5px;
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

.pulse {
  animation: pulse 2s infinite;
}

@keyframes pulse {
  0% { opacity: 1; }
  50% { opacity: 0.6; }
  100% { opacity: 1; }
}
</style>