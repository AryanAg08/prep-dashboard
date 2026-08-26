import { useEffect, useState } from 'react'
import { api, TOPICS } from '../api'

export default function Notes() {
  const [notes, setNotes] = useState([])
  const [filter, setFilter] = useState('')
  const [editing, setEditing] = useState(null) // note being edited, or null
  const [form, setForm] = useState({ topic: 'dsa', title: '', content: '', tags: '' })
  const [loading, setLoading] = useState(true)

  const load = () => {
    setLoading(true)
    api
      .listNotes(filter ? { topic: filter } : {})
      .then(setNotes)
      .finally(() => setLoading(false))
  }

  useEffect(load, [filter])

  const resetForm = () => {
    setForm({ topic: 'dsa', title: '', content: '', tags: '' })
    setEditing(null)
  }

  const submit = async (e) => {
    e.preventDefault()
    if (!form.title.trim()) return
    const payload = {
      topic: form.topic,
      title: form.title,
      content: form.content,
      tags: form.tags
        .split(',')
        .map((t) => t.trim())
        .filter(Boolean),
    }
    if (editing) await api.updateNote(editing, payload)
    else await api.createNote(payload)
    resetForm()
    load()
  }

  const startEdit = (n) => {
    setEditing(n.id)
    setForm({ topic: n.topic, title: n.title, content: n.content, tags: (n.tags || []).join(', ') })
    window.scrollTo({ top: 0, behavior: 'smooth' })
  }

  const remove = async (id) => {
    await api.deleteNote(id)
    if (editing === id) resetForm()
    load()
  }

  return (
    <div>
      <div className="section-head">
        <h2>Notes</h2>
        <select value={filter} onChange={(e) => setFilter(e.target.value)}>
          <option value="">All topics</option>
          {TOPICS.map((t) => (
            <option key={t.key} value={t.key}>
              {t.icon} {t.name}
            </option>
          ))}
        </select>
      </div>

      <form className="card note-form" onSubmit={submit}>
        <div className="note-form-top">
          <select value={form.topic} onChange={(e) => setForm({ ...form, topic: e.target.value })}>
            {TOPICS.map((t) => (
              <option key={t.key} value={t.key}>
                {t.icon} {t.name}
              </option>
            ))}
          </select>
          <input
            className="grow"
            placeholder="Note title"
            value={form.title}
            onChange={(e) => setForm({ ...form, title: e.target.value })}
          />
        </div>
        <textarea
          rows="4"
          placeholder="Write your note… (concepts, gotchas, formulas)"
          value={form.content}
          onChange={(e) => setForm({ ...form, content: e.target.value })}
        />
        <div className="note-form-bottom">
          <input
            className="grow"
            placeholder="tags, comma separated"
            value={form.tags}
            onChange={(e) => setForm({ ...form, tags: e.target.value })}
          />
          {editing && (
            <button type="button" className="btn" onClick={resetForm}>
              Cancel
            </button>
          )}
          <button type="submit" className="btn primary">
            {editing ? 'Save' : '+ Add Note'}
          </button>
        </div>
      </form>

      {loading ? (
        <p className="muted">Loading…</p>
      ) : notes.length === 0 ? (
        <p className="muted">No notes yet.</p>
      ) : (
        <div className="notes-grid">
          {notes.map((n) => {
            const topic = TOPICS.find((t) => t.key === n.topic)
            return (
              <div key={n.id} className="note-card">
                <div className="note-card-head">
                  <span className="pill">{topic ? `${topic.icon} ${topic.name}` : n.topic}</span>
                  <div className="note-actions">
                    <button className="icon-btn" onClick={() => startEdit(n)} title="Edit">
                      ✎
                    </button>
                    <button className="icon-btn" onClick={() => remove(n.id)} title="Delete">
                      🗑
                    </button>
                  </div>
                </div>
                <div className="note-title">{n.title}</div>
                {n.content && <div className="note-content">{n.content}</div>}
                {n.tags?.length > 0 && (
                  <div className="tags">
                    {n.tags.map((t) => (
                      <span key={t} className="tag">
                        #{t}
                      </span>
                    ))}
                  </div>
                )}
              </div>
            )
          })}
        </div>
      )}
    </div>
  )
}
