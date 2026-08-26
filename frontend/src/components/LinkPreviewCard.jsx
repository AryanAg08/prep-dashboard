import { useEffect, useState } from 'react'
import { api } from '../api'

// Renders a rich preview card for a DSA question URL by asking the backend
// to scrape its OpenGraph / meta tags. Falls back to a plain link on error.
export default function LinkPreviewCard({ url, onRemove }) {
  const [data, setData] = useState(null)
  const [status, setStatus] = useState('loading') // loading | ok | error

  useEffect(() => {
    let alive = true
    setStatus('loading')
    api
      .linkPreview(url)
      .then((d) => {
        if (!alive) return
        setData(d)
        setStatus('ok')
      })
      .catch(() => alive && setStatus('error'))
    return () => {
      alive = false
    }
  }, [url])

  if (status === 'loading') {
    return (
      <div className="preview-card skeleton">
        <div className="preview-thumb sk" />
        <div className="preview-body">
          <div className="sk-line sk" />
          <div className="sk-line short sk" />
        </div>
      </div>
    )
  }

  if (status === 'error' || !data) {
    return (
      <div className="preview-card">
        <div className="preview-body">
          <a href={url} target="_blank" rel="noreferrer" className="preview-title">
            {url}
          </a>
          <div className="preview-desc">Could not load preview.</div>
        </div>
        {onRemove && (
          <button className="icon-btn" onClick={onRemove} title="Remove">
            ✕
          </button>
        )}
      </div>
    )
  }

  return (
    <a href={data.url} target="_blank" rel="noreferrer" className="preview-card link">
      {data.image ? (
        <img className="preview-thumb" src={data.image} alt="" loading="lazy" />
      ) : (
        <div className="preview-thumb placeholder">🔗</div>
      )}
      <div className="preview-body">
        <div className="preview-site">
          {data.favicon && <img className="favicon" src={data.favicon} alt="" />}
          <span>{data.siteName}</span>
        </div>
        <div className="preview-title">{data.title}</div>
        {data.description && <div className="preview-desc">{data.description}</div>}
      </div>
      {onRemove && (
        <button
          className="icon-btn"
          onClick={(e) => {
            e.preventDefault()
            onRemove()
          }}
          title="Remove"
        >
          ✕
        </button>
      )}
    </a>
  )
}
