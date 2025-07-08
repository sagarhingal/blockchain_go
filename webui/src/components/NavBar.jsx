import { Link } from "react-router-dom";
import { useAuth } from "../AuthContext";

export default function NavBar() {
  const { user, logout } = useAuth();

  if (!user) return null;

  return (
    <nav style={{ padding: "1rem", borderBottom: "1px solid #ccc" }}>
      <Link to="/" style={{ marginRight: "1rem" }}>
        Dashboard
      </Link>
      <Link to="/chain" style={{ marginRight: "1rem" }}>
        Chain
      </Link>
      <Link to="/transaction" style={{ marginRight: "1rem" }}>
        Add Tx
      </Link>
      <Link to="/validate" style={{ marginRight: "1rem" }}>
        Validate
      </Link>
      <Link to="/settings" style={{ marginRight: "1rem" }}>
        Settings
      </Link>
      <button onClick={logout}>Logout</button>
    </nav>
  );
}
