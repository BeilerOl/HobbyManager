/**
 * Import API: preview and execute Excel import.
 * Uses real backend or mock when VITE_USE_MOCK is set.
 */
import * as realApi from './realApi.js'
import * as mockApi from './mockBackend.js'

const api = import.meta.env.VITE_USE_MOCK === 'true' ? mockApi : realApi

/**
 * @typedef {Object} ImportPreviewRow
 * @property {number} row_index
 * @property {string} [error]
 * @property {import('./realApi.js').WorkCreate} [work]
 */

/**
 * @typedef {Object} ImportPreviewResponse
 * @property {ImportPreviewRow[]} rows
 */

/**
 * @param {File} file - .xls file
 * @returns {Promise<ImportPreviewResponse>}
 */
export async function importPreview(file) {
  return api.importPreview(file)
}

/**
 * @typedef {Object} ImportExecuteResponse
 * @property {number} created
 * @property {{ row_index: number, message: string }[]} [errors]
 */

/**
 * @param {File} file - .xls file
 * @param {number[]} [rowIndices] - optional indices of rows to import (if omitted, all valid rows)
 * @returns {Promise<ImportExecuteResponse>}
 */
export async function importExecute(file, rowIndices = null) {
  return api.importExecute(file, rowIndices)
}
