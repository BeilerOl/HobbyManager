<template>
  <div class="work-list page-shell">
    <header class="page-header">
      <div>
        <h1 class="page-title">Catalogue des œuvres</h1>
        <p class="page-subtitle">Gérez vos livres, mangas, BD, comics et films dans un tableau unifié.</p>
      </div>
    </header>

    <div class="filters">
      <label class="filter">
        <span>Type</span>
        <select v-model="filterType" class="control-select" @change="load">
          <option value="">Tous les types</option>
          <option v-for="t in WORK_TYPES" :key="t.value" :value="t.value">{{ t.label }}</option>
        </select>
      </label>
      <label class="filter">
        <span>Statut</span>
        <select v-model="filterSeen" class="control-select" @change="load">
          <option value="">Tous</option>
          <option value="false">A voir</option>
          <option value="true">Vu</option>
        </select>
      </label>
    </div>

    <p v-if="loading">Chargement…</p>
    <p v-else-if="error" class="inline-error">{{ error }}</p>
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
            <th>Actions</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="w in works" :key="w.id">
            <td>{{ w.id }}</td>
            <!-- Type: select dropdown -->
            <td>
              <template v-if="editingCell?.id === w.id && editingCell?.field === 'type'">
                <div class="edit-cell">
                  <select v-model="editValue" class="edit-select">
                    <option v-for="t in WORK_TYPES" :key="t.value" :value="t.value">{{ t.label }}</option>
                  </select>
                  <div class="edit-actions">
                    <button type="button" class="btn-icon save" title="Sauvegarder" @click="saveEdit(w)">✓</button>
                    <button type="button" class="btn-icon cancel" title="Annuler" @click="cancelEdit">✕</button>
                  </div>
                </div>
              </template>
              <template v-else>
                <span class="editable text-link" @click="startEdit(w, 'type', w.type)">{{ typeLabel(w.type) }}</span>
              </template>
            </td>
            <!-- Title: text input -->
            <td>
              <template v-if="editingCell?.id === w.id && editingCell?.field === 'title'">
                <div class="edit-cell">
                  <input v-model="editValue" type="text" class="edit-input" @keyup.enter="saveEdit(w)" @keyup.escape="cancelEdit" />
                  <div class="edit-actions">
                    <button type="button" class="btn-icon save" title="Sauvegarder" @click="saveEdit(w)">✓</button>
                    <button type="button" class="btn-icon cancel" title="Annuler" @click="cancelEdit">✕</button>
                  </div>
                </div>
              </template>
              <template v-else>
                <router-link :to="`/works/${w.id}`" class="title-link">{{ w.title }}</router-link>
                <button type="button" class="btn-edit" title="Modifier" @click.stop="startEdit(w, 'title', w.title)">Edit</button>
              </template>
            </td>
            <!-- Authors: text input (comma-separated) -->
            <td>
              <template v-if="editingCell?.id === w.id && editingCell?.field === 'authors'">
                <div class="edit-cell">
                  <input v-model="editValue" type="text" class="edit-input" placeholder="Auteur1, Auteur2" @keyup.enter="saveEdit(w)" @keyup.escape="cancelEdit" />
                  <div class="edit-actions">
                    <button type="button" class="btn-icon save" title="Sauvegarder" @click="saveEdit(w)">✓</button>
                    <button type="button" class="btn-icon cancel" title="Annuler" @click="cancelEdit">✕</button>
                  </div>
                </div>
              </template>
              <template v-else>
                <span class="editable" @click="startEdit(w, 'authors', w.authors?.join(', ') || '')">{{ w.authors?.join(', ') || '—' }}</span>
              </template>
            </td>
            <!-- Date: read-only per API spec -->
            <td>{{ formatDate(w.added_at) }}</td>
            <!-- Origin: text input -->
            <td>
              <template v-if="editingCell?.id === w.id && editingCell?.field === 'origin'">
                <div class="edit-cell">
                  <input v-model="editValue" type="text" class="edit-input" @keyup.enter="saveEdit(w)" @keyup.escape="cancelEdit" />
                  <div class="edit-actions">
                    <button type="button" class="btn-icon save" title="Sauvegarder" @click="saveEdit(w)">✓</button>
                    <button type="button" class="btn-icon cancel" title="Annuler" @click="cancelEdit">✕</button>
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
                  <input v-model="editValue" type="text" class="edit-input" @keyup.enter="saveEdit(w)" @keyup.escape="cancelEdit" />
                  <div class="edit-actions">
                    <button type="button" class="btn-icon save" title="Sauvegarder" @click="saveEdit(w)">✓</button>
                    <button type="button" class="btn-icon cancel" title="Annuler" @click="cancelEdit">✕</button>
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
                  <select v-model="editValue" class="edit-select">
                    <option :value="true">Vu</option>
                    <option :value="false">À voir</option>
                  </select>
                  <div class="edit-actions">
                    <button type="button" class="btn-icon save" title="Sauvegarder" @click="saveEdit(w)">✓</button>
                    <button type="button" class="btn-icon cancel" title="Annuler" @click="cancelEdit">✕</button>
                  </div>
                </div>
              </template>
              <template v-else>
                <span class="status-badge editable" :class="{ seen: w.seen }" @click="startEdit(w, 'seen', w.seen)">{{ w.seen ? 'Vu' : 'A voir' }}</span>
              </template>
            </td>
            <!-- Actions column -->
            <td>
              <router-link :to="`/works/${w.id}`" class="btn-action" title="Voir détails">Voir</router-link>
              <router-link :to="`/works/${w.id}/edit`" class="btn-action" title="Modifier tout">Editer</router-link>
            </td>
          </tr>
        </tbody>
      </table>
    </div>
    <p v-if="saveError" class="inline-error save-error">{{ saveError }}</p>
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
.filters {
  display: flex;
  flex-wrap: wrap;
  gap: 0.9rem;
  margin-bottom: 1.2rem;
}

.filter {
  display: grid;
  gap: 0.4rem;
  min-width: 13rem;
}

.filter span {
  font-size: 0.76rem;
  text-transform: uppercase;
  letter-spacing: 0.08em;
  color: var(--q-text-muted);
}

.table-wrapper {
  overflow-x: auto;
  border-radius: var(--q-radius-md);
  border: 1px solid var(--q-border);
  box-shadow: inset 0 1px 0 rgba(153, 191, 255, 0.08);
}

.works-table {
  width: 100%;
  border-collapse: collapse;
  background: var(--q-surface-2);
  min-width: 980px;
}

.works-table th,
.works-table td {
  padding: 0.72rem 0.8rem;
  text-align: left;
  border-bottom: 1px solid rgba(113, 148, 201, 0.18);
  vertical-align: middle;
}

.works-table th {
  color: var(--q-text-muted);
  font-weight: 600;
  font-size: 0.74rem;
  text-transform: uppercase;
  letter-spacing: 0.08em;
  background: rgba(12, 26, 48, 0.95);
}

.works-table tbody tr:last-child td {
  border-bottom: none;
}

.works-table tbody tr:hover {
  background: rgba(31, 61, 102, 0.24);
}

.title-link {
  color: #9ec8ff;
  text-decoration: none;
  font-weight: 700;
}

.title-link:hover {
  color: #c6e6ff;
}

.empty {
  color: var(--q-text-secondary);
}

.save-error {
  margin-top: 1rem;
}

.editable {
  cursor: pointer;
  border-radius: 999px;
  transition: background 0.15s ease, color 0.15s ease;
}

.editable:hover {
  background: rgba(55, 87, 134, 0.3);
}

.text-link {
  padding: 0.2rem 0.45rem;
  color: #afd2ff;
}

.btn-edit {
  margin-left: 0.5rem;
  padding: 0.16rem 0.52rem;
  background: rgba(25, 49, 80, 0.65);
  border: 1px solid var(--q-border);
  color: #b4d7ff;
  cursor: pointer;
  font-size: 0.78rem;
  border-radius: 999px;
  transition: background 0.15s ease, border-color 0.15s ease;
}

.btn-edit:hover {
  background: rgba(31, 65, 108, 0.88);
  border-color: var(--q-border-strong);
}

.edit-cell {
  display: flex;
  align-items: center;
  gap: 0.5rem;
}
.edit-input,
.edit-select {
  padding: 0.35rem 0.5rem;
  background: rgba(5, 12, 23, 0.9);
  border: 1px solid var(--q-border-strong);
  border-radius: var(--q-radius-sm);
  color: var(--q-text-primary);
  font-size: 0.9rem;
  min-width: 100px;
}
.edit-input:focus,
.edit-select:focus {
  outline: none;
  box-shadow: 0 0 0 3px rgba(47, 143, 255, 0.22);
}
.edit-actions {
  display: flex;
  gap: 0.25rem;
}
.btn-icon {
  padding: 0.24rem 0.46rem;
  border: 1px solid transparent;
  border-radius: var(--q-radius-sm);
  cursor: pointer;
  font-size: 0.85rem;
  transition: background 0.15s ease, border-color 0.15s ease;
}
.btn-icon.save {
  background: rgba(24, 116, 78, 0.66);
  color: #adf7d3;
  border-color: rgba(77, 230, 157, 0.4);
}
.btn-icon.save:hover {
  background: rgba(26, 138, 91, 0.76);
}
.btn-icon.cancel {
  background: rgba(49, 67, 96, 0.7);
  color: var(--q-text-secondary);
  border-color: var(--q-border);
}
.btn-icon.cancel:hover {
  background: rgba(62, 83, 118, 0.82);
}

.btn-action {
  display: inline-block;
  padding: 0.26rem 0.56rem;
  margin-right: 0.35rem;
  text-decoration: none;
  font-size: 0.78rem;
  letter-spacing: 0.04em;
  text-transform: uppercase;
  border-radius: 999px;
  border: 1px solid var(--q-border);
  color: var(--q-text-secondary);
  transition: border-color 0.15s ease, color 0.15s ease, background 0.15s ease;
}
.btn-action:hover {
  color: var(--q-text-primary);
  border-color: var(--q-border-strong);
  background: rgba(30, 63, 109, 0.45);
}
</style>
