import { useEffect, useState } from "react";
import { Typography } from "@mui/material";
import { listActors } from "../api";

export default function Marketplace() {
  return (
    <div style={{ padding: "1rem" }}>
      <Typography variant="h5" gutterBottom>
        Marketplace
      </Typography>
      <Typography>to do</Typography>
    </div>
  );
}
