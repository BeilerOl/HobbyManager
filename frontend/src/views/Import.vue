<template>
  <div class="import-page">
    <h1>Import Excel</h1>
    <p class="intro">Choisissez un fichier .xls pour prévisualiser les données, puis lancez l'import après contrôle.</p>

    <div class="upload-zone" @click="fileInput?.click()">
      <input
        ref="fileInput"
        type="file"
        accept=".xls"
        class="file-input"
        @change="onFileChange"
      />
      <p v-if="!file && !loading">Cliquez pour choisir un fichier .xls</p>
      <p v-else-if="file">{{ file.name }}</p>
      <p v-else-if="loading">Chargement…</p>
    </div>

    <p v-if="previewError" class="error">{{ previewError }}</p>

    <template v-if="previewRows.length > 0">
      <div class="table-wrapper">
        <table class="preview-table">
          <thead>
            <tr>
              <th class="col-include">Inclure</th>
              <th>Ligne</th>
              <th>Type</th>
              <th>Titre</th>
              <th>Auteur(s)</th>
              <th>Origine</th>
              <th>Disponibilité</th>
              <th>Vu</th>
              <th>Erreur</th>
            </tr>
          </thead>
          <tbody>
            <tr
              v-for="row in previewRows"
              :key="row.row_index"
              :class="{ 'row-error': row.error, 'row-skip': !row.work && !row.error }"
            >
              <td class="col-include">
                <input
                  v-if="row.work && !row.error"
                  v-model="selectedIndices"
                  type="checkbox"
                  :value="row.row_index"
                />
                <span v-else>—</span>
              </td>
              <td>{{ row.row_index }}</td>
              <td>{{ row.work ? typeLabel(row.work.type) : '—' }}</td>
              <td>{{ row.work?.title ?? '—' }}</td>
              <td>{{ row.work?.authors?.join(', ') ?? '—' }}</td>
              <td>{{ row.work?.origin ?? '—' }}</td>
              <td>{{ row.work?.availability ?? '—' }}</td>
              <td>{{ row.work ? (row.work.seen ? 'Oui' : 'Non') : '—' }}</td>
              <td class="cell-error">{{ row.error ?? '—' }}</td>
            </tr>
          </tbody>
        </table>
      </div>

      <div class="actions">
        <button
          type="button"
          class="btn primary"
          :disabled="executeLoading || selectedIndices.length === 0"
          @click="runImport"
        >
          {{ executeLoading ? 'Import en cours…' : 'Lancer l\'import' }}
        </button>
      </div>

      <div v-if="executeResult" class="result">
        <p class="success">{{ executeResult.created }} œuvre(s) créée(s).</p>
        <p v-if="executeResult.errors?.length" class="errors">
          Erreurs : {{ executeResult.errors.map((e) => `Ligne ${e.row_index}: ${e.message}`).join(' ; ') }}
        </p>
      </div>
    </template>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { importPreview, importExecute } from '../api/import'
import { WORK_TYPES } from '../api/works'

const fileInput = ref(null)
const file = ref(/** @type {File|null} */ (null))
const loading = ref(false)
const previewError = ref('')
const previewRows = ref(/** @type {{ row_index: number, error?: string, work?: import('../api/realApi').WorkCreate }[] } */ ([]))
const selectedIndices = ref(/** @type {number[]} */ ([]))
const executeLoading = ref(false)
const executeResult = ref(/** @type {{ created: number, errors?: { row_index: number, message: string }[] } | null} */ (null))

function typeLabel(type) {
  const t = WORK_TYPES.find((x) => x.value === type)
  return t ? t.label : type
}

function onFileChange(ev) {
  const f = ev.target.files?.[0]
  if (!f) return
  file.value = f
  loadPreview()
}

async function loadPreview() {
  if (!file.value) return
  loading.value = true
  previewError.value = ''
  previewRows.value = []
  selectedIndices.value = []
  executeResult.value = null
  try {
    const res = await importPreview(file.value)
    previewRows.value = res.rows || []
    selectedIndices.value = previewRows.value
      .filter((r) => r.work && !r.error)
      .map((r) => r.row_index)
  } catch (e) {
    previewError.value = e instanceof Error ? e.message : String(e)
  } finally {
    loading.value = false
  }
}

async function runImport() {
  if (!file.value || selectedIndices.value.length === 0) return
  executeLoading.value = true
  executeResult.value = null
  try {
    const res = await importExecute(file.value, selectedIndices.value)
    executeResult.value = res
  } catch (e) {
    previewError.value = e instanceof Error ? e.message : String(e)
  } finally {
    executeLoading.value = false
  }
}
</script>

<style scoped>
.import-page {
  max-width: 100%;
}
.intro {
  color: #a1a1aa;
  margin-bottom: 1rem;
}
.upload-zone {
  position: relative;
  border: 2px dashed #3f3f46;
  border-radius: 8px;
  padding: 2rem;
  text-align: center;
  margin-bottom: 1rem;
  cursor: pointer;
}
.upload-zone .file-input {
  position: absolute;
  width: 0;
  height: 0;
  opacity: 0;
  pointer-events: none;
}
.upload-zone p {
  margin: 0;
  color: #a1a1aa;
}
.table-wrapper {
  overflow-x: auto;
  margin: 1rem 0;
}
.preview-table {
  width: 100%;
  border-collapse: collapse;
  font-size: 0.875rem;
}
.preview-table th,
.preview-table td {
  padding: 0.5rem 0.75rem;
  text-align: left;
  border-bottom: 1px solid #3f3f46;
}
.preview-table th {
  background: #25262b;
  color: #a1a1aa;
  font-weight: 600;
}
.col-include {
  width: 4rem;
}
.cell-error {
  color: #f87171;
  max-width: 12rem;
}
.row-error {
  background: rgba(248, 113, 113, 0.08);
}
.actions {
  margin-top: 1rem;
}
.result {
  margin-top: 1rem;
  padding: 1rem;
  background: #25262b;
  border-radius: 8px;
}
.result .success {
  color: #86efac;
  margin: 0 0 0.5rem 0;
}
.result .errors {
  color: #f87171;
  margin: 0;
  font-size: 0.875rem;
}
.error {
  color: #f87171;
  margin: 0.5rem 0;
}
</style>
