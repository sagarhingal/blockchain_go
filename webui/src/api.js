const API_URL = import.meta.env.VITE_API_URL || "http://localhost:8080";
let ordersCache = null;
let token = localStorage.getItem("token") || null;

export function setToken(t) {
  token = t;
}

export function clearToken() {
  token = null;
}

async function request(path, options = {}) {
  const opts = { headers: {}, ...options };
  if (token) {
    opts.headers = { ...opts.headers, Authorization: `Bearer ${token}` };
  }
  const res = await fetch(`${API_URL}${path}`, opts);
  if (!res.ok) throw new Error(await res.text());
  return res.json();
}

export async function getOrders(force = false) {
  if (ordersCache && !force) return ordersCache;
  const data = await request("/chain");
  ordersCache = data;
  return data;
}

export async function createOrder() {
  const data = await request("/order", { method: "POST" });
  ordersCache = null;
  return data;
}

export async function updateOrderStatus(id, status) {
  return request(`/order/${id}/status`, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ status }),
  });
}

export async function getOrderEvents(id) {
  return request(`/order/${id}/events`);
}

export async function addRole(id, actor, role) {
  return request(`/order/${id}/roles`, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ Actor: actor, Role: role }),
  });
}

export async function inviteWatcher(id, actor) {
  return request(`/order/${id}/invite`, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ Actor: actor }),
  });
}

export async function addAddon(id, details) {
  return request(`/order/${id}/addon`, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ Details: details }),
  });
}

export async function listActors() {
  return request("/marketplace");
}

export async function login(email, password) {
  const data = await request("/login", {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ email, password }),
  });
  return data;
}

export async function signup(user) {
  const data = await request("/signup", {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(user),
  });
  return data;
}

export async function resetPassword(info) {
  return request("/reset", {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(info),
  });
}

export function logout() {
  localStorage.removeItem("token");
  ordersCache = null;
  clearToken();
}
