# HRIS App

A modern, full-stack Human Resource Information System with a React frontend and Go backend.
<img width="2310" height="1044" alt="Screenshot 2026-01-11 115810" src="https://github.com/user-attachments/assets/30082270-5577-44d8-9aec-1d2ad17932a4" />
<img width="2544" height="904" alt="Screenshot 2026-01-11 115822" src="https://github.com/user-attachments/assets/4e0970ab-5d48-4dbb-8730-15ab6508c36d" />
<img width="2527" height="904" alt="Screenshot 2026-01-11 115912" src="https://github.com/user-attachments/assets/0ba8045e-794f-4364-8176-10fd6b56e895" />


## Quick Start

**Prerequisites:** Docker and Docker Compose installed

```bash
# Clone the repository
git clone <repository-url>
cd hris-app

# Copy environment template
cp .env.example .env

# Start all services
docker compose up -d --build
```

**Access the application:**

- Main Application: http://hris.local (add to `/etc/hosts`: `127.0.0.1 hris.local`)
- MinIO Console: http://minio.hris.local (add to `/etc/hosts`: `127.0.0.1 minio.hris.local`)
- Health Check: http://hris.local/api/v1/health

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

- Secure authentication with JWT
- File upload and management
- Database migrations
- Clean architecture
- Docker-based deployment
- Responsive UI with dark mode
- **NGINX reverse proxy** with:
  - Subdomain routing (hris.local, minio.hris.local)
  - Gzip compression
  - WebSocket support for hot reload
  - Load balancing capabilities

## Project Structure

```
hris-app/
├── backend/         # Go backend API
├── frontend/        # React frontend application
├── gateway/         # NGINX reverse proxy configuration
│   └── nginx.conf   # NGINX configuration with routing rules
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
    ├─→ /api/v1/*      → Backend (Go API)
    ├─→ /*             → Frontend (React App)
    └─→ minio.hris.local → MinIO (Object Storage)
```

### Routing Configuration

The gateway is configured in `gateway/nginx.conf`:

1. **Main Application** (`hris.local`):

   - `/api/v1/*` → Proxies to Backend API (port 8081)
   - `/` → Proxies to Frontend (port 8080)
   - Supports WebSocket for React hot reload

2. **MinIO Console** (`minio.hris.local`):
   - `/` → MinIO API (port 9000)
   - `/console` → MinIO Console UI (port 9001)
   - WebSocket support for MinIO console

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
127.0.0.1 hris.local
127.0.0.1 minio.hris.local
127.0.0.1 storage.hris.local
```

For Windows:

```bash
# Run as Administrator
notepad C:\Windows\System32\drivers\etc\hosts

# Add these lines:
127.0.0.1 hris.local
127.0.0.1 console.hris.local
127.0.0.1 storage.hris.local
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

[MIT License](https://github.com/PickHD/hris-app/blob/main/LICENSE)
