import { useEffect, useState } from 'react';
import { Typography, Button, Alert } from '@mui/material';
import { getOrders, createOrder } from '../api';
import OrderTable from '../components/OrderTable';

export default function Orders() {
  const [orders, setOrders] = useState([]);
  const [error, setError] = useState('');
  const [msg, setMsg] = useState('');

  const fetchData = () => {
    getOrders(true)
      .then((data) => setOrders(data))
      .catch((err) => setError(err.message));
  };

  useEffect(() => { fetchData(); }, []);

  const create = async () => {
    try {
      const ord = await createOrder();
      setMsg(`Created order ${ord.ID}`);
      fetchData();
      setTimeout(() => setMsg(''), 3000);
    } catch (err) {
      setError(err.message);
    }
  };

  return (
    <div style={{ padding: '1rem' }}>
      <Typography variant="h5" gutterBottom>Orders</Typography>
      {error && <Typography color="error">{error}</Typography>}
      {msg && <Alert severity="success" sx={{ mb:1 }}>{msg}</Alert>}
      <Button variant="contained" onClick={create} sx={{ mb: 1 }}>New Order</Button>
      <OrderTable orders={orders} />
    </div>
  );
}
