import { useEffect, useState } from 'react';
import { Typography } from '@mui/material';
import { listActors } from '../api';

export default function Marketplace() {
  const [actors, setActors] = useState([]);
  const [error, setError] = useState('');

  useEffect(() => {
    listActors()
      .then(setActors)
      .catch((err) => setError(err.message));
  }, []);

  return (
    <div style={{ padding: '1rem' }}>
      <Typography variant="h5" gutterBottom>Marketplace</Typography>
      {error && <Typography color="error">{error}</Typography>}
      <ul>
        {actors.map((a) => (
          <li key={a.Username}>{a.Username} - {a.FirstName}</li>
        ))}
      </ul>
    </div>
  );
}
