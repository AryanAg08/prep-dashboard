import { useEffect, useState } from 'react'
import { api, TOPICS } from '../api'

export default function Overview({ refreshKey }) {
  const [stats, setStats] = useState(null)
  const [recent, setRecent] = useState([])

  useEffect(() => {
    api.getStats().then(setStats)
    api.listActivities().then((a) => setRecent(a.slice(0, 5)))
  }, [refreshKey])

  const topicHours = {}
  stats?.perTopic?.forEach((row) => {
    topicHours[row._id] = { count: row.count, hours: row.hours }
  })
  const maxHours = Math.max(1, ...Object.values(topicHours).map((t) => t.hours || 0))

  return (
    <div>
      <div className="section-head">
        <h2>Overview</h2>
      </div>

      <div className="stat-row">
        <div className="stat">
          <div className="stat-num">{stats?.total ?? '—'}</div>
          <div className="stat-label">Total Activities</div>
        </div>
        <div className="stat">
          <div className="stat-num">{stats?.done ?? '—'}</div>
          <div className="stat-label">Completed</div>
        </div>
        <div className="stat">
          <div className="stat-num">{stats?.today ?? '—'}</div>
          <div className="stat-label">Logged Today</div>
        </div>
        <div className="stat">
          <div className="stat-num">{stats ? Math.round(stats.totalHours) : '—'}</div>
          <div className="stat-label">Total Hours</div>
        </div>
      </div>

      <div className="guide-cols">
        <div className="card">
          <h3>⏱️ Hours by Topic</h3>
          {Object.keys(topicHours).length === 0 ? (
            <p className="muted">Log activities to see your breakdown.</p>
          ) : (
            <div className="bars">
              {TOPICS.filter((t) => topicHours[t.key]).map((t) => {
                const h = topicHours[t.key].hours || 0
                return (
                  <div key={t.key} className="bar-row">
                    <span className="bar-label">
                      {t.icon} {t.name}
                    </span>
                    <div className="bar-track">
                      <div className="bar-fill" style={{ width: `${(h / maxHours) * 100}%` }} />
                    </div>
                    <span className="bar-val">{h}h</span>
                  </div>
                )
              })}
            </div>
          )}
        </div>

        <div className="card">
          <h3>🕑 Recent Activity</h3>
          {recent.length === 0 ? (
            <p className="muted">Nothing yet.</p>
          ) : (
            <div className="list tight">
              {recent.map((a) => {
                const topic = TOPICS.find((t) => t.key === a.topic)
                return (
                  <div key={a.id} className="row tight">
                    <span className="pill sm">{topic?.icon}</span>
                    <div className="row-main">
                      <div className="row-title">{a.title}</div>
                    </div>
                    <span className="row-meta">{a.date}</span>
                  </div>
                )
              })}
            </div>
          )}
        </div>
      </div>
    </div>
  )
}
