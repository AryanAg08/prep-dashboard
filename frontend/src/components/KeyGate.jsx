import { useEffect, useState } from 'react'
import { getApiKey, clearApiKey, login, verifyStoredKey } from '../api'

// Wraps the app and withholds it until a correct password has been entered.
// The password lives in localStorage, so this is a one-time prompt per browser.
export default function KeyGate({ children }) {
  const [status, setStatus] = useState('checking') // checking | locked | ready
  const [input, setInput] = useState('')
  const [error, setError] = useState('')
  const [busy, setBusy] = useState(false)

  // Validate whatever is already stored before showing the dashboard.
  useEffect(() => {
    if (!getApiKey()) {
      setStatus('locked')
      return
    }
    verifyStoredKey()
      .then(() => setStatus('ready'))
      .catch((err) => {
        if (err.name === 'UnauthorizedError') {
          clearApiKey()
          setError('Your saved password was rejected. Enter it again.')
        } else {
          setError(`Could not reach the API: ${err.message}`)
        }
        setStatus('locked')
      })
  }, [])

  async function submit(e) {
    e.preventDefault()
    if (!input.trim() || busy) return

    setBusy(true)
    setError('')
    try {
      await login(input)
      setInput('')
      setStatus('ready')
    } catch (err) {
      clearApiKey()
      setError(err.name === 'UnauthorizedError' ? 'Incorrect password.' : err.message)
    } finally {
      setBusy(false)
    }
  }

  if (status === 'checking') {
    return <div className="gate"><p className="muted">Unlocking…</p></div>
  }

  if (status === 'ready') return children

  return (
    <div className="gate">
      <form className="gate-card" onSubmit={submit}>
        <div className="gate-icon">🔐</div>
        <h1 className="gate-title">Prep Dashboard</h1>
        <p className="muted">Enter your password to continue.</p>
        <input
          className="gate-input"
          type="password"
          value={input}
          onChange={(e) => setInput(e.target.value)}
          placeholder="Password"
          autoFocus
        />
        {error && <p className="gate-error">{error}</p>}
        <button className="btn btn-primary gate-btn" type="submit" disabled={busy}>
          {busy ? 'Checking…' : 'Sign in'}
        </button>
      </form>
    </div>
  )
}
