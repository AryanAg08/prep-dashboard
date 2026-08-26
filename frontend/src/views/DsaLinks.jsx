import { useEffect, useState } from 'react'
import LinkPreviewCard from '../components/LinkPreviewCard'

const STORE_KEY = 'prep.dsaLinks'

// Saved DSA question links live in localStorage — the rich preview for each is
// fetched fresh from the backend scraper on render.
function loadLinks() {
  try {
    return JSON.parse(localStorage.getItem(STORE_KEY) || '[]')
  } catch {
    return []
  }
}

export default function DsaLinks() {
  const [links, setLinks] = useState(loadLinks)
  const [url, setUrl] = useState('')

  useEffect(() => {
    localStorage.setItem(STORE_KEY, JSON.stringify(links))
  }, [links])

  const add = (e) => {
    e.preventDefault()
    const u = url.trim()
    if (!u) return
    if (links.includes(u)) {
      setUrl('')
      return
    }
    setLinks([u, ...links])
    setUrl('')
  }

  const remove = (u) => setLinks(links.filter((x) => x !== u))

  return (
    <div>
      <div className="section-head">
        <h2>DSA Question Links</h2>
        <span className="muted">{links.length} saved</span>
      </div>

      <form className="card link-form" onSubmit={add}>
        <input
          className="grow"
          placeholder="Paste a LeetCode / GFG / Codeforces question URL…"
          value={url}
          onChange={(e) => setUrl(e.target.value)}
        />
        <button type="submit" className="btn primary">
          + Save & Preview
        </button>
      </form>

      {links.length === 0 ? (
        <p className="muted">
          No links saved yet. Paste a DSA question URL above to get a rich preview card.
        </p>
      ) : (
        <div className="preview-list">
          {links.map((u) => (
            <LinkPreviewCard key={u} url={u} onRemove={() => remove(u)} />
          ))}
        </div>
      )}
    </div>
  )
}
