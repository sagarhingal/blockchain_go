import { useState, useMemo } from 'react';
import {
  Table, TableHead, TableRow, TableCell, TableBody,
  TablePagination, TextField
} from '@mui/material';

export default function BlockTable({ chain, showControls = true }) {
  const [search, setSearch] = useState('');
  const [page, setPage] = useState(0);
  const [rowsPerPage, setRowsPerPage] = useState(10);

  const filtered = useMemo(() => {
    const q = search.toLowerCase();
    return chain.filter(b => {
      const data = b.Data || {};
      const text = `${data.from || ''} ${data.to || ''} ${data.amount || ''} ${b.Hash}`.toLowerCase();
      return text.includes(q);
    });
  }, [chain, search]);

  const handleChangePage = (e, newPage) => setPage(newPage);
  const handleChangeRowsPerPage = e => { setRowsPerPage(parseInt(e.target.value, 10)); setPage(0); };

  const start = page * rowsPerPage;
  const visible = filtered.slice(start, start + rowsPerPage);

  return (
    <>
      {showControls && (
        <TextField
          size="small"
          label="Search"
          value={search}
          onChange={e => setSearch(e.target.value)}
          sx={{ mb: 1 }}
        />
      )}
      <Table size="small">
        <TableHead>
          <TableRow>
            <TableCell>#</TableCell>
            <TableCell>Hash</TableCell>
            <TableCell>Prev</TableCell>
            <TableCell>From</TableCell>
            <TableCell>To</TableCell>
            <TableCell align="right">Amount</TableCell>
            <TableCell>Timestamp</TableCell>
          </TableRow>
        </TableHead>
        <TableBody>
          {visible.map((b, i) => {
            const idx = start + i;
            const isGenesis = idx === 0;
            const data = b.Data || {};
            return (
              <TableRow key={idx} sx={isGenesis ? { bgcolor: '#f5f5f5' } : undefined}>
                <TableCell>{idx}</TableCell>
                <TableCell>{b.Hash ? `${b.Hash.slice(0,8)}...` : '-'}</TableCell>
                <TableCell>{b.PrevHash ? `${b.PrevHash.slice(0,8)}...` : '-'}</TableCell>
                <TableCell>{isGenesis ? '-' : data.from}</TableCell>
                <TableCell>{isGenesis ? '-' : data.to}</TableCell>
                <TableCell align="right">{isGenesis ? '-' : data.amount}</TableCell>
                <TableCell>{new Date(b.Timestamp).toLocaleString()}</TableCell>
              </TableRow>
            );
          })}
        </TableBody>
      </Table>
      {showControls && (
        <TablePagination
          rowsPerPageOptions={[5,10,25]}
          component="div"
          count={filtered.length}
          rowsPerPage={rowsPerPage}
          page={page}
          onPageChange={handleChangePage}
          onRowsPerPageChange={handleChangeRowsPerPage}
        />
      )}
    </>
  );
}
