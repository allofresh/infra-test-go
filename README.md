# Infra Test Go

A sample Go REST API for product management, built with Go 1.22 standard library.

## Prerequisites

- [Go 1.22+](https://go.dev/dl/)
- [Docker](https://docs.docker.com/get-docker/) (optional, for containerized run)
- `make` (included on macOS/Linux)

## Project Structure

```
infra-test-go/
├── cmd/api/main.go                    # Application entrypoint
├── internal/
│   ├── model/product.go               # Product model
│   ├── service/product_service.go     # Business logic (thread-safe)
│   └── handler/product_handler.go     # HTTP handlers (REST API)
├── Dockerfile                         # Multi-stage Docker build
├── Makefile                           # Build & test targets
└── go.mod
```

## Running the Application

### Local

```bash
go run ./cmd/api
```

The server starts on port `8080` by default. Override with the `PORT` environment variable:

```bash
PORT=3000 go run ./cmd/api
```

### Docker

```bash
# Build the image
make docker-build

# Run the container
make docker-run
```

Or manually:

```bash
docker build -t infra-test-go:latest .
docker run -p 8080:8080 infra-test-go:latest
```

## API Endpoints

| Method   | Endpoint         | Description          |
|----------|------------------|----------------------|
| `GET`    | `/health`        | Health check         |
| `GET`    | `/products`      | List all products    |
| `GET`    | `/products/{id}` | Get product by ID    |
| `POST`   | `/products`      | Add a new product    |
| `DELETE` | `/products/{id}` | Delete a product     |

### Example Requests

```bash
# Health check
curl http://localhost:8080/health

# Add a product
curl -X POST http://localhost:8080/products \
  -H "Content-Type: application/json" \
  -d '{"name":"Apple","price":1.50,"quantity":100}'

# List products
curl http://localhost:8080/products

# Get product by ID
curl http://localhost:8080/products/1

# Delete a product
curl -X DELETE http://localhost:8080/products/1
```

## Running Unit Tests

The test pipeline has 4 components matching `run_unit_test_components`:

```json
["vet", "unit-test", "race-test", "coverage"]
```

### Run All Tests (Recommended)

```bash
make all
```

This executes all 4 components in order: `vet` → `unit-test` → `race-test` → `coverage`.

### Run Individual Components

```bash
# 1. vet — static analysis
make vet

# 2. unit-test — run all unit tests
make unit-test

# 3. race-test — run tests with race detector
make race-test

# 4. coverage — run tests with coverage report (minimum 80%)
make coverage
```

### Run Tests with Go Directly

```bash
# Basic test run
go test ./...

# Verbose output
go test -v ./...

# With race detector
go test -race ./...

# With coverage
go test -cover ./...

# Generate coverage HTML report
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out -o coverage.html
```

## Build

```bash
# Build binary
make build

# Run the binary
./bin/infra-test-go
```

## Clean

```bash
make clean
```

Removes the `bin/` and `coverage/` directories.
