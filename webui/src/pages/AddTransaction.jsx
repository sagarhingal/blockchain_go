import { useState } from "react";
import { Button, TextField, Typography } from '@mui/material';
import { addTransaction } from "../api";

export default function AddTransaction() {
  const [form, setForm] = useState({ from: "", to: "", amount: "" });
  const [error, setError] = useState("");
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
      setForm({ from: "", to: "", amount: "" });
      setError("");
    } catch (err) {
      setError(err.message);
    }
  };

  return (
    <div style={{ padding: "1rem" }}>
      <Typography variant="h5" gutterBottom>Add Transaction</Typography>
      {error && <Typography color="error">{error}</Typography>}
      <form onSubmit={handleSubmit} style={{ marginBottom: "1rem" }}>
        <TextField
          name="from"
          value={form.from}
          onChange={handleChange}
          label="From"
          required
          sx={{ mr: 1 }}
        />
        <TextField
          name="to"
          value={form.to}
          onChange={handleChange}
          label="To"
          required
          sx={{ mr: 1 }}
        />
        <TextField
          name="amount"
          type="number"
          value={form.amount}
          onChange={handleChange}
          label="Amount"
          required
          sx={{ mr: 1, width: 100 }}
        />
        <Button type="submit" variant="contained">Submit</Button>
      </form>
    </div>
  );
}
