import { createContext, useContext, useState } from 'react';
import { login as apiLogin, signup as apiSignup, logout as apiLogout } from './api';

const AuthContext = createContext(null);

export function AuthProvider({ children }) {
  const [user, setUser] = useState(null);

  const login = async (username, password, isSignup = false) => {
    const fn = isSignup ? apiSignup : apiLogin;
    await fn(username, password);
    setUser({ name: username });
  };
  const logout = async () => {
    await apiLogout();
    setUser(null);
  };

  return (
    <AuthContext.Provider value={{ user, login, logout }}>
      {children}
    </AuthContext.Provider>
  );
}

export const useAuth = () => useContext(AuthContext);
