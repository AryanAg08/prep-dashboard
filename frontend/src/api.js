// Thin API client. In dev, Vite proxies /api to the Go backend (see vite.config.js).
const BASE = '/api'

async function req(path, options = {}) {
  const res = await fetch(BASE + path, {
    headers: { 'Content-Type': 'application/json' },
    ...options,
  })
  if (!res.ok) {
    const text = await res.text().catch(() => '')
    throw new Error(`${res.status} ${text}`)
  }
  if (res.status === 204) return null
  return res.json()
}

export const api = {
  // Study guides / materials
  getGuides: () => req('/guides'),

  // Activities
  listActivities: (params = {}) => {
    const q = new URLSearchParams(params).toString()
    return req('/activities' + (q ? `?${q}` : ''))
  },
  createActivity: (data) => req('/activities', { method: 'POST', body: JSON.stringify(data) }),
  updateActivity: (id, data) => req(`/activities/${id}`, { method: 'PATCH', body: JSON.stringify(data) }),
  deleteActivity: (id) => req(`/activities/${id}`, { method: 'DELETE' }),
  getStats: () => req('/stats'),

  // Notes
  listNotes: (params = {}) => {
    const q = new URLSearchParams(params).toString()
    return req('/notes' + (q ? `?${q}` : ''))
  },
  createNote: (data) => req('/notes', { method: 'POST', body: JSON.stringify(data) }),
  updateNote: (id, data) => req(`/notes/${id}`, { method: 'PATCH', body: JSON.stringify(data) }),
  deleteNote: (id) => req(`/notes/${id}`, { method: 'DELETE' }),

  // Link preview (DSA questions)
  linkPreview: (url) => req('/link-preview', { method: 'POST', body: JSON.stringify({ url }) }),
}

export const TOPICS = [
  { key: 'dsa', name: 'DSA', icon: '🧮' },
  { key: 'lld', name: 'LLD', icon: '🧱' },
  { key: 'projects', name: 'Projects', icon: '🚀' },
  { key: 'os', name: 'OS Fundamentals', icon: '💻' },
  { key: 'dbms', name: 'DBMS + SQL', icon: '🗄️' },
  { key: 'hr', name: 'HR Questions', icon: '🤝' },
]
