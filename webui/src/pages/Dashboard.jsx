import { useEffect, useState } from "react";
import { getChain } from "../api";

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

  return (
    <div style={{ padding: "1rem" }}>
      <h2>Dashboard</h2>
      {error && <p style={{ color: "red" }}>{error}</p>}
      <p>
        Blocks in chain: {chain.length}{" "}
        <button onClick={refresh}>Refresh</button>
      </p>
      <table style={{ width: "100%", borderCollapse: "collapse" }}>
        <thead>
          <tr>
            <th style={{ borderBottom: "1px solid #ccc", textAlign: "left" }}>
              #
            </th>
            <th style={{ borderBottom: "1px solid #ccc", textAlign: "left" }}>
              From
            </th>
            <th style={{ borderBottom: "1px solid #ccc", textAlign: "left" }}>
              To
            </th>
            <th style={{ borderBottom: "1px solid #ccc", textAlign: "left" }}>
              Amount
            </th>
          </tr>
        </thead>
        <tbody>
          {chain.map((b, idx) => (
            <tr key={idx}>
              <td>{idx}</td>
              <td>{b.Data?.from || "-"}</td>
              <td>{b.Data?.to || "-"}</td>
              <td>{b.Data?.amount || "-"}</td>
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  );
}
