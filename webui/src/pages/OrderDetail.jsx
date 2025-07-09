import { useEffect, useState } from 'react';
import { useParams } from 'react-router-dom';
import {
  Typography,
  TextField,
  Button,
  Alert,
  MenuItem,
} from '@mui/material';
import {
  getOrders,
  getOrderEvents,
  updateOrderStatus,
  addRole,
  inviteWatcher,
  addAddon,
  listActors,
} from '../api';

const roles = [
  'client',
  'supplier',
  'transporter',
  'warehouse',
  'retailer',
];

export default function OrderDetail() {
  const { id } = useParams();
  const [order, setOrder] = useState(null);
  const [events, setEvents] = useState([]);
  const [actors, setActors] = useState([]);
  const [status, setStatus] = useState('');
  const [roleActor, setRoleActor] = useState('');
  const [roleRole, setRoleRole] = useState('client');
  const [watcher, setWatcher] = useState('');
  const [addon, setAddon] = useState('');
  const [msg, setMsg] = useState('');
  const [error, setError] = useState('');

  useEffect(() => {
    getOrders()
      .then((data) => setOrder(data.find((o) => o.ID === id)))
      .catch((err) => setError(err.message));
    getOrderEvents(id)
      .then(setEvents)
      .catch((err) => setError(err.message));
    listActors().then(setActors).catch(() => {});
  }, [id]);

  const handle = async (fn, clear) => {
    try {
      await fn();
      setMsg('Saved');
      setTimeout(() => setMsg(''), 2000);
      clear();
      const updated = await getOrderEvents(id);
      setEvents(updated);
    } catch (err) {
      setError(err.message);
    }
  };

  if (!order) return <Typography sx={{ p: 2 }}>Loading...</Typography>;

  return (
    <div style={{ padding: '1rem' }}>
      <Typography variant="h5" gutterBottom>
        Order {order.ID}
      </Typography>
      {error && <Typography color="error">{error}</Typography>}
      {msg && <Alert severity="success" sx={{ mb: 1 }}>{msg}</Alert>}
      <Typography>Status: {order.Status}</Typography>
      <Typography>Owner: {order.Owner}</Typography>
      <Typography sx={{ mt: 1 }} variant="h6">
        Actors
      </Typography>
      <ul>
        {Object.entries(order.Actors || {}).map(([a, r]) => (
          <li key={a}>{a}: {r}</li>
        ))}
      </ul>
      <form
        onSubmit={(e) => {
          e.preventDefault();
          handle(() => updateOrderStatus(id, status), () => setStatus(''));
        }}
        style={{ marginTop: '1rem' }}
      >
        <TextField
          label="New Status"
          value={status}
          onChange={(e) => setStatus(e.target.value)}
          sx={{ mr: 1 }}
          size="small"
        />
        <Button type="submit" variant="contained" size="small">
          Update
        </Button>
      </form>
      <form
        onSubmit={(e) => {
          e.preventDefault();
          handle(() => addRole(id, roleActor, roleRole), () => setRoleActor(''));
        }}
        style={{ marginTop: '1rem' }}
      >
        <TextField
          select
          label="Actor"
          value={roleActor}
          onChange={(e) => setRoleActor(e.target.value)}
          sx={{ mr: 1 }}
          size="small"
        >
          {actors.map((a) => (
            <MenuItem key={a.Username} value={a.Username}>
              {a.Username}
            </MenuItem>
          ))}
        </TextField>
        <TextField
          select
          label="Role"
          value={roleRole}
          onChange={(e) => setRoleRole(e.target.value)}
          sx={{ mr: 1 }}
          size="small"
        >
          {roles.map((r) => (
            <MenuItem key={r} value={r}>
              {r}
            </MenuItem>
          ))}
        </TextField>
        <Button type="submit" variant="contained" size="small">
          Add Role
        </Button>
      </form>
      <form
        onSubmit={(e) => {
          e.preventDefault();
          handle(() => inviteWatcher(id, watcher), () => setWatcher(''));
        }}
        style={{ marginTop: '1rem' }}
      >
        <TextField
          select
          label="Watcher"
          value={watcher}
          onChange={(e) => setWatcher(e.target.value)}
          sx={{ mr: 1 }}
          size="small"
        >
          {actors.map((a) => (
            <MenuItem key={a.Username} value={a.Username}>
              {a.Username}
            </MenuItem>
          ))}
        </TextField>
        <Button type="submit" variant="contained" size="small">
          Invite
        </Button>
      </form>
      <form
        onSubmit={(e) => {
          e.preventDefault();
          handle(() => addAddon(id, addon), () => setAddon(''));
        }}
        style={{ marginTop: '1rem' }}
      >
        <TextField
          label="Add-on Request"
          value={addon}
          onChange={(e) => setAddon(e.target.value)}
          sx={{ mr: 1 }}
          size="small"
        />
        <Button type="submit" variant="contained" size="small">
          Add
        </Button>
      </form>
      <Typography sx={{ mt: 2 }} variant="h6">
        Events
      </Typography>
      <ul>
        {events.map((e, idx) => (
          <li key={idx}>{e.time} - {e.actor}: {e.message}</li>
        ))}
      </ul>
    </div>
  );
}
