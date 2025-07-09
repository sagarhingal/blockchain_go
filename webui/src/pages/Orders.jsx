import { useEffect, useState } from "react";
import { Typography } from "@mui/material";
import { getOrders } from "../api";
import OrderTable from "../components/OrderTable";

export default function Orders() {
  const [orders, setOrders] = useState([]);
  const [error, setError] = useState("");

  const fetchData = () => {
    getOrders(true)
      .then((data) => setOrders(data))
      .catch((err) => setError(err.message));
  };

  useEffect(() => {
    fetchData();
  }, []);

  return (
    <div style={{ padding: "1rem" }}>
      <Typography variant="h5" gutterBottom>
        Orders
      </Typography>
      {error && <Typography color="error">{error}</Typography>}
      <OrderTable orders={orders} />
    </div>
  );
}
