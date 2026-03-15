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
                <span class="editable" @click="startEdit(w, 'type', w.type)">{{ typeLabel(w.type) }}</span>
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
                <button type="button" class="btn-edit" title="Modifier" @click.stop="startEdit(w, 'title', w.title)">✎</button>
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
                <span class="badge editable" :class="{ seen: w.seen }" @click="startEdit(w, 'seen', w.seen)">{{ w.seen ? 'Vu' : 'À voir' }}</span>
              </template>
            </td>
            <!-- Actions column -->
            <td>
              <router-link :to="`/works/${w.id}`" class="btn-action" title="Voir détails">👁</router-link>
              <router-link :to="`/works/${w.id}/edit`" class="btn-action" title="Modifier tout">✎</router-link>
            </td>
          </tr>
        </tbody>
      </table>
    </div>
    <p v-if="saveError" class="error save-error">{{ saveError }}</p>
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
  vertical-align: middle;
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
.save-error {
  margin-top: 1rem;
}

.editable {
  cursor: pointer;
  padding: 0.25rem 0.5rem;
  border-radius: 4px;
  transition: background 0.15s;
}
.editable:hover {
  background: #3f3f46;
}

.btn-edit {
  margin-left: 0.5rem;
  padding: 0.15rem 0.4rem;
  background: #3f3f46;
  border: 1px solid #52525b;
  color: #93c5fd;
  cursor: pointer;
  font-size: 0.85rem;
  border-radius: 4px;
  transition: background 0.15s, border-color 0.15s;
}
.btn-edit:hover {
  background: #52525b;
  border-color: #93c5fd;
}

.edit-cell {
  display: flex;
  align-items: center;
  gap: 0.5rem;
}
.edit-input,
.edit-select {
  padding: 0.35rem 0.5rem;
  background: #1a1b1e;
  border: 1px solid #3b82f6;
  border-radius: 4px;
  color: #e4e4e7;
  font-size: 0.9rem;
  min-width: 100px;
}
.edit-input:focus,
.edit-select:focus {
  outline: none;
  border-color: #60a5fa;
}
.edit-actions {
  display: flex;
  gap: 0.25rem;
}
.btn-icon {
  padding: 0.25rem 0.5rem;
  border: none;
  border-radius: 4px;
  cursor: pointer;
  font-size: 0.85rem;
  transition: background 0.15s;
}
.btn-icon.save {
  background: #166534;
  color: #bbf7d0;
}
.btn-icon.save:hover {
  background: #15803d;
}
.btn-icon.cancel {
  background: #3f3f46;
  color: #a1a1aa;
}
.btn-icon.cancel:hover {
  background: #52525b;
}

.btn-action {
  padding: 0.25rem 0.5rem;
  margin-right: 0.25rem;
  text-decoration: none;
  font-size: 1rem;
  border-radius: 4px;
  transition: background 0.15s;
}
.btn-action:hover {
  background: #3f3f46;
}
</style>
