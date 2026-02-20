# BaseKarya Backend

A REST API built with Go following Clean Architecture principles for the Human Resource Information System.

## Tech Stack

- **Go 1.25.1** - Primary programming language
- **Echo v4** - High-performance web framework
- **MySQL 8.0** - Relational database
- **GORM** - ORM for database operations
- **MinIO** - S3-compatible object storage
- **Zap** - Structured logging
- **golang-migrate** - Database migration tool
- **NGINX** - Reverse proxy (see root README for gateway configuration)

## Architecture

This backend follows **Clean Architecture** principles with clear separation of concerns:

```
├── cmd/api/              # Application entry point
├── internal/
│   ├── bootstrap/        # Dependency injection
│   ├── config/           # Configuration management
│   ├── infrastructure/   # External services (MySQL, MinIO)
│   ├── modules/          # Business logic by domain
│   │   └── health/       # Example: health check module
│   │       ├── handler.go    # HTTP layer
│   │       ├── service.go    # Business logic
│   │       └── repository.go # Data access
│   └── routes/           # HTTP routing
├── pkg/                  # Public/reusable packages
└── migrations/           # Database migrations
```

### Layer Responsibilities

- **Handler Layer**: Handles HTTP requests/responses
- **Service Layer**: Business logic and orchestration
- **Repository Layer**: Database operations and data persistence

## Quick Start

### Using Docker Compose (Recommended)

From the root directory:

```bash
# Copy environment template
cp .env.example .env

# Start all services
docker compose up -d --build

# Check logs
docker compose logs -f backend
```

The API will be available at `http://localhost:8081` (internal)
Note: When using Docker Compose, access the API through the NGINX gateway at `http://basekarya.local/api/v1`

### Local Development

```bash
cd backend

# Install dependencies
go mod download

# Run database migrations
migrate -path ./migrations -database "mysql://user:password@tcp(localhost:3306)/basekarya_db" up

# Run the server
go run cmd/api/main.go
```

The API will be available at `http://localhost:8080`

## Adding a New Module

1. Create directory: `internal/modules/{module-name}/`
2. Implement three files:
   - `handler.go` - HTTP layer
   - `service.go` - Business logic
   - `repository.go` - Data access
3. Register routes in `internal/routes/api.go`
4. Wire dependencies in `internal/bootstrap/container.go`

## Environment Variables

Key variables (see root `.env.example` for complete list):

| Variable | Description | Default |
|----------|-------------|---------|
| `SERVER_PORT` | API server port | 8080 |
| `SERVER_ENV` | Environment | development |
| `MYSQL_HOST` | Database host | localhost |
| `MYSQL_PORT` | Database port | 3306 |
| `MYSQL_DATABASE` | Database name | basekarya_db |
| `JWT_SECRET` | JWT signing key | - |
| `JWT_EXPIRES_IN` | JWT expiration | 24h |
| `LOG_LEVEL` | Logging level | debug |
| `MINIO_ENDPOINT` | MinIO endpoint | - |
| `MINIO_BUCKET_NAME` | Default bucket | - |

## API Testing

### When using Docker Compose (with NGINX Gateway)

```bash
# Health check through gateway
curl http://basekarya.local/api/v1/health

# Expected response
{
  "messages": "OK",
  "data": true,
  "error": null
}
```

### When running locally

```bash
# Health check direct to backend
curl http://localhost:8080/health

# Expected response
{
  "messages": "OK",
  "data": true,
  "error": null
}
```

## Development Workflow

```bash
# Run tests
go test ./...

# Format code
go fmt ./...

# Run linter
go vet ./...

# Build binary
go build -o basekarya-be-service ./cmd/api
```

## Troubleshooting

**Database connection failed**
- Check MySQL is running: `docker compose ps db`
- Verify `.env` credentials
- Check logs: `docker compose logs backend`

**Migration errors**
- Verify database connectivity
- Check migration file syntax
- Review logs: `docker compose logs migrate`

**Port already in use**
- Change `SERVER_PORT` in `.env`
- Check what's using the port: `lsof -i :8080`

## License

See LICENSE file in the root directory.
