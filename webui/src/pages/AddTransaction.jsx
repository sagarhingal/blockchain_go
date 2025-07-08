import { useState } from "react";
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
      <h2>Add Transaction</h2>
      {error && <p style={{ color: "red" }}>{error}</p>}
      <form onSubmit={handleSubmit} style={{ marginBottom: "1rem" }}>
        <input
          name="from"
          list="names"
          value={form.from}
          onChange={handleChange}
          placeholder="From"
          required
        />
        <input
          name="to"
          list="names"
          value={form.to}
          onChange={handleChange}
          placeholder="To"
          required
        />
        <input
          name="amount"
          list="amounts"
          type="number"
          value={form.amount}
          onChange={handleChange}
          placeholder="Amount"
          required
        />
        <button type="submit">Submit</button>
        <datalist id="names">
          <option value="Alice" />
          <option value="Bob" />
          <option value="Charlie" />
        </datalist>
        <datalist id="amounts">
          <option value="1" />
          <option value="5" />
          <option value="10" />
        </datalist>
      </form>
    </div>
  );
}
