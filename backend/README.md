# HRIS Backend

A REST API built with Go following Clean Architecture principles for Human Resource Information System (HRIS).

## Tech Stack

The backend is built using the following technologies:

### Core Framework & Language
- **Go 1.25.1** - Primary programming language
- **Echo v4** - High-performance, minimalist Go web framework for routing and middleware

### Database & ORM
- **MySQL 8.0** - Relational database management system
- **GORM** - Go's ORM (Object-Relational Mapping) library for database operations
- **MySQL Driver for GORM** - Database driver for MySQL connections

### Object Storage
- **MinIO** - High-performance, S3-compatible object storage for file management
- **MinIO Go SDK v7** - Go client library for MinIO operations

### Logging & Configuration
- **Zap (Uber)** - Fast, structured, leveled logging library
- **godotenv** - Load environment variables from .env file

### Containerization & Deployment
- **Docker** - Containerization platform
- **Docker Compose** - Multi-container orchestration
- **Alpine Linux** - Minimal Docker image for production deployment

### Database Migration
- **golang-migrate/migrate** - Database migration tool for managing schema changes

## Folder Structure & Clean Architecture Pattern

This backend follows **Clean Architecture** principles, ensuring separation of concerns, testability, and maintainability.

### Architecture Layers

```
backend/
├── cmd/                      # Application entry points
│   └── api/
│       └── main.go          # HTTP server entry point
│
├── internal/                 # Private application code
│   ├── bootstrap/           # Dependency injection & container setup
│   │   └── container.go     # Wire dependencies and initialize components
│   │
│   ├── config/              # Configuration management
│   │   └── config.go        # Environment-based configuration loader
│   │
│   ├── infrastructure/      # External services integration
│   │   ├── mysql.go         # MySQL/GORM connection provider
│   │   └── minio.go         # MinIO storage provider
│   │
│   ├── middleware/          # Echo middleware components (empty, ready for use)
│   │
│   ├── modules/             # Business logic organized by domain
│   │   └── health/          # Health check module example
│   │       ├── handler.go   # HTTP layer - handles requests/responses
│   │       ├── service.go   # Business logic layer
│   │       └── repository.go# Data access layer
│   │
│   └── routes/              # HTTP routing configuration
│       └── api.go           # Route definitions and middleware setup
│
├── pkg/                     # Public/reusable packages
│   ├── logger/              # Logging utilities
│   │   └── logger.go        # Zap logger wrapper
│   └── response/            # HTTP response utilities
│       └── http.go          # Standardized response structures
│
├── migrations/              # Database migration files
│
├── Dockerfile               # Container image definition
├── go.mod                   # Go module definition
└── go.sum                   # Dependency checksums
```

### Clean Architecture Principles

#### 1. **Dependency Rule**
Dependencies flow inward, with the innermost layers containing business logic and outer layers containing infrastructure details.

#### 2. **Layer Separation**

**Handler Layer (Presentation)**
- Location: `internal/modules/{module}/handler.go`
- Responsibility: Handle HTTP requests/responses
- Depends on: Service layer interfaces
- Example: `health/handler.go:18-24`

**Service Layer (Business Logic)**
- Location: `internal/modules/{module}/service.go`
- Responsibility: Implement business rules and orchestrate operations
- Depends on: Repository layer interfaces
- Example: `health/service.go:15-17`

**Repository Layer (Data Access)**
- Location: `internal/modules/{module}/repository.go`
- Responsibility: Database operations and data persistence
- Depends on: Infrastructure (GORM)
- Example: `health/repository.go:21-35`

#### 3. **Dependency Injection**
- Implemented via `bootstrap/container.go`
- Creates and wires all dependencies
- Provides clean dependency management and testability
- Example: `bootstrap/container.go:17-33`

#### 4. **Interface-Based Design**
- Services and Repositories are defined as interfaces
- Enables easy mocking for testing
- Promotes loose coupling between layers
- Example: `health/service.go:3-5` and `health/repository.go:9-11`

#### 5. **Configuration Management**
- Environment-based configuration
- Centralized in `internal/config/config.go`
- Supports default values and environment variable overrides
- Example: `config/config.go:53-89`

## How to Run the Server/API

### Prerequisites

- Go 1.25.1 or higher
- Docker and Docker Compose
- Git

### Method 1: Using Docker Compose (Recommended)

This is the easiest way to run the complete application stack including database and storage.

1. **Clone the repository**
   ```bash
   git clone <repository-url>
   cd hris-app
   ```

2. **Configure environment variables**
   ```bash
   cp .env.example .env
   ```

   Edit the `.env` file with your configuration:
   ```env
   # Database Configuration
   MYSQL_ROOT_PASSWORD=root_password
   MYSQL_DATABASE=hris_db
   MYSQL_USER=hris_user
   MYSQL_PASSWORD=hris_password
   MYSQL_HOST=db
   MYSQL_PORT=3306
   MYSQL_SSLMODE=disable

   # Backend Configuration
   SERVER_ENV=development
   SERVER_PORT=8080
   JWT_SECRET=your-super-secret-jwt-key
   JWT_EXPIRES_IN=24h

   # MinIO Configuration
   MINIO_ROOT_USER=minioadmin
   MINIO_ROOT_PASSWORD=minioadmin
   MINIO_ENDPOINT=minio:9000
   MINIO_ACCESS_KEY=minioadmin
   MINIO_SECRET_KEY=minioadmin
   MINIO_BUCKET_NAME=hris-uploads
   MINIO_BUCKET_LOCATION=us-east-1
   MINIO_IS_SECURE=false
   MINIO_PORT=9000
   MINIO_PORT_ADMIN=9001

   # Logging Configuration
   LOG_LEVEL=debug

   # File Upload Configuration
   MAX_REQUEST_BODY_SIZE_MB=50
   MAX_FILE_SIZE_MB=40

   # Frontend Configuration
   VITE_PORT=5173
   VITE_API_URL=http://localhost:8080
   ```

3. **Start all services**
   ```bash
   docker-compose up --build
   ```

   This will start:
   - MySQL database (port 3306)
   - Database migration service
   - MinIO object storage (ports 9000, 9001)
   - Backend API (port 8080)
   - Frontend (port 5173)

4. **Access the API**
   - API Base URL: `http://localhost:8080`
   - Health Check: `http://localhost:8080/health`
   - MinIO Console: `http://localhost:9001`

5. **Stop services**
   ```bash
   docker-compose down
   ```

   To remove volumes as well:
   ```bash
   docker-compose down -v
   ```

### Method 2: Local Development (Without Docker)

For local development without Docker containers.

1. **Install dependencies**
   ```bash
   cd backend
   go mod download
   ```

2. **Set up MySQL database**
   - Install MySQL 8.0 locally
   - Create database: `CREATE DATABASE hris_db;`
   - Update `.env` with your MySQL credentials

3. **Set up MinIO (Optional)**
   - Install MinIO locally or use a cloud S3 service
   - Update `.env` with MinIO credentials
   - Skip file upload features if not needed

4. **Run database migrations**
   ```bash
   # Install golang-migrate if not installed
   # https://github.com/golang-migrate/migrate/blob/master/cmd/migrate/README.md

   migrate -path ./migrations -database "mysql://user:password@tcp(localhost:3306)/hris_db" up
   ```

5. **Run the application**
   ```bash
   go run cmd/api/main.go
   ```

   Or build and run:
   ```bash
   go build -o hris-be-service ./cmd/api
   ./hris-be-service
   ```

6. **Access the API**
   - API Base URL: `http://localhost:8080`
   - Health Check: `http://localhost:8080/health`

### API Testing

Test the API using curl, Postman, or any HTTP client:

```bash
# Health Check
curl http://localhost:8080/health

# Expected Response
{
  "messages": "OK",
  "data": true,
  "error": null
}
```

### Development Workflow

1. **Adding a New Module**
   - Create a new directory under `internal/modules/{module-name}/`
   - Implement `handler.go`, `service.go`, and `repository.go`
   - Register the handler in `internal/routes/api.go`
   - Wire dependencies in `internal/bootstrap/container.go`

2. **Database Migrations**
   - Create migration files in `migrations/` directory
   - Use naming convention: `{version}_{name}.up.sql` and `{version}_{name}.down.sql`
   - Run migrations using the migrate service or tool

3. **Running Tests**
   ```bash
   go test ./...
   ```

### Environment Variables Reference

| Variable | Description | Default |
|----------|-------------|---------|
| `MYSQL_HOST` | Database host | localhost |
| `MYSQL_PORT` | Database port | 3306 |
| `MYSQL_USER` | Database user | root |
| `MYSQL_PASSWORD` | Database password | root_password |
| `MYSQL_DATABASE` | Database name | hris_db |
| `SERVER_PORT` | API server port | 8080 |
| `SERVER_ENV` | Environment (development/production) | development |
| `JWT_SECRET` | JWT signing secret | your-super-secret-jwt-key |
| `JWT_EXPIRES_IN` | JWT expiration time | 24h |
| `LOG_LEVEL` | Logging level | debug |
| `MINIO_ENDPOINT` | MinIO endpoint URL | |
| `MINIO_ACCESS_KEY` | MinIO access key | |
| `MINIO_SECRET_KEY` | MinIO secret key | |
| `MINIO_BUCKET_NAME` | MinIO bucket name | |
| `MAX_FILE_SIZE_MB` | Maximum file upload size | 40 |

### Troubleshooting

**Database Connection Issues**
- Ensure MySQL container is running: `docker-compose ps`
- Check database logs: `docker-compose logs db`
- Verify environment variables in `.env`

**Port Conflicts**
- Change ports in `.env` if default ports are in use
- Ensure no other services are using the same ports

**Migration Failures**
- Check database connectivity
- Verify migration file syntax
- Review migration logs: `docker-compose logs migrate`

### Production Deployment

For production deployment:

1. Set `SERVER_ENV=production` in environment variables
2. Use strong, unique passwords for all services
3. Enable SSL/TLS for database connections (`MYSQL_SSLMODE=require`)
4. Use secure MinIO configuration (`MINIO_IS_SECURE=true`)
5. Set appropriate `LOG_LEVEL` (e.g., `info` or `error`)
6. Configure proper CORS and security middleware
7. Use a reverse proxy (nginx, traefik) for SSL termination
8. Implement proper health checks and monitoring

### License

See LICENSE file in the root directory.
