<template>
  <div class="work-list">
    <div class="page-header">
      <h1 class="page-title">Œuvres</h1>
    </div>

    <div class="filters">
      <div class="filter-group">
        <label class="qtm-form-label" for="filter-type">Type</label>
        <select id="filter-type" v-model="filterType" class="qtm-form-select" @change="load">
          <option value="">Tous les types</option>
          <option v-for="t in WORK_TYPES" :key="t.value" :value="t.value">{{ t.label }}</option>
        </select>
      </div>
      <div class="filter-group">
        <label class="qtm-form-label" for="filter-seen">Statut</label>
        <select id="filter-seen" v-model="filterSeen" class="qtm-form-select" @change="load">
          <option value="">Tous</option>
          <option value="false">Non vu</option>
          <option value="true">Déjà vu</option>
        </select>
      </div>
    </div>

    <div v-if="loading" class="loading-state">
      <i class="material-icons spin">autorenew</i>
      <span>Chargement…</span>
    </div>
    <div v-else-if="error" class="qtm-alert qtm-alert-error">
      <i class="material-icons">error_outline</i>
      <span>{{ error }}</span>
    </div>
    <div v-else-if="works.length === 0" class="empty-state">
      <i class="material-icons empty-icon">inventory_2</i>
      <p>Aucune œuvre trouvée.</p>
      <router-link to="/works/new" class="qtm-btn qtm-btn-primary">
        <i class="material-icons">add</i>
        <span>Ajouter une œuvre</span>
      </router-link>
    </div>
    <div v-else class="table-wrapper">
      <table class="qtm-table">
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
            <th>Actions</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="w in works" :key="w.id">
            <td class="cell-id">{{ w.id }}</td>
            <!-- Type: select dropdown -->
            <td>
              <template v-if="editingCell?.id === w.id && editingCell?.field === 'type'">
                <div class="edit-cell">
                  <select v-model="editValue" class="qtm-form-select edit-input">
                    <option v-for="t in WORK_TYPES" :key="t.value" :value="t.value">{{ t.label }}</option>
                  </select>
                  <div class="edit-actions">
                    <button type="button" class="qtm-btn qtm-btn-success qtm-btn-icon qtm-btn-sm" title="Sauvegarder" @click="saveEdit(w)">
                      <i class="material-icons">check</i>
                    </button>
                    <button type="button" class="qtm-btn qtm-btn-outline qtm-btn-icon qtm-btn-sm" title="Annuler" @click="cancelEdit">
                      <i class="material-icons">close</i>
                    </button>
                  </div>
                </div>
              </template>
              <template v-else>
                <span class="editable" @click="startEdit(w, 'type', w.type)">{{ typeLabel(w.type) }}</span>
              </template>
            </td>
            <!-- Title: text input -->
            <td>
              <template v-if="editingCell?.id === w.id && editingCell?.field === 'title'">
                <div class="edit-cell">
                  <input v-model="editValue" type="text" class="qtm-form-input edit-input" @keyup.enter="saveEdit(w)" @keyup.escape="cancelEdit" />
                  <div class="edit-actions">
                    <button type="button" class="qtm-btn qtm-btn-success qtm-btn-icon qtm-btn-sm" title="Sauvegarder" @click="saveEdit(w)">
                      <i class="material-icons">check</i>
                    </button>
                    <button type="button" class="qtm-btn qtm-btn-outline qtm-btn-icon qtm-btn-sm" title="Annuler" @click="cancelEdit">
                      <i class="material-icons">close</i>
                    </button>
                  </div>
                </div>
              </template>
              <template v-else>
                <router-link :to="`/works/${w.id}`" class="title-link">{{ w.title }}</router-link>
                <button type="button" class="inline-edit-btn" title="Modifier" @click.stop="startEdit(w, 'title', w.title)">
                  <i class="material-icons">edit</i>
                </button>
              </template>
            </td>
            <!-- Authors: text input (comma-separated) -->
            <td>
              <template v-if="editingCell?.id === w.id && editingCell?.field === 'authors'">
                <div class="edit-cell">
                  <input v-model="editValue" type="text" class="qtm-form-input edit-input" placeholder="Auteur1, Auteur2" @keyup.enter="saveEdit(w)" @keyup.escape="cancelEdit" />
                  <div class="edit-actions">
                    <button type="button" class="qtm-btn qtm-btn-success qtm-btn-icon qtm-btn-sm" title="Sauvegarder" @click="saveEdit(w)">
                      <i class="material-icons">check</i>
                    </button>
                    <button type="button" class="qtm-btn qtm-btn-outline qtm-btn-icon qtm-btn-sm" title="Annuler" @click="cancelEdit">
                      <i class="material-icons">close</i>
                    </button>
                  </div>
                </div>
              </template>
              <template v-else>
                <span class="editable" @click="startEdit(w, 'authors', w.authors?.join(', ') || '')">{{ w.authors?.join(', ') || '—' }}</span>
              </template>
            </td>
            <!-- Date: read-only per API spec -->
            <td class="cell-date">{{ formatDate(w.added_at) }}</td>
            <!-- Origin: text input -->
            <td>
              <template v-if="editingCell?.id === w.id && editingCell?.field === 'origin'">
                <div class="edit-cell">
                  <input v-model="editValue" type="text" class="qtm-form-input edit-input" @keyup.enter="saveEdit(w)" @keyup.escape="cancelEdit" />
                  <div class="edit-actions">
                    <button type="button" class="qtm-btn qtm-btn-success qtm-btn-icon qtm-btn-sm" title="Sauvegarder" @click="saveEdit(w)">
                      <i class="material-icons">check</i>
                    </button>
                    <button type="button" class="qtm-btn qtm-btn-outline qtm-btn-icon qtm-btn-sm" title="Annuler" @click="cancelEdit">
                      <i class="material-icons">close</i>
                    </button>
                  </div>
                </div>
              </template>
              <template v-else>
                <span class="editable" @click="startEdit(w, 'origin', w.origin || '')">{{ w.origin || '—' }}</span>
              </template>
            </td>
            <!-- Availability: text input -->
            <td>
              <template v-if="editingCell?.id === w.id && editingCell?.field === 'availability'">
                <div class="edit-cell">
                  <input v-model="editValue" type="text" class="qtm-form-input edit-input" @keyup.enter="saveEdit(w)" @keyup.escape="cancelEdit" />
                  <div class="edit-actions">
                    <button type="button" class="qtm-btn qtm-btn-success qtm-btn-icon qtm-btn-sm" title="Sauvegarder" @click="saveEdit(w)">
                      <i class="material-icons">check</i>
                    </button>
                    <button type="button" class="qtm-btn qtm-btn-outline qtm-btn-icon qtm-btn-sm" title="Annuler" @click="cancelEdit">
                      <i class="material-icons">close</i>
                    </button>
                  </div>
                </div>
              </template>
              <template v-else>
                <span class="editable" @click="startEdit(w, 'availability', w.availability || '')">{{ w.availability || '—' }}</span>
              </template>
            </td>
            <!-- Seen: toggle select -->
            <td>
              <template v-if="editingCell?.id === w.id && editingCell?.field === 'seen'">
                <div class="edit-cell">
                  <select v-model="editValue" class="qtm-form-select edit-input">
                    <option :value="true">Vu</option>
                    <option :value="false">À voir</option>
                  </select>
                  <div class="edit-actions">
                    <button type="button" class="qtm-btn qtm-btn-success qtm-btn-icon qtm-btn-sm" title="Sauvegarder" @click="saveEdit(w)">
                      <i class="material-icons">check</i>
                    </button>
                    <button type="button" class="qtm-btn qtm-btn-outline qtm-btn-icon qtm-btn-sm" title="Annuler" @click="cancelEdit">
                      <i class="material-icons">close</i>
                    </button>
                  </div>
                </div>
              </template>
              <template v-else>
                <span
                  class="qtm-tag editable"
                  :class="w.seen ? 'qtm-tag-success' : 'qtm-tag-neutral'"
                  @click="startEdit(w, 'seen', w.seen)"
                >
                  <i class="material-icons" style="font-size: 0.875rem">{{ w.seen ? 'check_circle' : 'schedule' }}</i>
                  {{ w.seen ? 'Vu' : 'À voir' }}
                </span>
              </template>
            </td>
            <!-- Actions column -->
            <td>
              <div class="actions-cell">
                <router-link :to="`/works/${w.id}`" class="qtm-btn qtm-btn-ghost qtm-btn-icon qtm-btn-sm" title="Voir détails">
                  <i class="material-icons">visibility</i>
                </router-link>
                <router-link :to="`/works/${w.id}/edit`" class="qtm-btn qtm-btn-ghost qtm-btn-icon qtm-btn-sm" title="Modifier tout">
                  <i class="material-icons">edit</i>
                </router-link>
              </div>
            </td>
          </tr>
        </tbody>
      </table>
    </div>
    <div v-if="saveError" class="qtm-alert qtm-alert-error save-error">
      <i class="material-icons">error_outline</i>
      <span>{{ saveError }}</span>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { listWorks, updateWork, WORK_TYPES } from '../api/works'

const works = ref([])
const loading = ref(true)
const error = ref('')
const filterType = ref('')
const filterSeen = ref('')

const editingCell = ref(null)
const editValue = ref(null)
const saveError = ref('')

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

function startEdit(work, field, value) {
  editingCell.value = { id: work.id, field }
  editValue.value = value
  saveError.value = ''
}

function cancelEdit() {
  editingCell.value = null
  editValue.value = null
}

async function saveEdit(work) {
  if (!editingCell.value) return
  
  const field = editingCell.value.field
  let newValue = editValue.value
  
  if (field === 'authors') {
    newValue = newValue
      .split(',')
      .map(s => s.trim())
      .filter(Boolean)
  }
  
  const payload = {
    type: work.type,
    title: work.title,
    authors: work.authors || [],
    origin: work.origin || '',
    availability: work.availability || '',
    seen: work.seen,
  }
  
  payload[field] = newValue
  
  saveError.value = ''
  try {
    const updated = await updateWork(work.id, payload)
    const idx = works.value.findIndex(w => w.id === work.id)
    if (idx !== -1) {
      works.value[idx] = updated
    }
    cancelEdit()
  } catch (e) {
    saveError.value = e.message || 'Erreur lors de la sauvegarde'
  }
}

onMounted(load)
</script>

<style scoped>
.page-header {
  margin-bottom: var(--qtm-space-xl);
}
.page-title {
  margin: 0;
  font-size: var(--qtm-font-size-xl);
  font-weight: 700;
  color: var(--qtm-text-primary);
}
.filters {
  display: flex;
  gap: var(--qtm-space-l);
  margin-bottom: var(--qtm-space-xl);
}
.filter-group {
  display: flex;
  flex-direction: column;
  min-width: 10rem;
}
.filter-group .qtm-form-select {
  padding: var(--qtm-space-s) var(--qtm-space-m);
}
.table-wrapper {
  overflow-x: auto;
  border-radius: var(--qtm-radius-l);
  box-shadow: var(--qtm-shadow-s);
}
.cell-id {
  color: var(--qtm-text-secondary);
  font-size: var(--qtm-font-size-sm);
  font-variant-numeric: tabular-nums;
}
.cell-date {
  white-space: nowrap;
  color: var(--qtm-text-secondary);
}
.title-link {
  color: var(--qtm-primary-400);
  font-weight: 500;
  text-decoration: none;
}
.title-link:hover {
  color: var(--qtm-primary-500);
  text-decoration: underline;
}
.editable {
  cursor: pointer;
  padding: var(--qtm-space-xxs) var(--qtm-space-xs);
  border-radius: var(--qtm-radius-s);
  transition: background var(--qtm-transition-fast);
}
.editable:hover {
  background: var(--qtm-bluegrey-100);
}
.inline-edit-btn {
  margin-left: var(--qtm-space-xs);
  padding: 2px;
  background: none;
  border: none;
  color: var(--qtm-text-secondary);
  cursor: pointer;
  border-radius: var(--qtm-radius-s);
  transition: all var(--qtm-transition-fast);
  vertical-align: middle;
  opacity: 0.5;
}
.inline-edit-btn .material-icons {
  font-size: 0.875rem;
}
.inline-edit-btn:hover {
  color: var(--qtm-primary-400);
  opacity: 1;
}
.edit-cell {
  display: flex;
  align-items: center;
  gap: var(--qtm-space-s);
}
.edit-input {
  min-width: 100px;
  font-size: var(--qtm-font-size-sm);
}
.edit-actions {
  display: flex;
  gap: var(--qtm-space-xxs);
}
.actions-cell {
  display: flex;
  gap: var(--qtm-space-xxs);
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
.empty-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  padding: var(--qtm-space-xxxl) var(--qtm-space-xl);
  color: var(--qtm-text-secondary);
  text-align: center;
}
.empty-icon {
  font-size: 3rem;
  color: var(--qtm-bluegrey-300);
  margin-bottom: var(--qtm-space-l);
}
.empty-state p {
  margin: 0 0 var(--qtm-space-xl);
  font-size: var(--qtm-font-size-md);
}
.save-error {
  margin-top: var(--qtm-space-l);
}
</style>
