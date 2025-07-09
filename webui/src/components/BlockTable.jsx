import { useState, useMemo } from 'react';
import {
  Table, TableHead, TableRow, TableCell, TableBody,
  TablePagination, TextField, TableSortLabel
} from '@mui/material';

export default function BlockTable({ chain, showControls = true }) {
  const [search, setSearch] = useState('');
  const [page, setPage] = useState(0);
  const [rowsPerPage, setRowsPerPage] = useState(10);
  const [orderBy, setOrderBy] = useState('index');
  const [order, setOrder] = useState('desc');

  const filtered = useMemo(() => {
    const q = search.toLowerCase();
    return chain
      .map((b, idx) => ({ ...b, index: idx }))
      .filter(b => {
        const data = b.Data || {};
        const text = `${data.from || ''} ${data.to || ''} ${data.amount || ''} ${b.Hash}`.toLowerCase();
        return text.includes(q);
      });
  }, [chain, search]);

  const sorted = useMemo(() => {
    const arr = [...filtered];
    const dir = order === 'asc' ? 1 : -1;
    arr.sort((a, b) => {
      if (orderBy === 'timestamp') return dir * (a.Timestamp - b.Timestamp);
      return dir * (a.index - b.index);
    });
    return arr;
  }, [filtered, order, orderBy]);

  const handleChangePage = (e, newPage) => setPage(newPage);
  const handleChangeRowsPerPage = e => { setRowsPerPage(parseInt(e.target.value, 10)); setPage(0); };

  const start = page * rowsPerPage;
  const visible = sorted.slice(start, start + rowsPerPage);

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
            <TableCell>
              <TableSortLabel
                active={orderBy === 'index'}
                direction={orderBy === 'index' ? order : 'asc'}
                onClick={() => {
                  const isAsc = orderBy === 'index' && order === 'asc';
                  setOrder(isAsc ? 'desc' : 'asc');
                  setOrderBy('index');
                }}
              >
                #
              </TableSortLabel>
            </TableCell>
            <TableCell>Hash</TableCell>
            <TableCell>Prev</TableCell>
            <TableCell>From</TableCell>
            <TableCell>To</TableCell>
            <TableCell align="right">Amount</TableCell>
            <TableCell>
              <TableSortLabel
                active={orderBy === 'timestamp'}
                direction={orderBy === 'timestamp' ? order : 'asc'}
                onClick={() => {
                  const isAsc = orderBy === 'timestamp' && order === 'asc';
                  setOrder(isAsc ? 'desc' : 'asc');
                  setOrderBy('timestamp');
                }}
              >
                Timestamp
              </TableSortLabel>
            </TableCell>
          </TableRow>
        </TableHead>
        <TableBody>
          {visible.map((b) => {
            const idx = b.index;
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
          count={sorted.length}
          rowsPerPage={rowsPerPage}
          page={page}
          onPageChange={handleChangePage}
          onRowsPerPageChange={handleChangeRowsPerPage}
        />
      )}
    </>
  );
}
