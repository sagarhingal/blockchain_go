import { useAuth } from '../AuthContext';

export default function Settings() {
  const { user } = useAuth();
  return (
    <div style={{ padding: '1rem' }}>
      <h2>User Settings</h2>
      <p>Logged in as: {user?.name}</p>
      <p>More settings could go here.</p>
    </div>
  );
}
