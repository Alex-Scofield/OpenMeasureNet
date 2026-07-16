import { useState } from 'react'
import { useNavigate, Link } from 'react-router-dom'
import { signup } from '../api'

export default function Signup() {
  const [email, setEmail] = useState('')
  const [password, setPassword] = useState('')
  const [code, setCode] = useState('')
  const [error, setError] = useState('')
  const [success, setSuccess] = useState(false)
  const navigate = useNavigate()

  async function handleSubmit(e) {
    e.preventDefault()
    setError('')
    const data = await signup(email, password, code)
    if (data.message) {
      setSuccess(true)
    } else {
      setError(data.error || 'Signup failed')
    }
  }

  if (success) {
    return (
      <div className="auth-page">
        <div className="auth-form">
          <h1>OpenMeasureNet</h1>
          <p className="success">Account created! You can now log in.</p>
          <Link to="/login"><button>Log in</button></Link>
        </div>
      </div>
    )
  }

  return (
    <div className="auth-page">
      <form className="auth-form" onSubmit={handleSubmit}>
        <h1>OpenMeasureNet</h1>
        <input
          type="email"
          placeholder="Email"
          value={email}
          onChange={e => setEmail(e.target.value)}
          required
        />
        <input
          type="password"
          placeholder="Password (min 8 characters)"
          value={password}
          onChange={e => setPassword(e.target.value)}
          minLength={8}
          required
        />
        <input
          type="text"
          placeholder="Invite code"
          value={code}
          onChange={e => setCode(e.target.value)}
          required
        />
        {error && <p className="error">{error}</p>}
        <button type="submit">Sign up</button>
        <p className="alt-link">
          Already have an account? <Link to="/login">Log in</Link>
        </p>
      </form>
    </div>
  )
}
