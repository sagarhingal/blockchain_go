import { createContext, useContext, useState, useEffect } from "react";
import {
  login as apiLogin,
  signup as apiSignup,
  logout as apiLogout,
  getOrders,
  setToken,
  clearToken,
} from "./api";

const AuthContext = createContext(null);

export function AuthProvider({ children }) {
  const [user, setUser] = useState(null);

  useEffect(() => {
    const storedToken = localStorage.getItem("token");
    const name = localStorage.getItem("email");
    if (storedToken && name) {
      setToken(storedToken);
      getOrders()
        .then(() => setUser({ name }))
        .catch(() => {
          localStorage.removeItem("token");
          localStorage.removeItem("email");
          clearToken();
          setUser(null);
        });
    }
  }, []);

  const login = async (email, password, isSignup = false, extra = {}) => {
    const fn = isSignup ? apiSignup : apiLogin;
    const data = isSignup
      ? await fn({ email, password, ...extra })
      : await fn(email, password);
    localStorage.setItem("token", data.token);
    localStorage.setItem("email", email);
    setToken(data.token);
    setUser({ name: email });
  };

  const logout = () => {
    apiLogout();
    localStorage.removeItem("token");
    localStorage.removeItem("email");
    clearToken();
    setUser(null);
  };

  return (
    <AuthContext.Provider value={{ user, login, logout }}>
      {children}
    </AuthContext.Provider>
  );
}

// eslint-disable-next-line react-refresh/only-export-components
export const useAuth = () => useContext(AuthContext);
