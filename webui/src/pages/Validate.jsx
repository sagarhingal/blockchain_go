import { useState } from 'react';
import { validateChain } from '../api';

export default function Validate() {
  const [result, setResult] = useState(null);
  const [error, setError] = useState('');

  const handleClick = async () => {
    try {
      const res = await validateChain();
      setResult(res.valid);
    } catch (err) {
      setError(err.message);
    }
  };

  return (
    <div style={{ padding: '1rem' }}>
      <h2>Validate Chain</h2>
      {error && <p style={{ color: 'red' }}>{error}</p>}
      <button onClick={handleClick}>Validate</button>
      {result !== null && <p>Chain valid: {String(result)}</p>}
    </div>
  );
}
