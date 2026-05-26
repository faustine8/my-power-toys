# my-power-toys

A tiny cross-platform project launcher for developers.

## Status

This project is in early development (skeleton stage).

## Quick Start

```bash
# Show help
go run .

# Print version
go run . version

# Run tests
go test ./...
```

## Build

```bash
go build -o mpt .
./mpt version
```

## Project Structure

```text
my-power-toys/
├── main.go              # Entry point
├── cmd/
│   └── root.go          # Cobra CLI commands
├── internal/
│   ├── project/         # Project management (TODO)
│   ├── opener/          # Opener execution (TODO)
│   ├── config/          # Configuration (TODO)
│   └── tui/             # Bubble Tea TUI (TODO)
├── AGENTS.md            # Rules for AI agents
├── go.mod
└── README.md
```

## License

MIT
