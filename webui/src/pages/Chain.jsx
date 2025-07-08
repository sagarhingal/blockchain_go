import { useEffect, useState } from 'react';
import { getChain } from '../api';

export default function Chain() {
  const [chain, setChain] = useState([]);
  const [error, setError] = useState('');

  useEffect(() => {
    getChain().then(data => setChain(data.Chain || [])).catch(err => setError(err.message));
  }, []);

  return (
    <div style={{ padding: '1rem' }}>
      <h2>Blockchain</h2>
      {error && <p style={{ color: 'red' }}>{error}</p>}
      <pre style={{ whiteSpace: 'pre-wrap' }}>{JSON.stringify(chain, null, 2)}</pre>
    </div>
  );
}
