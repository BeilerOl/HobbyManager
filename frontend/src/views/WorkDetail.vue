<template>
  <div class="work-detail">
    <div v-if="loading" class="loading-state">
      <i class="material-icons spin">autorenew</i>
      <span>Chargement…</span>
    </div>
    <div v-else-if="error" class="qtm-alert qtm-alert-error">
      <i class="material-icons">error_outline</i>
      <span>{{ error }}</span>
    </div>
    <template v-else-if="work">
      <div class="detail-header">
        <router-link to="/" class="back-link">
          <i class="material-icons">arrow_back</i>
          <span>Liste</span>
        </router-link>
      </div>

      <div class="detail-card qtm-card">
        <div class="card-header">
          <div class="card-title-row">
            <h1 class="card-title">{{ work.title }}</h1>
            <span
              class="qtm-tag"
              :class="work.seen ? 'qtm-tag-success' : 'qtm-tag-neutral'"
            >
              <i class="material-icons" style="font-size: 0.875rem">{{ work.seen ? 'check_circle' : 'schedule' }}</i>
              {{ work.seen ? 'Vu' : 'À voir' }}
            </span>
          </div>
        </div>

        <div class="card-body">
          <dl class="meta-grid">
            <div class="meta-item">
              <dt>
                <i class="material-icons">category</i>
                Type
              </dt>
              <dd>{{ typeLabel(work.type) }}</dd>
            </div>
            <div class="meta-item">
              <dt>
                <i class="material-icons">person</i>
                Auteur(s)
              </dt>
              <dd>{{ work.authors?.join(', ') || '—' }}</dd>
            </div>
            <div class="meta-item">
              <dt>
                <i class="material-icons">calendar_today</i>
                Date d'ajout
              </dt>
              <dd>{{ formatDate(work.added_at) }}</dd>
            </div>
            <div class="meta-item">
              <dt>
                <i class="material-icons">source</i>
                Origine
              </dt>
              <dd>{{ work.origin || '—' }}</dd>
            </div>
            <div class="meta-item">
              <dt>
                <i class="material-icons">storefront</i>
                Disponibilité
              </dt>
              <dd>{{ work.availability || '—' }}</dd>
            </div>
          </dl>
        </div>

        <div class="card-footer">
          <router-link :to="`/works/${work.id}/edit`" class="qtm-btn qtm-btn-primary">
            <i class="material-icons">edit</i>
            <span>Modifier</span>
          </router-link>
          <button type="button" class="qtm-btn qtm-btn-danger" @click="confirmDelete">
            <i class="material-icons">delete</i>
            <span>Supprimer</span>
          </button>
        </div>
      </div>

      <div v-if="showConfirm" class="confirm-overlay">
        <div class="confirm-dialog qtm-card">
          <div class="confirm-icon">
            <i class="material-icons">warning_amber</i>
          </div>
          <p class="confirm-text">Supprimer « {{ work.title }} » ?</p>
          <p class="confirm-subtext">Cette action est irréversible.</p>
          <div class="confirm-actions">
            <button type="button" class="qtm-btn qtm-btn-danger" @click="doDelete">
              <i class="material-icons">delete</i>
              Oui, supprimer
            </button>
            <button type="button" class="qtm-btn qtm-btn-outline" @click="showConfirm = false">
              Annuler
            </button>
          </div>
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
.detail-header {
  margin-bottom: var(--qtm-space-xl);
}
.back-link {
  display: inline-flex;
  align-items: center;
  gap: var(--qtm-space-xs);
  color: var(--qtm-text-secondary);
  text-decoration: none;
  font-size: var(--qtm-font-size-base);
  font-weight: 500;
  padding: var(--qtm-space-xs) var(--qtm-space-s);
  border-radius: var(--qtm-radius-m);
  transition: all var(--qtm-transition-fast);
  margin-left: calc(-1 * var(--qtm-space-s));
}
.back-link:hover {
  color: var(--qtm-primary-400);
  background: var(--qtm-primary-100);
  text-decoration: none;
}
.back-link .material-icons {
  font-size: 1.125rem;
}

.detail-card {
  overflow: hidden;
}
.card-header {
  padding: var(--qtm-space-xl) var(--qtm-space-xxl);
  border-bottom: 1px solid var(--qtm-border-default);
  background: var(--qtm-bluegrey-50);
}
.card-title-row {
  display: flex;
  align-items: center;
  gap: var(--qtm-space-l);
}
.card-title {
  flex: 1;
  margin: 0;
  font-size: var(--qtm-font-size-xl);
  font-weight: 700;
  color: var(--qtm-text-primary);
}
.card-body {
  padding: var(--qtm-space-xxl);
}
.meta-grid {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: var(--qtm-space-xl);
  margin: 0;
}
.meta-item dt {
  display: flex;
  align-items: center;
  gap: var(--qtm-space-xs);
  color: var(--qtm-text-secondary);
  font-size: var(--qtm-font-size-sm);
  font-weight: 500;
  text-transform: uppercase;
  letter-spacing: 0.03em;
  margin-bottom: var(--qtm-space-xs);
}
.meta-item dt .material-icons {
  font-size: 1rem;
  color: var(--qtm-bluegrey-400);
}
.meta-item dd {
  margin: 0;
  font-size: var(--qtm-font-size-md);
  color: var(--qtm-text-primary);
}

.card-footer {
  display: flex;
  gap: var(--qtm-space-s);
  padding: var(--qtm-space-xl) var(--qtm-space-xxl);
  border-top: 1px solid var(--qtm-border-default);
  background: var(--qtm-bluegrey-50);
}

.confirm-overlay {
  position: fixed;
  inset: 0;
  background: var(--qtm-bg-overlay);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 200;
}
.confirm-dialog {
  padding: var(--qtm-space-xxl);
  max-width: 24rem;
  width: 90%;
  text-align: center;
}
.confirm-icon {
  margin-bottom: var(--qtm-space-l);
}
.confirm-icon .material-icons {
  font-size: 2.5rem;
  color: var(--qtm-warning-400);
}
.confirm-text {
  margin: 0 0 var(--qtm-space-xs);
  font-size: var(--qtm-font-size-md);
  font-weight: 600;
  color: var(--qtm-text-primary);
}
.confirm-subtext {
  margin: 0 0 var(--qtm-space-xl);
  font-size: var(--qtm-font-size-base);
  color: var(--qtm-text-secondary);
}
.confirm-actions {
  display: flex;
  gap: var(--qtm-space-s);
  justify-content: center;
}

.loading-state {
  display: flex;
  align-items: center;
  gap: var(--qtm-space-s);
  padding: var(--qtm-space-xxl);
  justify-content: center;
  color: var(--qtm-text-secondary);
}
@keyframes spin {
  from { transform: rotate(0deg); }
  to { transform: rotate(360deg); }
}
.spin {
  animation: spin 1s linear infinite;
}

@media (max-width: 640px) {
  .meta-grid {
    grid-template-columns: 1fr;
  }
}
</style>
