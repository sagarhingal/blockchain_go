# Server

This document describes the HTTP API exposed by the Go server.

## Running

```bash
go mod tidy
go run ./cmd/server
```

## Endpoints

- `POST /signup` – register a new user
- `POST /login` – authenticate and receive a JWT token
- `POST /order` – create a new order (requires Authorization header)
- `POST /order/{id}/roles` – assign a role to another actor
- `POST /order/{id}/status` – update the status of an order
- `POST /order/{id}/invite` – invite another user as watcher
- `POST /order/{id}/addon` – record an add‑on request
- `GET  /order/{id}/events` – retrieve decrypted events
- `GET  /chain` – return all orders
- `GET  /marketplace` – list all registered users

See `docs/swagger/swagger.yaml` for the OpenAPI specification.
