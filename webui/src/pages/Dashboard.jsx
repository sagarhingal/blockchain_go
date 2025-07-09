import { useEffect, useState } from "react";
import { Typography, IconButton } from '@mui/material';
import RefreshIcon from '@mui/icons-material/Refresh';
import { getOrders } from "../api";
import OrderTable from "../components/OrderTable";

export default function Dashboard() {
  const [orders, setOrders] = useState([]);
  const [error, setError] = useState("");

  useEffect(() => {
    getOrders()
      .then((data) => setOrders(data))
      .catch((err) => setError(err.message));
  }, []);

  const refresh = async () => {
    try {
      const updated = await getOrders(true);
      setOrders(updated);
    } catch (err) {
      setError(err.message);
    }
  };

  return (
    <div style={{ padding: "1rem" }}>
      <Typography variant="h4" gutterBottom>Dashboard</Typography>
      {error && <Typography color="error">{error}</Typography>}
      <div style={{ display: 'flex', alignItems: 'center', marginBottom: '0.5rem' }}>
        <Typography sx={{ mr: 1 }}>Orders: {orders.length}</Typography>
        <IconButton onClick={refresh} aria-label="refresh" size="small"><RefreshIcon /></IconButton>
      </div>
      <OrderTable orders={orders.slice(-5).reverse()} />
    </div>
  );
}
