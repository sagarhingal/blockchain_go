# blockchain_go

This project provides a minimal blockchain implementation written in Go.
It now exposes an HTTP server for interacting with the blockchain and
demonstrates how simple smart contracts could be integrated.

## Running the server

```
# install dependencies
go mod tidy
# run the server
go run ./cmd/server
```

## Running the CLI

```
go run ./cmd/cli
```

The server exposes three endpoints:

- `POST /transaction` – add a new transaction block using JSON payload `{"from":"Alice","to":"Bob","amount":5}`.
- `GET /chain` – retrieve the full chain.
- `GET /validate` – verify the chain integrity.

Example usage with `curl`:

```bash
curl -X POST -d '{"from":"Alice","to":"Bob","amount":5}' \
  -H 'Content-Type: application/json' http://localhost:8080/transaction
curl http://localhost:8080/chain
curl http://localhost:8080/validate
```

## Package layout

- `internal/blockchain` – blockchain types and logic.
- `internal/contracts` – basic smart contract interface.
- `cmd/server` – HTTP server exposing blockchain APIs.
- `cmd/cli` – demonstration CLI application.

## Testing

Run unit tests with:

```
go test ./...
```

