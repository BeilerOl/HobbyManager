<template>
  <div class="work-detail">
    <p v-if="loading">Chargement…</p>
    <p v-else-if="error" class="error">{{ error }}</p>
    <template v-else-if="work">
      <div class="header">
        <router-link to="/" class="back">← Liste</router-link>
        <h1>{{ work.title }}</h1>
        <span class="badge" :class="{ seen: work.seen }">{{ work.seen ? 'Vu' : 'À voir' }}</span>
      </div>
      <dl class="meta">
        <dt>Type</dt>
        <dd>{{ typeLabel(work.type) }}</dd>
        <dt>Auteur(s)</dt>
        <dd>{{ work.authors?.join(', ') || '—' }}</dd>
        <dt>Date d'ajout</dt>
        <dd>{{ formatDate(work.added_at) }}</dd>
        <dt>Origine</dt>
        <dd>{{ work.origin || '—' }}</dd>
        <dt>Disponibilité</dt>
        <dd>{{ work.availability || '—' }}</dd>
      </dl>
      <div class="actions">
        <router-link :to="`/works/${work.id}/edit`" class="btn primary">Modifier</router-link>
        <button type="button" class="btn danger" @click="confirmDelete">Supprimer</button>
      </div>
      <div v-if="showConfirm" class="confirm">
        <p>Supprimer « {{ work.title }} » ?</p>
        <button type="button" class="btn danger" @click="doDelete">Oui, supprimer</button>
        <button type="button" class="btn" @click="showConfirm = false">Annuler</button>
      </div>
    </template>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { getWork, deleteWork, WORK_TYPES } from '../api/works'

const route = useRoute()
const router = useRouter()
const work = ref(null)
const loading = ref(true)
const error = ref('')
const showConfirm = ref(false)

const id = computed(() => Number(route.params.id))

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
    work.value = await getWork(id.value)
    if (!work.value) error.value = 'Œuvre introuvable.'
  } catch (e) {
    error.value = e.message || 'Erreur'
  } finally {
    loading.value = false
  }
}

function confirmDelete() {
  showConfirm.value = true
}

async function doDelete() {
  try {
    await deleteWork(id.value)
    router.push('/')
  } catch (e) {
    error.value = e.message || 'Erreur lors de la suppression'
  }
  showConfirm.value = false
}

onMounted(load)
</script>

<style scoped>
.work-detail .header {
  display: flex;
  align-items: center;
  gap: 1rem;
  margin-bottom: 1.5rem;
}
.back {
  color: #a1a1aa;
  text-decoration: none;
}
.back:hover {
  color: #fafafa;
}
.work-detail h1 {
  flex: 1;
  margin: 0;
  font-size: 1.5rem;
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
.meta {
  display: grid;
  grid-template-columns: auto 1fr;
  gap: 0.5rem 1.5rem;
  margin-bottom: 1.5rem;
}
.meta dt {
  color: #a1a1aa;
  margin: 0;
}
.meta dd {
  margin: 0;
}
.actions {
  display: flex;
  gap: 0.5rem;
}
.btn {
  padding: 0.5rem 1rem;
  border-radius: 6px;
  border: 1px solid #3f3f46;
  background: #25262b;
  color: #e4e4e7;
  cursor: pointer;
  text-decoration: none;
  font-size: 0.9rem;
}
.btn:hover {
  background: #3f3f46;
}
.btn.primary {
  border-color: #3b82f6;
  background: #2563eb;
  color: #fff;
}
.btn.primary:hover {
  background: #1d4ed8;
}
.btn.danger {
  border-color: #b91c1c;
  background: #dc2626;
  color: #fff;
}
.btn.danger:hover {
  background: #b91c1c;
}
.confirm {
  margin-top: 1rem;
  padding: 1rem;
  background: #25262b;
  border-radius: 8px;
  border: 1px solid #3f3f46;
}
.confirm p {
  margin-top: 0;
  margin-bottom: 0.75rem;
}
.confirm .btn {
  margin-right: 0.5rem;
}
.error {
  color: #f87171;
}
</style>
