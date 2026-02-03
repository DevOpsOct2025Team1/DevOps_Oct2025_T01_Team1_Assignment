# Devops Assignment
This assignment project uses a microservices architecture.
## Services

- **User Service** (Port 8080): Manages user data with MongoDB
- **Auth Service** (Port 8081): Handles authentication with JWT
- **API Gateway** (Port 3000): HTTP/REST interface for clients
- **Frontend** (Port 5173): React application for user interface

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

## Frontend Development

### Running the Frontend Locally

The frontend is a React 19 application built with React Router v7 and Tailwind CSS v4.

1. **Navigate to the frontend directory:**
   ```bash
   cd apps/frontend
   ```

2. **Install dependencies (if not already done):**
   ```bash
   npm install
   ```

3. **Start the development server:**
   ```bash
   npm run dev
   ```

   The frontend will be available at http://localhost:5173

4. **Configure the API endpoint (optional):**
   
   By default, the frontend connects to the API Gateway at `http://localhost:3001`. To change this:
   
   - Create a `.env` file in `apps/frontend/`:
     ```
     VITE_API_BASE_URL=http://localhost:3001
     ```
   
   - Or set the environment variable before running:
     ```bash
     VITE_API_BASE_URL=http://your-api-url npm run dev
     ```
