import axios from "axios";

const API_BASE = "http://localhost:8080/api";

const api = axios.create({
  baseURL: API_BASE, // change if needed
  withCredentials: true, // send cookies (refresh token)
});

// Attach authorization header from localStorage
api.interceptors.request.use((config) => {
  const token = localStorage.getItem("token");
  if (token && config && config.headers) {
    config.headers.Authorization = `Bearer ${token}`;
  }
  return config;
});

let isRefreshing = false;
let failedQueue: Array<{
  resolve: (value?: unknown) => void;
  reject: (err: unknown) => void;
  config: any;
}> = [];

function processQueue(error: unknown, token: string | null = null) {
  failedQueue.forEach(({ resolve, reject, config }) => {
    if (error) {
      reject(error);
    } else {
      if (token && config && config.headers) {
        config.headers.Authorization = `Bearer ${token}`;
      }
      resolve(api(config));
    }
  });
  failedQueue = [];
}

api.interceptors.response.use(
  (res) => res,
  async (err) => {
    const originalRequest = err.config;
    if (!originalRequest) return Promise.reject(err);

    if (err.response && err.response.status === 401 && !originalRequest._retry) {
      originalRequest._retry = true;

      if (isRefreshing) {
        return new Promise((resolve, reject) => {
          failedQueue.push({ resolve, reject, config: originalRequest });
        });
      }

      isRefreshing = true;

      try {
        // call refresh endpoint using plain axios to avoid interceptor loop
        const r = await axios.post(`${API_BASE}/auth/refresh`, null, { withCredentials: true });
        const newToken = r.data?.token;
        if (newToken) {
          localStorage.setItem("token", newToken);
          // notify app about token change
          try {
            window.dispatchEvent(new CustomEvent("auth:token", { detail: newToken }));
          } catch (e) {
            // ignore in environments without window
          }
          processQueue(null, newToken);
          if (originalRequest.headers) {
            originalRequest.headers.Authorization = `Bearer ${newToken}`;
          }
          return api(originalRequest);
        }
        processQueue(new Error("no token from refresh"), null);
        return Promise.reject(err);
      } catch (refreshErr) {
        processQueue(refreshErr, null);
        // dispatch logout event
        try {
          window.dispatchEvent(new CustomEvent("auth:token", { detail: null }));
        } catch (e) {}
        return Promise.reject(refreshErr);
      } finally {
        isRefreshing = false;
      }
    }

    return Promise.reject(err);
  }
);

export default api;