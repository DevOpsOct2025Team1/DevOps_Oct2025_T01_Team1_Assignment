# Devops Assignment
This assignment project uses a microservices architecture.
## Services

- **User Service** (Port 8080): Manages user data with MongoDB
- **Auth Service** (Port 8081): Handles authentication with JWT
- **API Gateway** (Port 3000): HTTP/REST interface for clients

## Prerequisites

### For Development (Docker Compose - Recommended)
- [Docker Desktop](https://www.docker.com/products/docker-desktop/) or Docker + Docker Compose

## Quick Start

### Setup env
```bash
make setup-envw
```

### Serve all services
```bash
make serve
```

This starts all services via NX in development mode.

### Docker Compose
To run with Docker Compose instead:

Ensure you have created your .env file, then run:
```bash
docker compose up
```

This starts:
- API Gateway: http://localhost:3000

**Stop everything:**
```bash
docker compose down
```

**View logs:**
```bash
docker compose logs -f
```

### Run services individually:
```bash
./nx serve user-service

./nx serve auth-service

./nx serve api-gateway
```