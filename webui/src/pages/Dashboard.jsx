import { useEffect, useState } from "react";
import { Typography, IconButton, Alert } from '@mui/material';
import RefreshIcon from '@mui/icons-material/Refresh';
import CheckIcon from '@mui/icons-material/CheckCircle';
import { getChain, validateChain } from "../api";
import BlockTable from "../components/BlockTable";

export default function Dashboard() {
  const [chain, setChain] = useState([]);
  const [error, setError] = useState("");
  const [message, setMessage] = useState("");

  useEffect(() => {
    getChain()
      .then((data) => setChain(data.Chain || []))
      .catch((err) => setError(err.message));
    const interval = setInterval(() => {
      getChain(true)
        .then((data) => setChain(data.Chain || []))
        .catch(() => {});
    }, 30000);
    return () => clearInterval(interval);
  }, []);

  const refresh = async () => {
    try {
      const updated = await getChain(true);
      setChain(updated.Chain || []);
    } catch (err) {
      setError(err.message);
    }
  };

  const validate = async () => {
    try {
      const res = await validateChain();
      setMessage(`Chain valid: ${res.valid}`);
    } catch (err) {
      setError(err.message);
    }
  };

  useEffect(() => {
    if (!message) return;
    const t = setTimeout(() => setMessage(""), 3000);
    return () => clearTimeout(t);
  }, [message]);

  return (
    <div style={{ padding: "1rem" }}>
      <Typography variant="h4" gutterBottom>Dashboard</Typography>
      {error && <Typography color="error">{error}</Typography>}
      {message && <Alert severity="info" sx={{ mb:1 }}>{message}</Alert>}
      <div style={{ display: 'flex', alignItems: 'center', marginBottom: '0.5rem' }}>
        <Typography sx={{ mr: 1 }}>Blocks in chain: {chain.length}</Typography>
        <IconButton onClick={refresh} aria-label="refresh" size="small"><RefreshIcon /></IconButton>
        <IconButton onClick={validate} aria-label="validate" size="small"><CheckIcon /></IconButton>
      </div>
      <BlockTable chain={chain} />
    </div>
  );
}
