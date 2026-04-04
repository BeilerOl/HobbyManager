<template>
  <div class="work-import">
    <h1>Importer depuis Excel</h1>
    <p class="hint">Fichier <code>.xlsx</code> — la première ligne doit contenir les en-têtes (ex. titre, type, auteurs…).</p>

    <section v-if="step === 1" class="panel">
      <label class="file-label">
        <span class="btn-primary">Choisir un fichier</span>
        <input type="file" accept=".xlsx,application/vnd.openxmlformats-officedocument.spreadsheetml.sheet" class="hidden-input" @change="onFile" />
      </label>
      <p v-if="fileName" class="file-name">{{ fileName }}</p>
      <button type="button" class="btn-primary" :disabled="!file || parsing" @click="runPreview">
        {{ parsing ? 'Analyse…' : 'Analyser le fichier' }}
      </button>
      <p v-if="error" class="error">{{ error }}</p>
    </section>

    <section v-else-if="step === 2" class="panel">
      <div v-if="sheetWarnings.length" class="warnings">
        <strong>Avertissements :</strong>
        <ul>
          <li v-for="(w, i) in sheetWarnings" :key="i">{{ w }}</li>
        </ul>
      </div>
      <p class="summary">
        {{ validCount }} ligne(s) prête(s) à importer sur {{ previewRows.length }}.
        <span v-if="invalidCount > 0" class="warn-text">{{ invalidCount }} ligne(s) avec erreurs seront ignorées.</span>
      </p>
      <div class="table-wrapper">
        <table class="works-table">
          <thead>
            <tr>
              <th>Ligne</th>
              <th>Type</th>
              <th>Titre</th>
              <th>Auteur(s)</th>
              <th>Origine</th>
              <th>Disponibilité</th>
              <th>Vu</th>
              <th>Statut</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="(r, idx) in previewRows" :key="idx" :class="{ invalid: r.errors?.length }">
              <td>{{ r.row_index }}</td>
              <td>{{ typeLabel(r.work.type) }}</td>
              <td>{{ r.work.title || '—' }}</td>
              <td>{{ (r.work.authors || []).join(', ') || '—' }}</td>
              <td>{{ r.work.origin || '—' }}</td>
              <td>{{ r.work.availability || '—' }}</td>
              <td>{{ r.work.seen ? 'Oui' : 'Non' }}</td>
              <td>
                <span v-if="!r.errors?.length" class="ok">OK</span>
                <span v-else class="err">{{ r.errors.join(' ; ') }}</span>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
      <div class="actions">
        <button type="button" class="btn-secondary" :disabled="importing" @click="reset">Autre fichier</button>
        <button type="button" class="btn-primary" :disabled="importing || validCount === 0" @click="runImport">
          {{ importing ? 'Import…' : 'Importer dans la base' }}
        </button>
      </div>
      <p v-if="error" class="error">{{ error }}</p>
    </section>

    <section v-else class="panel done">
      <p class="success-msg">Import terminé : {{ resultCreated }} œuvre(s) créée(s).</p>
      <p v-if="resultFailed.length" class="warn-text">Échecs : {{ resultFailed.length }}</p>
      <ul v-if="resultFailed.length" class="fail-list">
        <li v-for="(f, i) in resultFailed" :key="i">#{{ f.index + 1 }} — {{ f.message }}</li>
      </ul>
      <div class="actions">
        <router-link to="/" class="btn-primary link-btn">Retour à la liste</router-link>
        <button type="button" class="btn-secondary" @click="fullReset">Nouvel import</button>
      </div>
    </section>
  </div>
</template>

<script setup>
import { ref, computed } from 'vue'
import { previewWorkImport, importWorks, WORK_TYPES } from '../api/works'

const step = ref(1)
const file = ref(null)
const fileName = ref('')
const parsing = ref(false)
const importing = ref(false)
const error = ref('')
const previewRows = ref([])
const sheetWarnings = ref([])
const resultCreated = ref(0)
const resultFailed = ref([])

const validCount = computed(() => previewRows.value.filter((r) => !r.errors?.length).length)
const invalidCount = computed(() => previewRows.value.filter((r) => r.errors?.length).length)

function typeLabel(value) {
  return (WORK_TYPES.find((t) => t.value === value)?.label ?? value) || '—'
}

function onFile(e) {
  const f = e.target.files?.[0]
  file.value = f || null
  fileName.value = f?.name || ''
  error.value = ''
}

async function runPreview() {
  if (!file.value) return
  parsing.value = true
  error.value = ''
  try {
    const data = await previewWorkImport(file.value)
    sheetWarnings.value = data.sheet_warnings || []
    previewRows.value = data.rows || []
    step.value = 2
  } catch (e) {
    error.value = e.message || 'Erreur lors de l’analyse'
  } finally {
    parsing.value = false
  }
}

async function runImport() {
  const items = previewRows.value.filter((r) => !r.errors?.length).map((r) => r.work)
  if (!items.length) return
  importing.value = true
  error.value = ''
  try {
    const res = await importWorks(items)
    resultCreated.value = (res.created || []).length
    resultFailed.value = res.failed || []
    step.value = 3
  } catch (e) {
    error.value = e.message || 'Erreur lors de l’import'
  } finally {
    importing.value = false
  }
}

function reset() {
  step.value = 1
  file.value = null
  fileName.value = ''
  previewRows.value = []
  sheetWarnings.value = []
  error.value = ''
}

function fullReset() {
  step.value = 1
  file.value = null
  fileName.value = ''
  previewRows.value = []
  sheetWarnings.value = []
  resultCreated.value = 0
  resultFailed.value = []
  error.value = ''
}
</script>

<style scoped>
.work-import h1 {
  margin-top: 0;
  font-size: 1.5rem;
}
.hint {
  color: #a1a1aa;
  margin-bottom: 1.25rem;
  font-size: 0.95rem;
}
.hint code {
  background: #3f3f46;
  padding: 0.1rem 0.35rem;
  border-radius: 4px;
  font-size: 0.85rem;
}
.panel {
  background: #25262b;
  border: 1px solid #3f3f46;
  border-radius: 8px;
  padding: 1.25rem;
}
.file-label {
  display: inline-block;
  cursor: pointer;
}
.hidden-input {
  position: absolute;
  width: 0;
  height: 0;
  opacity: 0;
}
.file-name {
  margin: 0.75rem 0;
  color: #a1a1aa;
  font-size: 0.9rem;
}
.btn-primary,
.btn-secondary {
  padding: 0.5rem 1rem;
  border-radius: 6px;
  border: 1px solid transparent;
  cursor: pointer;
  font-size: 0.95rem;
}
.btn-primary {
  background: #3b82f6;
  color: #fff;
}
.btn-primary:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}
.btn-secondary {
  background: #3f3f46;
  color: #e4e4e7;
  border-color: #52525b;
}
.warnings {
  background: #422006;
  border: 1px solid #854d0e;
  border-radius: 6px;
  padding: 0.75rem 1rem;
  margin-bottom: 1rem;
  color: #fef3c7;
  font-size: 0.9rem;
}
.warnings ul {
  margin: 0.5rem 0 0 1.25rem;
  padding: 0;
}
.summary {
  margin-bottom: 1rem;
  color: #d4d4d8;
}
.warn-text {
  color: #fbbf24;
}
.table-wrapper {
  overflow-x: auto;
  margin-bottom: 1rem;
}
.works-table {
  width: 100%;
  border-collapse: collapse;
  background: #1a1b1e;
  border: 1px solid #3f3f46;
  border-radius: 8px;
  overflow: hidden;
  font-size: 0.88rem;
}
.works-table th,
.works-table td {
  padding: 0.6rem 0.5rem;
  text-align: left;
  border-bottom: 1px solid #3f3f46;
  vertical-align: top;
}
.works-table th {
  color: #a1a1aa;
  font-weight: 600;
  font-size: 0.75rem;
  text-transform: uppercase;
}
.works-table tbody tr.invalid {
  background: rgba(127, 29, 29, 0.25);
}
.ok {
  color: #86efac;
}
.err {
  color: #fca5a5;
  font-size: 0.8rem;
}
.actions {
  display: flex;
  gap: 0.75rem;
  flex-wrap: wrap;
  align-items: center;
}
.error {
  color: #f87171;
  margin-top: 0.75rem;
}
.done .success-msg {
  color: #86efac;
  font-size: 1.05rem;
}
.fail-list {
  color: #fca5a5;
  margin: 0.5rem 0 1rem 1.25rem;
}
.link-btn {
  display: inline-block;
  text-decoration: none;
  text-align: center;
}
</style>
