import { createContext, useContext, useState, useEffect } from 'react';
import {
  login as apiLogin,
  signup as apiSignup,
  logout as apiLogout,
  getOrders,
  setToken,
  clearToken,
} from './api';

const AuthContext = createContext(null);

export function AuthProvider({ children }) {
  const [user, setUser] = useState(null);

  useEffect(() => {
    const storedToken = localStorage.getItem('token');
    const name = localStorage.getItem('username');
    if (storedToken && name) {
      setToken(storedToken);
      getOrders()
        .then(() => setUser({ name }))
        .catch(() => {
          localStorage.removeItem('token');
          localStorage.removeItem('username');
          clearToken();
          setUser(null);
        });
    }
  }, []);

  const login = async (username, password, isSignup = false, extra = {}) => {
    const fn = isSignup ? apiSignup : apiLogin;
    const data = isSignup
      ? await fn({ username, password, ...extra })
      : await fn(username, password);
    localStorage.setItem('token', data.token);
    localStorage.setItem('username', username);
    setToken(data.token);
    setUser({ name: username });
  };

  const logout = () => {
    apiLogout();
    localStorage.removeItem('token');
    localStorage.removeItem('username');
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
