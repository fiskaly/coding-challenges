# Signature Service - Go Backend

This is the Go backend implementation for the signature service challenge.

## Getting Started

### Prerequisites

- Go 1.21+
- Docker & Docker Compose (for PostgreSQL)

### Running with Docker

```bash
docker-compose up
```

The backend will be available at `http://localhost:8080`.

### Running Locally

1. Start PostgreSQL:
```bash
docker-compose up postgres
```

2. Run the application:
```bash
go run main.go
```

### Database

The PostgreSQL schema is automatically initialized when you start the containers.
You can find the schema in `db/schema.sql`.

## Project Structure

```
backend-go/
├── api/              # HTTP handlers and routing
├── crypto/           # Cryptographic operations (RSA, ECC)
├── domain/           # Domain models
├── persistence/      # Database access layer
├── db/               # Database schemas and migrations
├── main.go           # Application entry point
└── docker-compose.yml
```

## API Endpoints

### Health Check
```
GET /api/v1/health
```

### TODO: Implement these endpoints

```
POST   /api/v1/devices
GET    /api/v1/devices
POST   /api/v1/transactions
GET    /api/v1/transactions
```

## Testing

```bash
go test ./...
```
