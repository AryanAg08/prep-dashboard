import { useEffect, useState } from 'react'
import { api } from '../api'

export default function StudyMaterials() {
  const [guides, setGuides] = useState([])
  const [active, setActive] = useState(null)

  useEffect(() => {
    api.getGuides().then((g) => {
      setGuides(g)
      setActive(g[0]?.key)
    })
  }, [])

  const guide = guides.find((g) => g.key === active)

  return (
    <div>
      <div className="section-head">
        <h2>Study Materials</h2>
        <span className="muted">Curated to crack top companies</span>
      </div>

      <div className="topic-tabs">
        {guides.map((g) => (
          <button
            key={g.key}
            className={`topic-tab ${active === g.key ? 'active' : ''}`}
            onClick={() => setActive(g.key)}
          >
            <span>{g.icon}</span> {g.name}
          </button>
        ))}
      </div>

      {guide && (
        <div className="guide">
          <p className="guide-summary">{guide.summary}</p>

          <div className="guide-cols">
            <div className="card guide-roadmap">
              <h3>🗺️ Roadmap</h3>
              <ol>
                {guide.roadmap.map((step, i) => (
                  <li key={i}>{step}</li>
                ))}
              </ol>
            </div>

            <div className="guide-resources">
              <h3>📚 Resources</h3>
              {guide.resources.map((r, i) => (
                <a key={i} href={r.url} target="_blank" rel="noreferrer" className="resource">
                  <span className={`res-type ${r.type}`}>{r.type}</span>
                  <div className="res-body">
                    <div className="res-title">{r.title}</div>
                    {r.note && <div className="res-note">{r.note}</div>}
                  </div>
                  <span className="res-arrow">↗</span>
                </a>
              ))}
            </div>
          </div>
        </div>
      )}
    </div>
  )
}
