# MCP-GO-DEMO

A minimal MCP server using the official [`modelcontextprotocol/go-sdk`](https://github.com/modelcontextprotocol/go-sdk).

## Run

```bash
go run .
```

The server communicates over stdio and exposes a single tool:

- `ping` → returns `{"message":"pong"}`

## Test

```bash
go test ./...
```
