import { Link } from "react-router-dom";
import { AppBar, Toolbar, Button } from '@mui/material';
import { useAuth } from "../AuthContext";

export default function NavBar() {
  const { user, logout } = useAuth();

  if (!user) return null;

  return (
    <AppBar position="static">
      <Toolbar>
        <Button color="inherit" component={Link} to="/">Dashboard</Button>
        <Button color="inherit" component={Link} to="/chain">Chain</Button>
        <Button color="inherit" component={Link} to="/transaction">Add Tx</Button>
        <Button color="inherit" component={Link} to="/settings">Settings</Button>
        <Button color="inherit" onClick={logout}>Logout</Button>
      </Toolbar>
    </AppBar>
  );
}
