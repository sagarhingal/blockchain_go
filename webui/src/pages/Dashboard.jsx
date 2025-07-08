import { useEffect, useState } from 'react';
import { getChain, addTransaction } from '../api';

export default function Dashboard() {
  const [chain, setChain] = useState([]);
  const [form, setForm] = useState({ from: '', to: '', amount: '' });
  const [error, setError] = useState('');

  useEffect(() => {
    getChain().then(data => setChain(data.Chain || [])).catch(err => setError(err.message));
  }, []);

  const handleChange = (e) => {
    setForm({ ...form, [e.target.name]: e.target.value });
  };

  const handleSubmit = async (e) => {
    e.preventDefault();
    try {
      await addTransaction({
        from: form.from,
        to: form.to,
        amount: parseFloat(form.amount),
      });
      const updated = await getChain(true);
      setChain(updated.Chain || []);
      setForm({ from: '', to: '', amount: '' });
    } catch (err) {
      setError(err.message);
    }
  };

  return (
    <div style={{ padding: '1rem' }}>
      <h2>Dashboard</h2>
      {error && <p style={{ color: 'red' }}>{error}</p>}
      <p>Blocks in chain: {chain.length}</p>
      <form onSubmit={handleSubmit} style={{ marginBottom: '1rem' }}>
        <input name="from" value={form.from} onChange={handleChange} placeholder="From" required />
        <input name="to" value={form.to} onChange={handleChange} placeholder="To" required />
        <input name="amount" type="number" value={form.amount} onChange={handleChange} placeholder="Amount" required />
        <button type="submit">Add Transaction</button>
      </form>
    </div>
  );
}
