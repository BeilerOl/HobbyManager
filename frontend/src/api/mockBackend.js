/**
 * In-memory mock API for demo (e.g. GitHub Pages). Same interface as real API.
 */
import mockWorksData from '../data/mockWorks.json'

/** @typedef {'roman'|'livre_culture_generale'|'film'|'serie_tv'|'jeu_societe'|'jeu_video'} WorkType */
/**
 * @typedef {Object} Work
 * @property {number} id
 * @property {WorkType} type
 * @property {string} title
 * @property {string[]} authors
 * @property {string} added_at
 * @property {string} origin
 * @property {string} availability
 * @property {boolean} seen
 */
/**
 * @typedef {Object} WorkCreate
 * @property {WorkType} type
 * @property {string} title
 * @property {string[]} authors
 * @property {string} origin
 * @property {string} availability
 * @property {boolean} seen
 */

let store = [...mockWorksData]
let nextId = Math.max(0, ...store.map((w) => w.id)) + 1

function clone(work) {
  return { ...work, authors: [...(work.authors || [])] }
}

/**
 * @param {{ type?: string, seen?: boolean }} [params]
 * @returns {Promise<Work[]>}
 */
export async function listWorks(params = {}) {
  await Promise.resolve()
  let list = store.map(clone)
  if (params.type) list = list.filter((w) => w.type === params.type)
  if (params.seen !== undefined) list = list.filter((w) => w.seen === params.seen)
  return list
}

/**
 * @param {number} id
 * @returns {Promise<Work|null>}
 */
export async function getWork(id) {
  await Promise.resolve()
  const work = store.find((w) => w.id === Number(id))
  return work ? clone(work) : null
}

/**
 * @param {WorkCreate} data
 * @returns {Promise<Work>}
 */
export async function createWork(data) {
  await Promise.resolve()
  const added_at = new Date().toISOString()
  const work = {
    id: nextId++,
    type: data.type,
    title: data.title,
    authors: Array.isArray(data.authors) ? data.authors : [],
    added_at,
    origin: data.origin || '',
    availability: data.availability || '',
    seen: Boolean(data.seen),
  }
  store.push(work)
  return clone(work)
}

/**
 * @param {number} id
 * @param {WorkCreate} data
 * @returns {Promise<Work>}
 */
export async function updateWork(id, data) {
  await Promise.resolve()
  const idx = store.findIndex((w) => w.id === Number(id))
  if (idx === -1) throw new Error('Not found')
  const work = {
    ...store[idx],
    type: data.type,
    title: data.title,
    authors: Array.isArray(data.authors) ? data.authors : [],
    origin: data.origin || '',
    availability: data.availability || '',
    seen: Boolean(data.seen),
  }
  store[idx] = work
  return clone(work)
}

/**
 * @param {number} id
 * @returns {Promise<void>}
 */
export async function deleteWork(id) {
  await Promise.resolve()
  const idx = store.findIndex((w) => w.id === Number(id))
  if (idx !== -1) store.splice(idx, 1)
}

export const WORK_TYPES = [
  { value: 'roman', label: 'Roman' },
  { value: 'livre_culture_generale', label: 'Livre de culture générale' },
  { value: 'film', label: 'Film' },
  { value: 'serie_tv', label: 'Série TV' },
  { value: 'jeu_societe', label: 'Jeu de société' },
  { value: 'jeu_video', label: 'Jeu vidéo' },
]

/**
 * Mock: return fake preview rows (from first 3 mock works as WorkCreate).
 * @param {File} [file] - ignored in mock
 * @returns {Promise<{ rows: { row_index: number, error?: string, work?: WorkCreate }[] }>}
 */
export async function importPreview(file) {
  void file
  await Promise.resolve()
  const sample = store.slice(0, 3).map((w) => ({
    type: w.type,
    title: w.title,
    authors: w.authors || [],
    origin: w.origin || '',
    availability: w.availability || '',
    seen: w.seen,
  }))
  const rows = sample.map((work, i) => ({
    row_index: i,
    work,
  }))
  return { rows }
}

/**
 * Mock: simulate import by creating works for selected rows (or all preview rows).
 * @param {File} [file] - ignored in mock
 * @param {number[]|null} [rowIndices]
 * @returns {Promise<{ created: number, errors?: { row_index: number, message: string }[] }>}
 */
export async function importExecute(file, rowIndices = null) {
  void file
  await Promise.resolve()
  const sample = store.slice(0, 3).map((w) => ({
    type: w.type,
    title: `[Import] ${w.title}`,
    authors: w.authors || [],
    origin: w.origin || '',
    availability: w.availability || '',
    seen: w.seen,
  }))
  const indices = Array.isArray(rowIndices) && rowIndices.length > 0
    ? rowIndices
    : sample.map((_, i) => i)
  let created = 0
  for (const idx of indices) {
    if (idx >= 0 && idx < sample.length) {
      await createWork(sample[idx])
      created++
    }
  }
  return { created, errors: [] }
}
