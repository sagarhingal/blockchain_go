import { useState } from 'react';
import { Link } from 'react-router-dom';
import { Button, TextField, Typography } from '@mui/material';
import { useAuth } from '../AuthContext';

export default function Login() {
  const { login } = useAuth();
  const [name, setName] = useState('');
  const [pass, setPass] = useState('');
  const [error, setError] = useState('');

  const handleSubmit = async (e) => {
    e.preventDefault();
    try {
      await login(name.trim(), pass);
    } catch (err) {
      setError(err.message);
    }
  };

  return (
    <form onSubmit={handleSubmit} style={{ maxWidth: 300, margin: '2rem auto' }}>
      <Typography variant="h5" gutterBottom>Login</Typography>
      {error && <Typography color="error">{error}</Typography>}
      <TextField
        fullWidth
        value={name}
        onChange={(e) => setName(e.target.value)}
        label="Username"
        margin="normal"
        required
      />
      <TextField
        fullWidth
        type="password"
        value={pass}
        onChange={(e) => setPass(e.target.value)}
        label="Password"
        margin="normal"
        required
      />
      <Button type="submit" variant="contained" fullWidth sx={{ mt: 2 }}>Login</Button>
      <Button component={Link} to="/signup" fullWidth sx={{ mt: 1 }}>Sign Up</Button>
      <Button component={Link} to="/reset" fullWidth sx={{ mt: 1 }}>Reset Password</Button>
    </form>
  );
}
