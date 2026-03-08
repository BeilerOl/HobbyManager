/**
 * API client: real backend (fetch) or mock (in-memory) when VITE_USE_MOCK is set (e.g. GitHub Pages build).
 */
import * as realApi from './realApi.js'
import * as mockApi from './mockBackend.js'

const api = import.meta.env.VITE_USE_MOCK === 'true' ? mockApi : realApi

export const listWorks = api.listWorks
export const getWork = api.getWork
export const createWork = api.createWork
export const updateWork = api.updateWork
export const deleteWork = api.deleteWork
export const WORK_TYPES = api.WORK_TYPES
