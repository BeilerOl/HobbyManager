<template>
  <div class="work-list">
    <h1>Œuvres</h1>
    <div class="filters">
      <select v-model="filterType" @change="load">
        <option value="">Tous les types</option>
        <option v-for="t in WORK_TYPES" :key="t.value" :value="t.value">{{ t.label }}</option>
      </select>
      <select v-model="filterSeen" @change="load">
        <option value="">Tous</option>
        <option value="false">Non vu</option>
        <option value="true">Déjà vu</option>
      </select>
    </div>
    <p v-if="loading">Chargement…</p>
    <p v-else-if="error" class="error">{{ error }}</p>
    <ul v-else-if="works.length === 0" class="empty">Aucune œuvre.</ul>
    <ul v-else class="list">
      <li v-for="w in works" :key="w.id" class="item">
        <router-link :to="`/works/${w.id}`" class="link">
          <span class="title">{{ w.title }}</span>
          <span class="meta">{{ typeLabel(w.type) }} · {{ w.authors?.join(', ') || '—' }}</span>
          <span class="badge" :class="{ seen: w.seen }">{{ w.seen ? 'Vu' : 'À voir' }}</span>
        </router-link>
      </li>
    </ul>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { listWorks, WORK_TYPES } from '../api/works'

const works = ref([])
const loading = ref(true)
const error = ref('')
const filterType = ref('')
const filterSeen = ref('')

function typeLabel(value) {
  return WORK_TYPES.find(t => t.value === value)?.label ?? value
}

async function load() {
  loading.value = true
  error.value = ''
  try {
    const params = {}
    if (filterType.value) params.type = filterType.value
    if (filterSeen.value !== '') params.seen = filterSeen.value === 'true'
    works.value = await listWorks(params)
  } catch (e) {
    error.value = e.message || 'Erreur lors du chargement'
  } finally {
    loading.value = false
  }
}

onMounted(load)
</script>

<style scoped>
.work-list h1 {
  margin-top: 0;
  font-size: 1.5rem;
}
.filters {
  display: flex;
  gap: 0.75rem;
  margin-bottom: 1rem;
}
.filters select {
  padding: 0.5rem 0.75rem;
  background: #25262b;
  border: 1px solid #3f3f46;
  border-radius: 6px;
  color: #e4e4e7;
}
.list {
  list-style: none;
  padding: 0;
  margin: 0;
}
.item {
  margin-bottom: 0.5rem;
}
.link {
  display: flex;
  align-items: center;
  gap: 0.75rem;
  padding: 0.75rem 1rem;
  background: #25262b;
  border-radius: 8px;
  text-decoration: none;
  color: inherit;
  border: 1px solid transparent;
}
.link:hover {
  border-color: #71717a;
}
.title {
  font-weight: 600;
  flex: 1;
}
.meta {
  color: #a1a1aa;
  font-size: 0.9rem;
}
.badge {
  font-size: 0.75rem;
  padding: 0.25rem 0.5rem;
  border-radius: 4px;
  background: #3f3f46;
  color: #a1a1aa;
}
.badge.seen {
  background: #166534;
  color: #bbf7d0;
}
.empty, .error {
  color: #a1a1aa;
}
.error {
  color: #f87171;
}
</style>
