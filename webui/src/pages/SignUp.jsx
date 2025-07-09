import { useState } from 'react';
import { Link, useNavigate } from 'react-router-dom';
import { Button, TextField, Typography } from '@mui/material';
import { useAuth } from '../AuthContext';

export default function SignUp() {
  const { login } = useAuth();
  const navigate = useNavigate();
  const [name, setName] = useState('');
  const [pass, setPass] = useState('');
  const [confirm, setConfirm] = useState('');
  const [error, setError] = useState('');

  const handleSubmit = async (e) => {
    e.preventDefault();
    if (pass !== confirm) {
      setError('Passwords do not match');
      return;
    }
    try {
      await login(name, pass, true);
      navigate('/');
    } catch (err) {
      setError(err.message);
    }
  };

  return (
    <form onSubmit={handleSubmit} style={{ maxWidth: 300, margin: '2rem auto' }}>
      <Typography variant="h5" gutterBottom>Sign Up</Typography>
      {error && <Typography color="error">{error}</Typography>}
      <TextField fullWidth value={name} onChange={(e)=>setName(e.target.value)} label="Username" margin="normal" required />
      <TextField fullWidth type="password" value={pass} onChange={(e)=>setPass(e.target.value)} label="Password" margin="normal" required />
      <TextField fullWidth type="password" value={confirm} onChange={(e)=>setConfirm(e.target.value)} label="Confirm Password" margin="normal" required />
      <Button type="submit" variant="contained" fullWidth sx={{ mt:2 }}>Create Account</Button>
      <Button component={Link} to="/login" fullWidth sx={{ mt:1 }}>Back to Login</Button>
    </form>
  );
}
