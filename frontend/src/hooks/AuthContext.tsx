import React, { createContext, useContext, useEffect, useState } from "react";

interface AuthContextType {
  token: string | null;
  login: (jwt: string) => void;
  logout: () => void;
}

const AuthContext = createContext<AuthContextType | undefined>(undefined);

export const AuthProvider: React.FC<{ children: React.ReactNode }> = ({ children }) => {
  const [token, setToken] = useState<string | null>(
    localStorage.getItem("token")
  );

  const login = (jwt: string) => {
    localStorage.setItem("token", jwt);
    setToken(jwt);
  };

  const logout = () => {
    localStorage.removeItem("token");
    setToken(null);
  };

  useEffect(() => {
    // listen for token updates from axios interceptor (refresh)
    const handler = (e: Event) => {
      // event may be CustomEvent with detail = token
      const ce = e as CustomEvent;
      const t = ce?.detail ?? null;
      if (t) {
        setToken(t);
        localStorage.setItem("token", t);
      } else {
        setToken(null);
        localStorage.removeItem("token");
      }
    };
    window.addEventListener("auth:token", handler as EventListener);
    return () => window.removeEventListener("auth:token", handler as EventListener);
  }, []);

  return (
    <AuthContext.Provider value={{ token, login, logout }}>
      {children}
    </AuthContext.Provider>
  );
};

// possibly fix
// eslint-disable-next-line react-refresh/only-export-components
export const useAuth = (): AuthContextType => {
  const ctx = useContext(AuthContext);
  if (!ctx) {
    throw new Error("useAuth must be used inside <AuthProvider>");
  }
  return ctx;
};
