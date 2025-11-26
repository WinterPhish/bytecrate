import { useState } from "react";
import { register } from "../api/auth";

export default function Register() {
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [error, setError] = useState("");

  async function handleRegister(e: React.FormEvent) {
    e.preventDefault();
    setError("");

    try {
      const data = await register(
        email,
        password,
    );

      localStorage.setItem("token", data.token);
      alert("Account created!");
    } catch (err: unknown) {
      setError(err instanceof Error ? err.message : "Login failed");
    }
  }

  return (
    <div style={{ padding: "2rem" }}>
      <h1>Register</h1>
      <form onSubmit={handleRegister}>
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

        <button>Create Account</button>

        {error && <p style={{ color: "red" }}>{error}</p>}
      </form>
    </div>
  );
}