# BaseKarya

A modern, full-stack Human Resource Information System with a React frontend and Go backend.

![CI](https://github.com/PickHD/basekarya/actions/workflows/ci.yml/badge.svg)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://github.com/PickHD/basekarya/blob/main/LICENSE)
[![Maintainability](https://qlty.sh/gh/PickHD/projects/basekarya/maintainability.svg)](https://qlty.sh/gh/PickHD/projects/basekarya)

<img width="2428" height="1202" alt="image" src="https://github.com/user-attachments/assets/2c54e865-2c6d-47cc-8cc1-881a688a22eb" />

## Quick Start

**Prerequisites:** Docker and Docker Compose installed

```bash
# Clone the repository
git clone <repository-url>
cd basekarya

# Copy environment template
cp .env.example .env

# Start all services
docker compose up -d --build
```

**Access the application:**

- Main Application: http://basekarya.local (add to `/etc/hosts`: `127.0.0.1 basekarya.local`)
- API Health Check: http://basekarya.local/api/v1/health

**Stop services:**

```bash
docker compose down
```

## Tech Stack

### Frontend

- React 19.2.0 + TypeScript 5.9.3
- Vite 7.2.4
- TailwindCSS 3.4.17
- TanStack Query
- Radix UI

### Backend

- Go 1.25.1
- Echo v4
- MySQL 8.0 + GORM
- MinIO (S3-compatible storage)
- Zap logging

### Infrastructure

- **NGINX** - Reverse proxy and load balancer
- **Docker Compose** - Container orchestration
- **MinIO** - S3-compatible object storage

## Features

- Secure authentication with JWT & role-based access control (RBAC)
- Multi-tenant SaaS architecture with subscription-based module gating
- Comprehensive Employee Management (Admin)
- Attendance tracking with face recognition and geolocation
- Company profile & organizational configuration
- Real-time Notifications via WebSockets
- Automated Payroll generation and email delivery
- Employee Reimbursement tracking & approval workflow
- Leave management with auto balance generation
- Overtime tracking with configurable rates
- Loan/kasbon management with installment deductions
- Finance & accounting (transactions, categories, dashboard)
- Recruitment pipeline (requisitions, applicants)
- Employee onboarding workflows & templates
- Contract management with expiry notifications
- File upload and management (MinIO S3-compatible)
- Database migrations
- Clean architecture (handler/service/repository)
- Docker-based deployment
- Responsive UI with dark mode
- **NGINX reverse proxy** with:
  - Subdomain routing (basekarya.local, storage.basekarya.local)
  - Gzip compression
  - WebSocket support for hot reload
  - Load balancing capabilities

## Testing

### Backend (Go)

```bash
make test-be             # Run all tests
cd backend && go test ./... -cover -coverprofile=coverage.out   # With coverage
```

### Frontend (React + Vitest)

```bash
make test-fe             # Run all tests
cd frontend && pnpm test:coverage   # With coverage
```

### Test Coverage

| Module   | Tests | Description                                                                  |
| -------- | ----- | ---------------------------------------------------------------------------- |
| Backend  | 410+  | Unit tests across all modules (auth, user, payroll, attendance, leave, etc.) |
| Frontend | 188   | Hook and utility tests across 25 feature modules                             |

## CI/CD

GitHub Actions runs on every push/PR to `main` or `develop`:

- **Backend**: Build + test + coverage report (`go test -cover`)
- **Frontend**: Install + test + coverage report (Vitest with v8 coverage)
- **Dependabot**: Weekly dependency update checks for Go, npm, Docker, and GitHub Actions

Coverage reports appear in the **Actions > Summary** tab. Artifacts are uploaded for detailed HTML/LCOV reports.

## Project Structure

```
basekarya/
├── backend/         # Go backend API
├── frontend/        # React frontend application
├── gateway/         # NGINX reverse proxy configuration
│   └── nginx.conf   # NGINX configuration with routing rules
├── .github/
│   ├── workflows/
│   │   ├── ci.yml       # CI pipeline (test + coverage)
│   │   └── dependabot.yml  # Weekly dependency updates
├── docker-compose.yml
├── Makefile
└── .env.example
```

## Available Commands

```bash
make help          # Show all commands
make run-docker    # Run with Docker
make build         # Build both services
make run           # Run both locally
make run-be        # Run backend only
make run-fe        # Run frontend only
make migrate-up    # Run database migrations
make test          # Run all tests (backend + frontend)
make test-be       # Run backend tests
make test-fe       # Run frontend tests
```

## Local Development

**Backend:**

```bash
cd backend
go mod download
go run cmd/api/main.go
```

**Frontend:**

```bash
cd frontend
pnpm install
pnpm dev
```

## NGINX Gateway

This project uses **NGINX as a reverse proxy** to route traffic between the frontend and backend services.

### Architecture

```
Internet (Port 80/443)
    ↓
[NGINX Gateway]
    ├─→ basekarya.local/api/v1/*      → Backend (Go API)
    ├─→ basekarya.local/*             → Frontend (React App)
    └─→ storage.basekarya.local/ → MinIO (Object Storage)
```

### Routing Configuration

The gateway is configured in `gateway/nginx.conf`:

1. **Main Application** (`basekarya.local`):
   - `/api/v1/*` → Proxies to Backend API (port 8081)
   - `/` → Proxies to Frontend (port 8080)
   - Supports WebSocket for React hot reload

2. **MinIO API S3** (`storage.basekarya.local`):
   - `/` → MinIO API S3 (port 9000)

### Features

- **Gzip Compression**: Compresses text-based responses (JSON, CSS, JS, HTML)
- **WebSocket Support**: Enables hot reload during development and MinIO console
- **Health Checks**: Backend health check at `/api/v1/health`
- **Performance Optimizations**:
  - Sendfile enabled
  - TCP optimizations (nopush, nodelay)
  - Keep-alive connections
- **File Upload**: Supports up to 100MB file uploads

### Setup Local Hosts

To access the application locally, add these entries to your `/etc/hosts` file:

```bash
# Linux/macOS
sudo nano /etc/hosts

# Add these lines:
127.0.0.1 basekarya.local
127.0.0.1 storage.basekarya.local
```

For Windows:

```bash
# Run as Administrator
notepad C:\Windows\System32\drivers\etc\hosts

# Add these lines:
127.0.0.1 basekarya.local
127.0.0.1 storage.basekarya.local
```

### Customizing NGINX Configuration

To modify the gateway configuration:

1. Edit `gateway/nginx.conf`
2. Restart the gateway service:
   ```bash
   docker compose restart gateway
   ```

### SSL/HTTPS Setup (Optional)

The configuration includes commented-out volumes for Let's Encrypt certificates. To enable HTTPS:

1. Uncomment the certbot volumes in `docker-compose.yml`:

   ```yaml
   volumes:
     - ./gateway/nginx.conf:/etc/nginx/nginx.conf:ro
     - ./certbot/conf:/etc/letsencrypt
     - ./certbot/www:/var/www/certbot
   ```

2. Update `gateway/nginx.conf` to include SSL configuration

3. Use Certbot to generate certificates automatically

## Documentation

- [Backend Documentation](./backend/README.md) - Architecture, API, and development guide
- [Frontend Documentation](./frontend/README.md) - Components, styling, and setup

## License

[MIT License](https://github.com/PickHD/basekarya/blob/main/LICENSE)
