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
    <p v-else-if="works.length === 0" class="empty">Aucune œuvre.</p>
    <div v-else class="table-wrapper">
      <table class="works-table">
        <thead>
          <tr>
            <th>ID</th>
            <th>Type</th>
            <th>Titre</th>
            <th>Auteur(s)</th>
            <th>Date d'ajout</th>
            <th>Origine</th>
            <th>Disponibilité</th>
            <th>Statut</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="w in works" :key="w.id">
            <td>{{ w.id }}</td>
            <td>{{ typeLabel(w.type) }}</td>
            <td>
              <router-link :to="`/works/${w.id}`" class="title-link">{{ w.title }}</router-link>
            </td>
            <td>{{ w.authors?.join(', ') || '—' }}</td>
            <td>{{ formatDate(w.added_at) }}</td>
            <td>{{ w.origin || '—' }}</td>
            <td>{{ w.availability || '—' }}</td>
            <td>
              <span class="badge" :class="{ seen: w.seen }">{{ w.seen ? 'Vu' : 'À voir' }}</span>
            </td>
          </tr>
        </tbody>
      </table>
    </div>
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

function formatDate(s) {
  if (!s) return '—'
  try {
    const d = new Date(s)
    return d.toLocaleDateString('fr-FR', { dateStyle: 'medium' })
  } catch {
    return s
  }
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
.table-wrapper {
  overflow-x: auto;
}
.works-table {
  width: 100%;
  border-collapse: collapse;
  background: #25262b;
  border: 1px solid #3f3f46;
  border-radius: 8px;
  overflow: hidden;
}
.works-table th,
.works-table td {
  padding: 0.75rem;
  text-align: left;
  border-bottom: 1px solid #3f3f46;
  vertical-align: top;
}
.works-table th {
  color: #a1a1aa;
  font-weight: 600;
  font-size: 0.85rem;
  text-transform: uppercase;
  letter-spacing: 0.02em;
}
.works-table tbody tr:last-child td {
  border-bottom: none;
}
.title-link {
  color: #93c5fd;
  text-decoration: none;
  font-weight: 600;
}
.title-link:hover {
  text-decoration: underline;
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
.empty,
.error {
  color: #a1a1aa;
}
.error {
  color: #f87171;
}
</style>
