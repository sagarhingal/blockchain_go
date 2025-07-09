import { createContext, useContext, useState, useEffect } from 'react';
import { login as apiLogin, signup as apiSignup, logout as apiLogout, getChain } from './api';

const AuthContext = createContext(null);

export function AuthProvider({ children }) {
  const [user, setUser] = useState(null);

  useEffect(() => {
    const stored = localStorage.getItem('username');
    if (stored) {
      getChain()
        .then(() => setUser({ name: stored }))
        .catch(() => {
          localStorage.removeItem('username');
          setUser(null);
        });
    }
  }, []);

  const login = async (username, password, isSignup = false) => {
    const fn = isSignup ? apiSignup : apiLogin;
    await fn(username, password);
    localStorage.setItem('username', username);
    setUser({ name: username });
  };
  const logout = async () => {
    await apiLogout();
    localStorage.removeItem('username');
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
