<template>
  <div class="work-detail page-shell">
    <p v-if="loading">Chargement…</p>
    <p v-else-if="error" class="inline-error">{{ error }}</p>
    <template v-else-if="work">
      <header class="page-header work-header">
        <div>
          <router-link to="/" class="back">← Retour au catalogue</router-link>
          <h1 class="page-title">{{ work.title }}</h1>
          <p class="page-subtitle">Fiche détaillée de l'œuvre sélectionnée.</p>
        </div>
        <span class="status-badge" :class="{ seen: work.seen }">{{ work.seen ? 'Vu' : 'A voir' }}</span>
      </header>

      <dl class="meta-grid">
        <div class="meta-item">
          <dt>Type</dt>
          <dd>{{ typeLabel(work.type) }}</dd>
        </div>
        <div class="meta-item">
          <dt>Auteur(s)</dt>
          <dd>{{ work.authors?.join(', ') || '—' }}</dd>
        </div>
        <div class="meta-item">
          <dt>Date d'ajout</dt>
          <dd>{{ formatDate(work.added_at) }}</dd>
        </div>
        <div class="meta-item">
          <dt>Origine</dt>
          <dd>{{ work.origin || '—' }}</dd>
        </div>
        <div class="meta-item">
          <dt>Disponibilité</dt>
          <dd>{{ work.availability || '—' }}</dd>
        </div>
      </dl>

      <div class="actions">
        <router-link :to="`/works/${work.id}/edit`" class="action-btn primary">Modifier</router-link>
        <button type="button" class="action-btn danger" @click="confirmDelete">Supprimer</button>
      </div>

      <div v-if="showConfirm" class="confirm-card">
        <p>Supprimer « {{ work.title }} » ?</p>
        <div class="confirm-actions">
          <button type="button" class="action-btn danger" @click="doDelete">Oui, supprimer</button>
          <button type="button" class="action-btn" @click="showConfirm = false">Annuler</button>
        </div>
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
.work-header {
  margin-bottom: 1.1rem;
}

.back {
  color: var(--q-text-secondary);
  text-decoration: none;
  display: inline-flex;
  margin-bottom: 0.45rem;
  font-size: 0.88rem;
  letter-spacing: 0.04em;
  text-transform: uppercase;
}

.back:hover {
  color: var(--q-text-primary);
}

.meta-grid {
  margin: 0 0 1.5rem;
  padding: 0;
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(180px, 1fr));
  gap: 0.8rem;
}

.meta-item {
  margin: 0;
  border-radius: var(--q-radius-md);
  border: 1px solid var(--q-border);
  background: rgba(8, 18, 34, 0.75);
  padding: 0.72rem 0.8rem;
}

.meta-item dt {
  margin: 0;
  color: var(--q-text-muted);
  text-transform: uppercase;
  letter-spacing: 0.07em;
  font-size: 0.72rem;
}

.meta-item dd {
  margin: 0.38rem 0 0;
  color: var(--q-text-primary);
}

.actions {
  display: flex;
  flex-wrap: wrap;
  gap: 0.6rem;
}

.confirm-card {
  margin-top: 1rem;
  padding: 1rem 1.05rem;
  border-radius: var(--q-radius-md);
  border: 1px solid rgba(255, 95, 143, 0.33);
  background: rgba(41, 16, 29, 0.5);
}

.confirm-card p {
  margin-top: 0;
  margin-bottom: 0.9rem;
}

.confirm-actions {
  display: flex;
  flex-wrap: wrap;
  gap: 0.55rem;
}
</style>
