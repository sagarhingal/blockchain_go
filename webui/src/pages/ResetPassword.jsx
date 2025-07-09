import { useState } from 'react';
import { useNavigate, Link } from 'react-router-dom';
import { Button, TextField, Typography } from '@mui/material';
import { resetPassword } from '../api';

export default function ResetPassword() {
  const [pw, setPw] = useState('');
  const [confirm, setConfirm] = useState('');
  const [msg, setMsg] = useState('');
  const navigate = useNavigate();

  const handleSubmit = async (e) => {
    e.preventDefault();
    if (pw !== confirm) {
      setMsg('Passwords do not match');
      return;
    }
    try {
      await resetPassword(pw);
      setMsg('Password updated');
      setTimeout(() => navigate('/login'), 1000);
    } catch (err) {
      setMsg(err.message);
    }
  };

  return (
    <form onSubmit={handleSubmit} style={{ maxWidth:300, margin:'2rem auto' }}>
      <Typography variant="h5" gutterBottom>Reset Password</Typography>
      {msg && <Typography>{msg}</Typography>}
      <TextField fullWidth type="password" value={pw} onChange={(e)=>setPw(e.target.value)} label="New Password" margin="normal" required />
      <TextField fullWidth type="password" value={confirm} onChange={(e)=>setConfirm(e.target.value)} label="Confirm Password" margin="normal" required />
      <Button type="submit" variant="contained" fullWidth sx={{ mt:2 }}>Reset</Button>
      <Button component={Link} to="/login" fullWidth sx={{ mt:1 }}>Back</Button>
    </form>
  );
}
