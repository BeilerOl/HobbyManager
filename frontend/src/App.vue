<template>
  <div class="app">
    <header class="qtm-header">
      <div class="qtm-header-brand">
        <router-link to="/" class="logo">
          <i class="material-icons logo-icon">home</i>
          <span class="logo-text">HobbyManager</span>
        </router-link>
      </div>
      <div class="qtm-header-divider"></div>
      <nav class="qtm-header-nav">
        <router-link to="/" class="nav-link" active-class="nav-link-active">
          <i class="material-icons">list</i>
          <span>Liste</span>
        </router-link>
        <router-link to="/works/new" class="nav-link" active-class="nav-link-active">
          <i class="material-icons">add_circle_outline</i>
          <span>Ajouter</span>
        </router-link>
      </nav>
      <div class="qtm-header-actions">
        <button class="qtm-btn qtm-btn-icon theme-toggle" @click="toggleTheme" :title="isDark ? 'Switch to light mode' : 'Switch to dark mode'">
          <i class="material-icons">{{ isDark ? 'light_mode' : 'dark_mode' }}</i>
        </button>
      </div>
    </header>
    <main class="main">
      <router-view />
    </main>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'

const isDark = ref(false)

function toggleTheme() {
  isDark.value = !isDark.value
  const theme = isDark.value ? 'dark' : 'light'
  document.documentElement.setAttribute('data-theme', theme)
  localStorage.setItem('theme', theme)
}

onMounted(() => {
  const saved = localStorage.getItem('theme') || 'light'
  isDark.value = saved === 'dark'
  document.documentElement.setAttribute('data-theme', saved)
})
</script>

<style>
.app {
  min-height: 100vh;
  display: flex;
  flex-direction: column;
}
.qtm-header {
  display: flex;
  align-items: center;
  padding: 0 var(--qtm-space-xl);
  background: var(--qtm-bg-elevated);
  border-bottom: 1px solid var(--qtm-border-default);
  position: sticky;
  top: 0;
  z-index: 100;
  height: 3.5rem;
  box-shadow: var(--qtm-shadow-s);
}
.qtm-header-brand {
  display: flex;
  align-items: center;
}
.logo {
  display: flex;
  align-items: center;
  gap: var(--qtm-space-s);
  text-decoration: none;
  color: var(--qtm-primary-500);
}
.logo:hover {
  text-decoration: none;
}
.logo-icon {
  font-size: 1.5rem;
  color: var(--qtm-primary-400);
}
.logo-text {
  font-weight: 700;
  font-size: var(--qtm-font-size-lg);
  color: var(--qtm-primary-500);
}
.qtm-header-divider {
  width: 1px;
  height: 1.75rem;
  background: var(--qtm-border-default);
  margin: 0 var(--qtm-space-xl);
}
.qtm-header-nav {
  display: flex;
  align-items: center;
  gap: var(--qtm-space-xs);
  height: 100%;
}
.nav-link {
  display: flex;
  align-items: center;
  gap: var(--qtm-space-xs);
  padding: var(--qtm-space-s) var(--qtm-space-m);
  color: var(--qtm-text-secondary);
  text-decoration: none;
  font-size: var(--qtm-font-size-base);
  font-weight: 500;
  border-radius: var(--qtm-radius-m);
  transition: all var(--qtm-transition-fast);
  height: 100%;
  border-bottom: 2px solid transparent;
  border-radius: 0;
}
.nav-link:hover {
  color: var(--qtm-primary-400);
  background: var(--qtm-primary-100);
  text-decoration: none;
}
.nav-link-active {
  color: var(--qtm-primary-400);
  border-bottom-color: var(--qtm-primary-400);
}
.nav-link .material-icons {
  font-size: 1.125rem;
}
.qtm-header-actions {
  margin-left: auto;
  display: flex;
  align-items: center;
}
.theme-toggle {
  color: var(--qtm-text-secondary);
  background: transparent;
  border-color: transparent;
}
.theme-toggle:hover {
  color: var(--qtm-primary-400);
  background: var(--qtm-primary-100);
}
.main {
  flex: 1;
  padding: var(--qtm-space-xxl) var(--qtm-space-xl);
  max-width: 72rem;
  margin: 0 auto;
  width: 100%;
}
</style>
