# Web UI

This directory contains a small React application built with [Vite](https://vitejs.dev/) that provides a basic interface for the Go blockchain server.

## Features

- Login/logout flow with an in-memory session
- Dashboard to view chain length and submit new transactions
- View the raw blockchain data
- Validate the chain state
- Simple user settings page
- Minimal caching for the chain request to reduce network traffic

## Development

1. Install dependencies:
   ```bash
   npm install
   ```
2. Start the development server:
   ```bash
   npm run dev
   ```
   The app will be available at `http://localhost:5173` by default.

The UI expects the blockchain server to be running on `http://localhost:8080`. If your server runs elsewhere, set the `VITE_API_URL` environment variable when starting the dev server:

```bash
VITE_API_URL=http://otherhost:8080 npm run dev
```

## Building for production

```bash
npm run build
```

This will output static files in the `dist` directory which can be served by any HTTP server.
