import { useEffect, useState } from 'react'
import { api, TOPICS } from '../api'

const today = () => new Date().toISOString().slice(0, 10)

export default function ActivityLog({ onChange }) {
  const [activities, setActivities] = useState([])
  const [filter, setFilter] = useState('')
  const [form, setForm] = useState({
    date: today(),
    topic: 'dsa',
    title: '',
    description: '',
    hours: 1,
  })
  const [loading, setLoading] = useState(true)

  const load = () => {
    setLoading(true)
    api
      .listActivities(filter ? { topic: filter } : {})
      .then(setActivities)
      .finally(() => setLoading(false))
  }

  useEffect(load, [filter])

  const submit = async (e) => {
    e.preventDefault()
    if (!form.title.trim()) return
    await api.createActivity({ ...form, hours: Number(form.hours) })
    setForm({ ...form, title: '', description: '' })
    load()
    onChange?.()
  }

  const toggleDone = async (a) => {
    await api.updateActivity(a.id, { done: !a.done })
    load()
    onChange?.()
  }

  const remove = async (id) => {
    await api.deleteActivity(id)
    load()
    onChange?.()
  }

  return (
    <div>
      <div className="section-head">
        <h2>Daily Activity Log</h2>
        <select value={filter} onChange={(e) => setFilter(e.target.value)}>
          <option value="">All topics</option>
          {TOPICS.map((t) => (
            <option key={t.key} value={t.key}>
              {t.icon} {t.name}
            </option>
          ))}
        </select>
      </div>

      <form className="card form-grid" onSubmit={submit}>
        <input
          type="date"
          value={form.date}
          onChange={(e) => setForm({ ...form, date: e.target.value })}
        />
        <select value={form.topic} onChange={(e) => setForm({ ...form, topic: e.target.value })}>
          {TOPICS.map((t) => (
            <option key={t.key} value={t.key}>
              {t.icon} {t.name}
            </option>
          ))}
        </select>
        <input
          className="grow"
          placeholder="What did you study / do?"
          value={form.title}
          onChange={(e) => setForm({ ...form, title: e.target.value })}
        />
        <input
          type="number"
          min="0"
          step="0.5"
          className="hours"
          value={form.hours}
          onChange={(e) => setForm({ ...form, hours: e.target.value })}
          title="Hours"
        />
        <input
          className="grow full"
          placeholder="Notes / details (optional)"
          value={form.description}
          onChange={(e) => setForm({ ...form, description: e.target.value })}
        />
        <button type="submit" className="btn primary">
          + Add
        </button>
      </form>

      {loading ? (
        <p className="muted">Loading…</p>
      ) : activities.length === 0 ? (
        <p className="muted">No activities yet. Log your first one above!</p>
      ) : (
        <div className="list">
          {activities.map((a) => {
            const topic = TOPICS.find((t) => t.key === a.topic)
            return (
              <div key={a.id} className={`row ${a.done ? 'done' : ''}`}>
                <input type="checkbox" checked={a.done} onChange={() => toggleDone(a)} />
                <span className="pill">{topic ? `${topic.icon} ${topic.name}` : a.topic}</span>
                <div className="row-main">
                  <div className="row-title">{a.title}</div>
                  {a.description && <div className="row-desc">{a.description}</div>}
                </div>
                <span className="row-meta">{a.date}</span>
                <span className="row-meta">{a.hours}h</span>
                <button className="icon-btn" onClick={() => remove(a.id)} title="Delete">
                  🗑
                </button>
              </div>
            )
          })}
        </div>
      )}
    </div>
  )
}
