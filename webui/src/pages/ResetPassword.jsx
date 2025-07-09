import { useState } from "react";
import { useNavigate } from "react-router-dom";
import { Button, TextField, Typography } from "@mui/material";
import { resetPassword } from "../api";
import { useAuth } from "../AuthContext";

export default function ResetPassword() {
  const { user } = useAuth();
  const navigate = useNavigate();
  const [email, setEmail] = useState("");
  const [pass, setPass] = useState("");
  const [error, setError] = useState("");

  const handleSubmit = async (e) => {
    e.preventDefault();
    try {
      await resetPassword({ email: email.trim(), password: pass });
      navigate("/login");
    } catch (err) {
      setError(err.message);
    }
  };

  return (
    <form
      onSubmit={handleSubmit}
      style={{ maxWidth: 300, margin: "2rem auto" }}
    >
      <Typography variant="h5" gutterBottom>
        Reset Password
      </Typography>
      {error && <Typography color="error">{error}</Typography>}
      {!user && (
        <TextField
          fullWidth
          value={email}
          onChange={(e) => setEmail(e.target.value)}
          label="Email"
          margin="normal"
          required
        />
      )}
      <TextField
        fullWidth
        type="password"
        value={pass}
        onChange={(e) => setPass(e.target.value)}
        label="New Password"
        margin="normal"
        required
      />
      <Button type="submit" variant="contained" fullWidth sx={{ mt: 2 }}>
        Reset
      </Button>
    </form>
  );
}
