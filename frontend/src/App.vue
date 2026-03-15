<template>
  <div class="app">
    <header class="header">
      <router-link to="/" class="logo">HobbyManager</router-link>
      <nav>
        <router-link to="/">Liste</router-link>
        <router-link to="/works/new">Ajouter</router-link>
        <router-link to="/import">Import</router-link>
      </nav>
    </header>
    <main class="main">
      <router-view />
    </main>
  </div>
</template>

<script setup>
</script>

<style>
* {
  box-sizing: border-box;
}
body {
  margin: 0;
  font-family: system-ui, -apple-system, sans-serif;
  background: #1a1b1e;
  color: #e4e4e7;
}
.app {
  min-height: 100vh;
  display: flex;
  flex-direction: column;
}
.header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 1rem 1.5rem;
  background: #25262b;
  border-bottom: 1px solid #3f3f46;
}
.logo {
  font-weight: 700;
  font-size: 1.25rem;
  color: #fafafa;
  text-decoration: none;
}
.header nav a {
  margin-left: 1rem;
  color: #a1a1aa;
  text-decoration: none;
}
.header nav a.router-link-active {
  color: #fafafa;
}
.main {
  flex: 1;
  padding: 1.5rem;
  max-width: 56rem;
  margin: 0 auto;
  width: 100%;
}
</style>
