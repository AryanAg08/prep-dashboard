import { useState } from 'react'
import KeyGate from './components/KeyGate'
import Overview from './views/Overview'
import ActivityLog from './views/ActivityLog'
import Notes from './views/Notes'
import StudyMaterials from './views/StudyMaterials'
import DsaLinks from './views/DsaLinks'

const NAV = [
  { key: 'overview', label: 'Overview', icon: '📊' },
  { key: 'activity', label: 'Activity Log', icon: '📝' },
  { key: 'notes', label: 'Notes', icon: '🗒️' },
  { key: 'materials', label: 'Study Materials', icon: '📚' },
  { key: 'dsa', label: 'DSA Links', icon: '🔗' },
]

export default function App() {
  const [view, setView] = useState('overview')
  // Bumped whenever activities change so Overview re-fetches its stats.
  const [refreshKey, setRefreshKey] = useState(0)
  const bump = () => setRefreshKey((k) => k + 1)

  return (
    <KeyGate>
      <div className="app">
        <aside className="sidebar">
          <div className="brand">
            <span className="brand-icon">🎯</span>
            <div>
              <div className="brand-title">Prep Dashboard</div>
              <div className="brand-sub">Placement Tracker</div>
            </div>
          </div>
          <nav>
            {NAV.map((n) => (
              <button
                key={n.key}
                className={`nav-item ${view === n.key ? 'active' : ''}`}
                onClick={() => setView(n.key)}
              >
                <span className="nav-icon">{n.icon}</span>
                {n.label}
              </button>
            ))}
          </nav>
          <div className="sidebar-foot">Stay consistent. Ship daily. 🚀</div>
        </aside>

        <main className="content">
          {view === 'overview' && <Overview refreshKey={refreshKey} />}
          {view === 'activity' && <ActivityLog onChange={bump} />}
          {view === 'notes' && <Notes />}
          {view === 'materials' && <StudyMaterials />}
          {view === 'dsa' && <DsaLinks />}
        </main>
      </div>
    </KeyGate>
  )
}
