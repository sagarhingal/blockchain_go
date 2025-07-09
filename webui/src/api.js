const API_URL = import.meta.env.VITE_API_URL || 'http://localhost:8080';
let chainCache = null;

async function request(path, options = {}) {
  const res = await fetch(`${API_URL}${path}`, {
    credentials: 'include',
    ...options,
  });
  if (!res.ok) throw new Error(await res.text());
  return res.json();
}

export async function getChain(force = false) {
  if (chainCache && !force) return chainCache;
  const data = await request('/chain');
  chainCache = data;
  return data;
}

export async function addTransaction(tx) {
  const data = await request('/transaction', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(tx),
  });
  chainCache = null;
  return data;
}

export async function validateChain() {
  return request('/validate');
}

export async function login(username, password) {
  return request('/login', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ username, password }),
  });
}

export async function signup(username, password) {
  return request('/signup', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ username, password }),
  });
}

export async function resetPassword(password) {
  return request('/reset', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ password }),
  });
}

export async function logout() {
  return request('/logout');
}
