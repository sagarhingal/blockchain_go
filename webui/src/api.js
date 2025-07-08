const API_URL = import.meta.env.VITE_API_URL || 'http://localhost:8080';
let chainCache = null;

export async function getChain(force = false) {
  if (chainCache && !force) return chainCache;
  const res = await fetch(`${API_URL}/chain`);
  if (!res.ok) throw new Error('Failed to fetch chain');
  const data = await res.json();
  chainCache = data;
  return data;
}

export async function addTransaction(tx) {
  const res = await fetch(`${API_URL}/transaction`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(tx),
  });
  if (!res.ok) throw new Error('Failed to add transaction');
  chainCache = null;
  return res.json();
}

export async function validateChain() {
  const res = await fetch(`${API_URL}/validate`);
  if (!res.ok) throw new Error('Failed to validate chain');
  return res.json();
}
