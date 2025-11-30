import { useState } from "react";
import { useNavigate } from "react-router-dom";
import { register } from "../api/auth";
import { useAuth } from "../hooks/AuthContext";

export default function Register() {
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [error, setError] = useState("");
  const { login: setAuth } = useAuth();
  const navigate = useNavigate();

  async function handleRegister(e: React.FormEvent) {
    e.preventDefault();
    setError("");

    try {
      const data = await register(
        email,
        password,
    );

      setAuth(data.token);
      navigate("/dashboard");
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