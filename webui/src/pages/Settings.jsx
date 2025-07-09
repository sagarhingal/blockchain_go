import { Typography } from '@mui/material';
import { useAuth } from '../AuthContext';

export default function Settings() {
  const { user } = useAuth();
  return (
    <div style={{ padding: '1rem' }}>
      <Typography variant="h5" gutterBottom>User Settings</Typography>
      <Typography>Logged in as: {user?.name}</Typography>
    </div>
  );
}
