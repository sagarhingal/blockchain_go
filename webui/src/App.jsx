import { BrowserRouter as Router, Routes, Route } from "react-router-dom";
import { AuthProvider, useAuth } from "./AuthContext";
import NavBar from "./components/NavBar";
import Login from "./pages/Login";
import Dashboard from "./pages/Dashboard";
import AddTransaction from "./pages/AddTransaction";
import Chain from "./pages/Chain";
import Settings from "./pages/Settings";
import Validate from "./pages/Validate";

function AppRoutes() {
  const { user } = useAuth();
  return (
    <>
      <NavBar />
      <Routes>
        {user ? (
          <>
            <Route path="/" element={<Dashboard />} />
            <Route path="/transaction" element={<AddTransaction />} />
            <Route path="/chain" element={<Chain />} />
            <Route path="/validate" element={<Validate />} />
            <Route path="/settings" element={<Settings />} />
          </>
        ) : (
          <Route path="/*" element={<Login />} />
        )}
      </Routes>
    </>
  );
}

export default function App() {
  return (
    <Router>
      <AuthProvider>
        <AppRoutes />
      </AuthProvider>
    </Router>
  );
}
