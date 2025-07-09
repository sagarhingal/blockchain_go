import { useState } from 'react';
import { Link, useNavigate } from 'react-router-dom';
import { Button, TextField, Typography } from '@mui/material';
import { useAuth } from '../AuthContext';

export default function SignUp() {
  const { login } = useAuth();
  const navigate = useNavigate();
  const [form, setForm] = useState({
    username: '',
    password: '',
    confirm: '',
    firstName: '',
    lastName: '',
    mobile: '',
    pinCode: '',
    state: '',
    city: '',
    country: '',
  });
  const [error, setError] = useState('');

  const handleChange = (e) => {
    setForm({ ...form, [e.target.name]: e.target.value });
  };

  const handleSubmit = async (e) => {
    e.preventDefault();
    if (form.password !== form.confirm) {
      setError('Passwords do not match');
      return;
    }
    try {
      await login(form.username, form.password, true, {
        first_name: form.firstName,
        last_name: form.lastName,
        mobile: form.mobile,
        pin_code: form.pinCode,
        state: form.state,
        city: form.city,
        country: form.country,
      });
      navigate('/');
    } catch (err) {
      setError(err.message);
    }
  };

  return (
    <form onSubmit={handleSubmit} style={{ maxWidth: 300, margin: '2rem auto' }}>
      <Typography variant="h5" gutterBottom>Sign Up</Typography>
      {error && <Typography color="error">{error}</Typography>}
      <TextField fullWidth name="username" value={form.username} onChange={handleChange} label="Username" margin="normal" required />
      <TextField fullWidth type="password" name="password" value={form.password} onChange={handleChange} label="Password" margin="normal" required />
      <TextField fullWidth type="password" name="confirm" value={form.confirm} onChange={handleChange} label="Confirm Password" margin="normal" required />
      <TextField fullWidth name="firstName" value={form.firstName} onChange={handleChange} label="First Name" margin="normal" />
      <TextField fullWidth name="lastName" value={form.lastName} onChange={handleChange} label="Last Name" margin="normal" />
      <TextField fullWidth name="mobile" value={form.mobile} onChange={handleChange} label="Mobile" margin="normal" />
      <TextField fullWidth name="pinCode" value={form.pinCode} onChange={handleChange} label="Pin Code" margin="normal" />
      <TextField fullWidth name="state" value={form.state} onChange={handleChange} label="State" margin="normal" />
      <TextField fullWidth name="city" value={form.city} onChange={handleChange} label="City" margin="normal" />
      <TextField fullWidth name="country" value={form.country} onChange={handleChange} label="Country" margin="normal" />
      <Button type="submit" variant="contained" fullWidth sx={{ mt:2 }}>Create Account</Button>
      <Button component={Link} to="/login" fullWidth sx={{ mt:1 }}>Back to Login</Button>
    </form>
  );
}
