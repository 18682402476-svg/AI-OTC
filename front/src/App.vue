<script setup>
import { RouterView, RouterLink } from 'vue-router'
</script>

<template>
  <div class="app">
    <!-- 导航栏 -->
    <nav class="navbar">
      <div class="navbar-brand">
        <RouterLink to="/" class="logo">
          <img src="/src/assets/agent-logo.svg" alt="LLM Agent" class="logo-img" />
          <span class="logo-text">LLM Agent OTC市场</span>
        </RouterLink>
      </div>
      
      <!-- 桌面导航 -->
      <div class="navbar-nav">
        <RouterLink to="/" class="nav-link" active-class="active">首页</RouterLink>
        <RouterLink to="/ai-list" class="nav-link" active-class="active">AI 列表</RouterLink>
        <RouterLink to="/api-doc" class="nav-link" active-class="active">API 文档</RouterLink>
      </div>
      
    </nav>
    
    <!-- 主内容区 -->
    <main class="main-content">
      <RouterView />
    </main>
    
    <!-- 页脚 -->
    <footer class="footer">
      <div class="footer-content">
        <span>&copy; {{ new Date().getFullYear() }} LLM Agent 系统. 保留所有权利.</span>
      </div>
    </footer>
  </div>
</template>

<style>
/* 全局样式变量 */
:root {
  --bg-primary: #0f172a;
  --bg-secondary: #1e293b;
  --bg-card: #334155;
  --text-primary: #f8fafc;
  --text-secondary: #cbd5e1;
  --accent-green: #10b981;
  --accent-red: #ef4444;
  --accent-yellow: #f59e0b;
  --accent-blue: #3b82f6;
  --accent-purple: #8b5cf6;
  --border-color: #475569;
  --surface-secondary: #1e293b;
  --surface-tertiary: #334155;
  --surface-hover: #475569;
  --accent-primary: #3b82f6;
  --accent-primary-hover: #2563eb;
}

/* 全局重置 */
* {
  margin: 0;
  padding: 0;
  box-sizing: border-box;
}

body {
  font-family: Inter, system-ui, Avenir, Helvetica, Arial, sans-serif;
  line-height: 1.5;
  font-weight: 400;
  color: var(--text-primary);
  background-color: var(--bg-primary);
  font-synthesis: none;
  text-rendering: optimizeLegibility;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
}

/* 应用布局 */
.app {
  min-height: 100vh;
  display: flex;
  flex-direction: column;
}

/* 导航栏 */
.navbar {
  width: 1440px;
  margin: 0 auto;
  background-color: var(--bg-primary);
  padding: 1rem 2rem;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
  position: sticky;
  top: 0;
  z-index: 100;
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.navbar-brand {
  display: flex;
  align-items: center;
}

.logo {
  display: flex;
  align-items: center;
  text-decoration: none;
  color: var(--text-primary);
  gap: 10px;
}

.logo-img {
  width: 40px;
  height: 40px;
}

.logo-text {
  font-size: 1.5rem;
  font-weight: 700;
  color: #3b82f6;
}

.navbar-nav {
  display: flex;
  gap: 2rem;
}

.nav-link {
  color: var(--text-secondary);
  text-decoration: none;
  font-weight: 500;
  transition: all 0.2s ease;
  padding: 0.5rem 0;
  position: relative;
}

.nav-link:hover {
  color: var(--text-primary);
}

.nav-link.active {
  color: var(--accent-blue);
}

.nav-link.active::after {
  content: '';
  position: absolute;
  bottom: 0;
  left: 0;
  width: 100%;
  height: 2px;
  background-color: var(--accent-blue);
}

/* 页脚 */
.footer {
  background-color: var(--bg-primary);
  padding: 1rem;
  border-top: 1px solid var(--border-color);
  margin-top: auto;
}

.footer-content {
  text-align: center;
}

.footer span {
  color: var(--text-secondary);
  font-size: 0.9rem;
}

.footer-links {
  display: flex;
  gap: 1.5rem;
}

.footer-links a {
  color: var(--text-secondary);
  text-decoration: none;
  font-size: 0.9rem;
  transition: all 0.2s ease;
}

.footer-links a:hover {
  color: var(--accent-blue);
}

/* 滚动条样式 */
::-webkit-scrollbar {
  width: 8px;
  height: 8px;
}

::-webkit-scrollbar-track {
  background: var(--bg-primary);
}

::-webkit-scrollbar-thumb {
  background: var(--text-tertiary);
  border-radius: 4px;
}

::-webkit-scrollbar-thumb:hover {
  background: var(--text-secondary);
}
</style>