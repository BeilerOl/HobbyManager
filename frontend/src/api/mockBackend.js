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
