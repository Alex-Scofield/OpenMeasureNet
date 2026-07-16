const API_URL = '/api'

async function request(path, options = {}) {
  const token = localStorage.getItem('token')
  const headers = {
    'Content-Type': 'application/json',
    ...options.headers,
  }
  if (token) {
    headers['Authorization'] = `Bearer ${token}`
  }

  const res = await fetch(`${API_URL}${path}`, { ...options, headers })

  if (res.status === 401) {
    localStorage.removeItem('token')
    window.location.href = '/login'
    return
  }

  return res
}

export async function signup(email, password, code) {
  const res = await request('/auth/signup', {
    method: 'POST',
    body: JSON.stringify({ email, password, code }),
  })
  return res.json()
}

export async function login(email, password) {
  const res = await request('/auth/login', {
    method: 'POST',
    body: JSON.stringify({ email, password }),
  })
  return res.json()
}

export async function logout() {
  await request('/auth/logout', { method: 'POST' })
  localStorage.removeItem('token')
}

export async function getNodes() {
  const res = await request('/nodes')
  return res.json()
}

export async function getNode(id) {
  const res = await request(`/nodes/${id}`)
  return res.json()
}

export async function getMap() {
  const res = await request('/map')
  return res.json()
}
