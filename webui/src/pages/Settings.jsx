import { useState } from 'react';
import { useAuth } from '../AuthContext';
import { resetPassword } from '../api';
import { Button, TextField, Typography } from '@mui/material';

export default function Settings() {
  const { user } = useAuth();
  const [pw, setPw] = useState('');
  const [msg, setMsg] = useState('');

  const handleSubmit = async (e) => {
    e.preventDefault();
    try {
      await resetPassword(pw);
      setMsg('Password updated');
      setPw('');
    } catch (err) {
      setMsg(err.message);
    }
  };

  return (
    <div style={{ padding: '1rem' }}>
      <Typography variant="h5" gutterBottom>User Settings</Typography>
      <Typography>Logged in as: {user?.name}</Typography>
      <form onSubmit={handleSubmit} style={{ marginTop: '1rem' }}>
        <TextField
          type="password"
          value={pw}
          onChange={(e) => setPw(e.target.value)}
          label="New Password"
          required
          sx={{ mr: 1 }}
        />
        <Button type="submit" variant="contained">Reset Password</Button>
      </form>
      {msg && <Typography color="primary">{msg}</Typography>}
    </div>
  );
}
