import { useEffect, useState } from "react";
import {
  Typography,
  IconButton,
  TextField,
  TablePagination,
} from "@mui/material";
import RefreshIcon from "@mui/icons-material/Refresh";
import SortIcon from "@mui/icons-material/Sort";
import { getOrders } from "../api";
import OrderTable from "../components/OrderTable";

export default function Dashboard() {
  const [orders, setOrders] = useState([]);
  const [error, setError] = useState("");
  const [search, setSearch] = useState("");
  const [page, setPage] = useState(0);
  const [rowsPerPage, setRowsPerPage] = useState(5);
  const [asc, setAsc] = useState(false);

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

  const filtered = orders.filter(
    (o) =>
      o.ID.includes(search) ||
      o.Owner.toLowerCase().includes(search.toLowerCase()),
  );
  const sorted = [...filtered].sort((a, b) => {
    const t1 = new Date(a.Created).getTime();
    const t2 = new Date(b.Created).getTime();
    return asc ? t1 - t2 : t2 - t1;
  });
  const paginated = sorted.slice(
    page * rowsPerPage,
    page * rowsPerPage + rowsPerPage,
  );

  return (
    <div style={{ padding: "1rem" }}>
      <Typography variant="h4" gutterBottom>
        Dashboard
      </Typography>
      {error && <Typography color="error">{error}</Typography>}
      <div
        style={{
          display: "flex",
          alignItems: "center",
          marginBottom: "0.5rem",
        }}
      >
        <Typography sx={{ mr: 1 }}>Orders: {orders.length}</Typography>
        <IconButton onClick={refresh} aria-label="refresh" size="small">
          <RefreshIcon />
        </IconButton>
        <IconButton onClick={() => setAsc(!asc)} size="small">
          <SortIcon style={{ transform: asc ? "rotate(180deg)" : "none" }} />
        </IconButton>
        <TextField
          size="small"
          value={search}
          onChange={(e) => {
            setSearch(e.target.value);
            setPage(0);
          }}
          placeholder="search"
          sx={{ ml: 1 }}
        />
      </div>
      <OrderTable orders={paginated} />
      <TablePagination
        component="div"
        count={sorted.length}
        page={page}
        onPageChange={(_, p) => setPage(p)}
        rowsPerPage={rowsPerPage}
        onRowsPerPageChange={(e) => {
          setRowsPerPage(parseInt(e.target.value, 10));
          setPage(0);
        }}
      />
    </div>
  );
}
