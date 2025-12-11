# Vector Backend

Go backend service for the Vector monorepo project.

## 📋 Overview

This is the backend service built with Go, providing API endpoints and business logic for the Vector application.

## 🛠️ Tech Stack

- **Language**: Go (Golang)
- **Architecture**: Clean Architecture with `cmd` and `internal` packages
- **Linting**: golangci-lint with comprehensive rules
- **Package Manager**: Go Modules

## 📁 Project Structure

```
backend/
├── cmd/
│   └── vector/          # Application entry point
├── internal/            # Private application code
├── .env                 # Environment variables (not in git)
├── .env.sample          # Environment variables template
├── .gitignore           # Git ignore rules
├── .golangci.yml        # Linting configuration
├── go.mod               # Go module dependencies
├── go.sum               # Dependency checksums
└── README.md            # This file
```

## 🚀 Getting Started

### Prerequisites

- Go 1.21 or higher
- golangci-lint (for linting)

### Installation

1. **Install Go dependencies:**
   ```bash
   go mod download
   ```

2. **Set up environment variables:**
   ```bash
   cp .env.sample .env
   # Edit .env with your configuration
   ```

### Development

**Run the application:**
```bash
go run cmd/vector/main.go
```

**Build the application:**
```bash
go build -o bin/vector cmd/vector/main.go
```

**Run the built binary:**
```bash
./bin/vector
```

## 🧪 Testing

**Run all tests:**
```bash
go test ./...
```

**Run tests with coverage:**
```bash
go test -cover ./...
```

**Generate coverage report:**
```bash
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

## 🔍 Code Quality

### Linting

This project uses [golangci-lint](https://golangci-lint.run/) with a comprehensive configuration.

**Run linter:**
```bash
golangci-lint run
```

**Auto-fix issues:**
```bash
golangci-lint run --fix
```

**Install golangci-lint:**
```bash
# macOS
brew install golangci-lint

# Linux/Windows
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
```

### Formatting

**Format code:**
```bash
go fmt ./...
```

**Organize imports:**
```bash
goimports -w .
```

## 📦 Building for Production

**Build optimized binary:**
```bash
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o bin/vector cmd/vector/main.go
```

**Build flags:**
- `CGO_ENABLED=0` - Disable CGO for static binary
- `-ldflags="-w -s"` - Strip debug information for smaller binary size
- `GOOS` and `GOARCH` - Target operating system and architecture

## 🔧 Configuration

Configuration is managed through environment variables. See `.env.sample` for available options.

### Environment Variables

Create a `.env` file based on `.env.sample`:

```bash
cp .env.sample .env
```

## 📝 Development Guidelines

### Code Organization

- **`cmd/`** - Application entry points (main packages)
- **`internal/`** - Private application code (cannot be imported by other projects)
  - `internal/api` - HTTP handlers and routes
  - `internal/service` - Business logic
  - `internal/repository` - Data access layer
  - `internal/models` - Data structures

### Best Practices

1. **Error Handling**: Always check and handle errors appropriately
2. **Context Usage**: Pass `context.Context` for cancellation and timeouts
3. **Logging**: Use structured logging for better observability
4. **Testing**: Write unit tests for business logic
5. **Documentation**: Add comments for exported functions and types

## 🐛 Debugging

**Run with race detector:**
```bash
go run -race cmd/vector/main.go
```

**Enable verbose logging:**
```bash
LOG_LEVEL=debug go run cmd/vector/main.go
```

## 📚 Additional Resources

- [Go Documentation](https://go.dev/doc/)
- [Effective Go](https://go.dev/doc/effective_go)
- [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)
- [golangci-lint Linters](https://golangci-lint.run/usage/linters/)

## 🤝 Contributing

1. Follow the existing code style and structure
2. Run linters and tests before committing
3. Write meaningful commit messages
4. Keep functions small and focused
5. Document exported functions and types

## 📄 License

ISC
