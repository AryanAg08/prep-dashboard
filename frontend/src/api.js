// Thin API client. In dev, Vite proxies /api to the Go backend (see vite.config.js).
const BASE = '/api'
const KEY_STORAGE = 'prep-api-key'

// The password is typed in by the user and kept in localStorage rather than
// baked into the bundle -- this is a static SPA, so anything built in ships to
// every visitor in plain sight. It doubles as the API key on later requests.
export function getApiKey() {
  return localStorage.getItem(KEY_STORAGE) || ''
}

export function setApiKey(key) {
  localStorage.setItem(KEY_STORAGE, key.trim())
}

export function clearApiKey() {
  localStorage.removeItem(KEY_STORAGE)
}

// Thrown on 401 so the UI can send the user back to the key prompt.
export class UnauthorizedError extends Error {
  constructor() {
    super('Invalid or missing API key')
    this.name = 'UnauthorizedError'
  }
}

async function req(path, options = {}) {
  const res = await fetch(BASE + path, {
    ...options,
    headers: {
      'Content-Type': 'application/json',
      'X-API-Key': getApiKey(),
      ...(options.headers || {}),
    },
  })
  if (res.status === 401) {
    throw new UnauthorizedError()
  }
  if (!res.ok) {
    const text = await res.text().catch(() => '')
    throw new Error(`${res.status} ${text}`)
  }
  if (res.status === 204) return null
  return res.json()
}

// Verifies a password against the backend and stores it on success. Throws
// UnauthorizedError on a bad password, or a plain Error if the API is down or
// the attempt was rate limited.
export async function login(password) {
  const res = await fetch(BASE + '/login', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ password }),
  })
  if (res.status === 401) throw new UnauthorizedError()
  if (res.status === 429) throw new Error('Too many attempts. Wait a minute and try again.')
  if (!res.ok) throw new Error(`${res.status} ${await res.text().catch(() => '')}`)

  setApiKey(password)
}

// Confirms an already-stored password still works, without a login round trip.
export async function verifyStoredKey() {
  await req('/guides')
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
