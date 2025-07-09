import { useEffect, useState } from "react";
import { useParams } from "react-router-dom";
import { Typography } from "@mui/material";
import { getOrders, getOrderEvents } from "../api";

export default function OrderDetail() {
  const { id } = useParams();
  const [order, setOrder] = useState(null);
  const [events, setEvents] = useState([]);
  const [error, setError] = useState("");

  useEffect(() => {
    getOrders()
      .then((data) => setOrder(data.find((o) => o.ID === id)))
      .catch((err) => setError(err.message));
    getOrderEvents(id)
      .then(setEvents)
      .catch((err) => setError(err.message));
  }, [id]);

  if (!order) return <Typography sx={{ p: 2 }}>Loading...</Typography>;

  return (
    <div style={{ padding: "1rem" }}>
      <Typography variant="h5" gutterBottom>
        Order {order.ID}
      </Typography>
      {error && <Typography color="error">{error}</Typography>}
      <Typography>Status: {order.Status}</Typography>
      <Typography>Owner: {order.Owner}</Typography>
      <Typography sx={{ mt: 1 }} variant="h6">
        Actors
      </Typography>
      <ul>
        {Object.entries(order.Actors || {}).map(([a, r]) => (
          <li key={a}>
            {a}: {r}
          </li>
        ))}
      </ul>
      <Typography sx={{ mt: 2 }} variant="h6">
        Events
      </Typography>
      <ul>
        {events.map((e, idx) => (
          <li key={idx}>
            {e.time} - {e.actor}: {e.message}
          </li>
        ))}
      </ul>
    </div>
  );
}
