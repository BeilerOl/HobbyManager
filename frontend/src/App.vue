<template>
  <div class="app-shell">
    <header class="app-header">
      <div class="brand-block">
        <p class="brand-kicker">Quantum-inspired workspace</p>
        <router-link to="/" class="logo">HobbyManager</router-link>
      </div>
      <nav class="main-nav">
        <router-link to="/">Catalogue</router-link>
        <router-link to="/works/new">Ajouter</router-link>
      </nav>
    </header>
    <main class="app-content">
      <router-view />
    </main>
  </div>
</template>

<script setup>
</script>

<style scoped>
.app-shell {
  min-height: 100vh;
  display: flex;
  flex-direction: column;
}

.app-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 1rem;
  padding: 1.15rem 1.75rem;
  background: linear-gradient(180deg, rgba(15, 25, 44, 0.92), rgba(8, 14, 28, 0.92));
  border-bottom: 1px solid var(--q-border);
  backdrop-filter: blur(12px);
  position: sticky;
  top: 0;
  z-index: 100;
}

.brand-block {
  display: grid;
  gap: 0.15rem;
}

.brand-kicker {
  margin: 0;
  font-size: 0.72rem;
  text-transform: uppercase;
  letter-spacing: 0.08em;
  color: var(--q-text-muted);
}

.logo {
  font-weight: 800;
  font-size: 1.3rem;
  color: var(--q-text-primary);
  text-decoration: none;
  letter-spacing: 0.02em;
}

.main-nav {
  display: flex;
  align-items: center;
  gap: 0.55rem;
}

.main-nav a {
  border-radius: 999px;
  padding: 0.42rem 0.88rem;
  color: var(--q-text-secondary);
  text-decoration: none;
  border: 1px solid transparent;
  transition: border-color 0.18s ease, color 0.18s ease, background 0.18s ease;
}

.main-nav a:hover {
  border-color: var(--q-border);
  color: var(--q-text-primary);
}

.main-nav a.router-link-active {
  color: var(--q-text-primary);
  border-color: rgba(66, 217, 255, 0.55);
  background: rgba(27, 81, 126, 0.36);
}

.app-content {
  flex: 1;
  padding: 1.6rem;
  max-width: 72rem;
  margin: 0 auto;
  width: 100%;
}
</style>
