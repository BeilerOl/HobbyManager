const API_BASE = '/api/v1'

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

/**
 * @param {{ type?: string, seen?: boolean }} [params]
 * @returns {Promise<Work[]>}
 */
export async function listWorks(params = {}) {
  const q = new URLSearchParams()
  if (params.type) q.set('type', params.type)
  if (params.seen !== undefined) q.set('seen', String(params.seen))
  const url = `${API_BASE}/works${q.toString() ? '?' + q : ''}`
  const res = await fetch(url)
  if (!res.ok) throw new Error(await res.text())
  return res.json()
}

/**
 * @param {number} id
 * @returns {Promise<Work>}
 */
export async function getWork(id) {
  const res = await fetch(`${API_BASE}/works/${id}`)
  if (res.status === 404) return null
  if (!res.ok) throw new Error(await res.text())
  return res.json()
}

/**
 * @param {WorkCreate} data
 * @returns {Promise<Work>}
 */
export async function createWork(data) {
  const res = await fetch(`${API_BASE}/works`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(data),
  })
  if (!res.ok) throw new Error(await res.text())
  return res.json()
}

/**
 * @param {number} id
 * @param {WorkCreate} data
 * @returns {Promise<Work>}
 */
export async function updateWork(id, data) {
  const res = await fetch(`${API_BASE}/works/${id}`, {
    method: 'PUT',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(data),
  })
  if (!res.ok) throw new Error(await res.text())
  return res.json()
}

/**
 * @param {number} id
 * @returns {Promise<void>}
 */
export async function deleteWork(id) {
  const res = await fetch(`${API_BASE}/works/${id}`, { method: 'DELETE' })
  if (res.status === 404) return
  if (!res.ok) throw new Error(await res.text())
}

export const WORK_TYPES = [
  { value: 'roman', label: 'Roman' },
  { value: 'livre_culture_generale', label: 'Livre de culture générale' },
  { value: 'film', label: 'Film' },
  { value: 'serie_tv', label: 'Série TV' },
  { value: 'jeu_societe', label: 'Jeu de société' },
  { value: 'jeu_video', label: 'Jeu vidéo' },
]
