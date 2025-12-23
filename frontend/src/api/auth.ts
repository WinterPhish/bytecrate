import axios from "axios";
import api from "./axios";

export async function login(email: string, password: string) {
  try {
    const res = await api.post("/auth/login", { email, password });
    return res.data; // parsed JSON
  } catch (err: unknown) {
    if (axios.isAxiosError(err)) {
      throw new Error(err.response?.data || err.message);
    }
    if (err instanceof Error) {
      throw new Error(err.message);
    }
    throw new Error("Login failed");
  }
}

export async function register(email: string, password: string) {
  try {
    const res = await api.post("/auth/register", { email, password });
    return res.data;
  } catch (err: unknown) {
    if (axios.isAxiosError(err)) {
      throw new Error(err.response?.data || err.message);
    }
    if (err instanceof Error) {
      throw new Error(err.message);
    }
    throw new Error("Register failed");
  }
}

// call refresh endpoint to rotate refresh token and obtain a new access token
export async function refreshToken() {
  try {
    const res = await axios.post("http://localhost:8080/api/auth/refresh", null, { withCredentials: true });
    return res.data;
  } catch (err: unknown) {
    if (axios.isAxiosError(err)) {
      throw new Error(err.response?.data || err.message);
    }
    if (err instanceof Error) {
      throw new Error(err.message);
    }
    throw new Error("Refresh failed");
  }
}