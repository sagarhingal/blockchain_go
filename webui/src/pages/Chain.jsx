import { useEffect, useState } from 'react';
import { Typography } from '@mui/material';
import { getChain } from '../api';

export default function Chain() {
  const [chain, setChain] = useState([]);
  const [error, setError] = useState('');

  useEffect(() => {
    getChain().then(data => setChain(data.Chain || [])).catch(err => setError(err.message));
  }, []);

  return (
    <div style={{ padding: '1rem' }}>
      <Typography variant="h5" gutterBottom>Blockchain</Typography>
      {error && <Typography color="error">{error}</Typography>}
      <pre style={{ whiteSpace: 'pre-wrap' }}>{JSON.stringify(chain, null, 2)}</pre>
    </div>
  );
}
