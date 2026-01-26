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

### Frontend Features

- **Login Page** (`/login`): Public authentication page
- **Admin Dashboard** (`/admin`): Admin-only interface for user management
  - Create new users
  - Delete existing users
  - View recent admin actions
- **User Dashboard** (`/dashboard`): Authenticated user interface
  - File management placeholders (upload/download/delete)
  - Note: Backend file management APIs are not yet implemented
- **Logout** (`/logout`): Clears session and redirects to login

### Authentication & Authorization

- Token-based authentication using localStorage
- Route protection:
  - Unauthenticated users → redirected to `/login`
  - Non-admin users → cannot access `/admin`
  - Admin users → cannot access `/dashboard`
- Auto-logout on 401 responses

### Available API Endpoints

The frontend integrates with these API Gateway endpoints:

- `POST /api/login` - User authentication
- `POST /api/admin/create_user` - Create new user (admin only)
- `DELETE /api/admin/delete_user` - Delete user (admin only)
- `GET /health` - Health check

### Running Tests

Run the frontend tests with:

```bash
cd apps/frontend
npm run test:browser
```

Tests include:
- Login form rendering and validation
- Protected route redirects
- Logout functionality

### Building for Production

```bash
npm run build
npm run start
```

The production build will be served on http://localhost:3000 (React Router serve)