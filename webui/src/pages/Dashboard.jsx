import { useEffect, useState } from "react";
import { Button, Table, TableBody, TableCell, TableHead, TableRow, Typography } from '@mui/material';
import { getChain, validateChain } from "../api";

export default function Dashboard() {
  const [chain, setChain] = useState([]);
  const [error, setError] = useState("");

  useEffect(() => {
    getChain()
      .then((data) => setChain(data.Chain || []))
      .catch((err) => setError(err.message));
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
      alert(`Chain valid: ${res.valid}`);
    } catch (err) {
      setError(err.message);
    }
  };

  return (
    <div style={{ padding: "1rem" }}>
      <Typography variant="h4" gutterBottom>Dashboard</Typography>
      {error && <p style={{ color: "red" }}>{error}</p>}
      <p>
        Blocks in chain: {chain.length}{" "}
        <Button variant="outlined" onClick={refresh}>Refresh</Button>{" "}
        <Button variant="outlined" onClick={validate}>Validate</Button>
      </p>
      <Table>
        <TableHead>
          <TableRow>
            <TableCell>#</TableCell>
            <TableCell>From</TableCell>
            <TableCell>To</TableCell>
            <TableCell>Amount</TableCell>
          </TableRow>
        </TableHead>
        <TableBody>
          {chain.map((b, idx) => (
            <TableRow key={idx}>
              <TableCell>{idx}</TableCell>
              <TableCell>{b.Data?.from || "-"}</TableCell>
              <TableCell>{b.Data?.to || "-"}</TableCell>
              <TableCell>{b.Data?.amount || "-"}</TableCell>
            </TableRow>
          ))}
        </TableBody>
      </Table>
    </div>
  );
}
