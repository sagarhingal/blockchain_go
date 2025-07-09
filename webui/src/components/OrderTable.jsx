import { Table, TableHead, TableRow, TableCell, TableBody } from '@mui/material';
import { Link } from 'react-router-dom';

export default function OrderTable({ orders }) {
  return (
    <Table size="small">
      <TableHead>
        <TableRow>
          <TableCell>ID</TableCell>
          <TableCell>Owner</TableCell>
          <TableCell>Status</TableCell>
        </TableRow>
      </TableHead>
      <TableBody>
        {orders.map((o) => (
          <TableRow key={o.ID}>
            <TableCell>
              <Link to={`/orders/${o.ID}`}>{o.ID}</Link>
            </TableCell>
            <TableCell>{o.Owner}</TableCell>
            <TableCell>{o.Status}</TableCell>
          </TableRow>
        ))}
      </TableBody>
    </Table>
  );
}
