import { useState } from "react";
import { login } from "../api/auth";

export default function Login() {
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [error, setError] = useState("");

  async function handleLogin(e: React.FormEvent) {
    e.preventDefault();
    setError("");

    try {
      const data = await login(
        email,
        password,
      )

      localStorage.setItem("token", data.token);
      alert("Logged in!");
    } catch (err: unknown) {
      setError(err instanceof Error ? err.message : "Login failed");
    }
  }

  return (
    <div style={{ padding: "2rem" }}>
      <h1>Login</h1>
      <form onSubmit={handleLogin}>
        <input
          placeholder="Email"
          value={email}
          onChange={e => setEmail(e.target.value)}
        /><br/>

        <input
          placeholder="Password"
          type="password"
          value={password}
          onChange={e => setPassword(e.target.value)}
        /><br/>

        <button>Login</button>

        {error && <p style={{ color: "red" }}>{error}</p>}
      </form>
    </div>
  );
}