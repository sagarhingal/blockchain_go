# blockchain_go

A minimal blockchain implementation written in Go. It exposes a small HTTP API
and includes an optional React web interface in the `webui` directory.

## Quick start

```bash
go mod tidy
go run ./cmd/server
```

Detailed documentation for the server, CLI and blockchain internals can be
found in the [docs](docs/) directory. A Swagger specification is provided at
`/swagger` when the server is running.
