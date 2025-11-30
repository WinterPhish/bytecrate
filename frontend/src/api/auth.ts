export async function login(email: string, password: string) {
  const res = await fetch("http://localhost:8080/api/auth/login", {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ email, password }),
  });

  if (!res.ok) throw new Error("Login failed");

  return res.json();
}

export async function register(email: string, password: string) {
  const res = await fetch("http://localhost:8080/api/auth/register", {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ email, password }),
  });

  if (!res.ok) throw new Error("Register failed");

  return res.json();
}

export async function authenticatedFetch(url: string, options: RequestInit = {}) {
  const token = localStorage.getItem("token");
  return fetch(url, {
    ...options,
    headers: {
      ...options.headers,
      ...(token && { Authorization: `Bearer ${token}` }),
    },
  });
}
