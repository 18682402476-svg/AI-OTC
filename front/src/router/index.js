// 路由配置
import { createRouter, createWebHashHistory } from 'vue-router'

// 导入页面组件
import LLMAgentOTCSystem from '../views/LLMAgentOTCSystem.vue'
import AIList from '../views/AIList.vue'
import APIDoc from '../views/APIDoc.vue'
import AgentDetailPage from '../views/AgentDetailPage.vue'

const routes = [
  {
    path: '/',
    name: 'Home',
    component: LLMAgentOTCSystem
  },
  {
    path: '/ai-list',
    name: 'AIList',
    component: AIList
  },
  {
    path: '/api-doc',
    name: 'APIDoc',
    component: APIDoc
  },
  {
    path: '/agent/:id',
    name: 'AgentDetail',
    component: AgentDetailPage,
    props: true
  }
]

const router = createRouter({
  history: createWebHashHistory(),
  routes
})

export default router