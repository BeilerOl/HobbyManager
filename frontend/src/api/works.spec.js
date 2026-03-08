import { describe, it, expect, vi, beforeEach } from 'vitest'
import { listWorks, getWork, createWork, updateWork, deleteWork, WORK_TYPES } from './works'

describe('works API', () => {
  beforeEach(() => {
    vi.stubGlobal('fetch', vi.fn())
  })

  it('exports WORK_TYPES with expected values', () => {
    expect(WORK_TYPES).toBeDefined()
    expect(WORK_TYPES.length).toBeGreaterThan(0)
    expect(WORK_TYPES.find(t => t.value === 'roman')).toEqual({ value: 'roman', label: 'Roman' })
  })

  it('listWorks calls GET /api/v1/works', async () => {
    fetch.mockResolvedValueOnce({ ok: true, json: async () => [] })
    await listWorks()
    expect(fetch).toHaveBeenCalledWith('/api/v1/works')
  })

  it('listWorks appends query params when provided', async () => {
    fetch.mockResolvedValueOnce({ ok: true, json: async () => [] })
    await listWorks({ type: 'film', seen: true })
    expect(fetch).toHaveBeenCalledWith(expect.stringContaining('type=film'))
    expect(fetch).toHaveBeenCalledWith(expect.stringContaining('seen=true'))
  })

  it('getWork returns null on 404', async () => {
    fetch.mockResolvedValueOnce({ status: 404 })
    const result = await getWork(1)
    expect(result).toBeNull()
  })

  it('createWork sends POST with JSON body', async () => {
    const data = { type: 'roman', title: 'T', authors: ['A'], origin: '', availability: '', seen: false }
    fetch.mockResolvedValueOnce({ ok: true, json: async () => ({ id: 1, ...data }) })
    await createWork(data)
    expect(fetch).toHaveBeenCalledWith('/api/v1/works', expect.objectContaining({
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(data),
    }))
  })
})
