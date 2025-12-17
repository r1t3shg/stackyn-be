# MVP Backend - PaaS Deployment Engine

A Go-based backend API and deployment engine for a Platform-as-a-Service (PaaS) MVP.

## Architecture

The backend consists of two main components:

1. **API Server** (`cmd/api`) - REST API for managing apps and deployments
2. **Deployment Worker** (`cmd/worker`) - Background worker that processes deployments

## Features

- **App Management**: Create, list, and delete applications
- **Deployment Pipeline**: 
  - Clone Git repositories
  - Build Docker images
  - Run containers with Traefik labels
  - Track deployment status and logs
- **Database**: PostgreSQL for persistence
- **Docker Integration**: Full Docker SDK integration for building and running containers

## Project Structure

```
backend/
  cmd/
    api/        # HTTP server
    worker/     # deployment worker (engine loop)
  internal/
    config/     # config loading
    db/         # DB connection & migrations
    apps/       # app model + queries
    deployments/# deployment model + queries
    engine/     # core deployment pipeline
    gitrepo/    # cloning repositories
    dockerbuild/# building images
    dockerrun/  # starting containers
    logs/       # deployment logs parsing
  internal/db/migrations/   # SQL migration files
  go.mod
  go.sum
```

## Prerequisites

- Go 1.21 or later
- PostgreSQL
- Docker Engine (accessible via socket or TCP)
- Git (for cloning repositories)
- Traefik (for routing - should be running and watching Docker containers)

## Configuration

The application uses environment variables for configuration:

- `DATABASE_URL` - PostgreSQL connection string (default: `postgres://postgres:postgres@localhost:5432/mvp?sslmode=disable`)
- `DOCKER_HOST` - Docker daemon address (default: `unix:///var/run/docker.sock`)
- `BASE_DOMAIN` - Base domain for subdomain routing (default: `localhost`)
- `PORT` - API server port (default: `8080`)

## Setup

1. **Install dependencies:**
   ```bash
   cd backend
   go mod download
   ```

2. **Set up PostgreSQL:**
   ```bash
   createdb mvp
   # Or use your preferred method to create the database
   ```

3. **Configure environment variables:**
   ```bash
   export DATABASE_URL="postgres://user:password@localhost:5432/mvp?sslmode=disable"
   export DOCKER_HOST="unix:///var/run/docker.sock"
   export BASE_DOMAIN="yourdomain.com"
   export PORT="8080"
   ```

## Running

### API Server

```bash
go run cmd/api/main.go
```

Or build and run:
```bash
go build -o bin/api cmd/api/main.go
./bin/api
```

The API will be available at `http://localhost:8080`

### Deployment Worker

```bash
go run cmd/worker/main.go
```

Or build and run:
```bash
go build -o bin/worker cmd/worker/main.go
./bin/worker
```

**Note:** Both the API server and worker need to be running. The API server handles HTTP requests, while the worker processes deployments in the background.

## API Endpoints

### Apps

- `GET /api/v1/apps` - List all apps
- `POST /api/v1/apps` - Create a new app
  ```json
  {
    "name": "my-app",
    "repo_url": "https://github.com/user/repo.git"
  }
  ```
- `GET /api/v1/apps/{id}` - Get app by ID
- `DELETE /api/v1/apps/{id}` - Delete an app
- `GET /api/v1/apps/{id}/deployments` - List deployments for an app

### Deployments

- `GET /api/v1/deployments/{id}` - Get deployment by ID

### Health Check

- `GET /health` - Health check endpoint

## Deployment Flow

1. **Create App**: POST to `/api/v1/apps` with name and repo URL
2. **Automatic Deployment**: A deployment is automatically created with status `pending`
3. **Worker Processing**: The worker picks up pending deployments and:
   - Clones the repository
   - Builds a Docker image
   - Runs a container with Traefik labels
   - Updates deployment status to `running`
4. **Access**: The app becomes available at `{subdomain}.{BASE_DOMAIN}`

## Traefik Integration

The deployment engine automatically sets Traefik labels on containers:

- `traefik.enable=true`
- `traefik.http.routers.{subdomain}.rule=Host(\`{subdomain}.{baseDomain}\`)`
- `traefik.http.services.{subdomain}.loadbalancer.server.port=80`

Make sure Traefik is configured to watch Docker containers and has access to the Docker socket.

## Database Migrations

Migrations are automatically applied when the application starts. Migration files are located in `internal/db/migrations/`.

## Development

### Building

```bash
go build ./cmd/api
go build ./cmd/worker
```

### Testing

```bash
go test ./...
```

## Notes

- The deployment worker polls for pending deployments every 2 seconds
- Build logs are captured and stored in the database
- Containers are named using the subdomain pattern: `{app-name}-{deployment-id}`
- Images are named: `mvp-{app-name}:{deployment-id}`
- Repository clones are stored in `/tmp/mvp-deployments/` (configurable)

## Future Enhancements

- Webhook support for Git repositories
- Support for custom Dockerfile paths
- Environment variable injection
- Resource limits and constraints
- Multi-stage deployments (staging/production)
- Log streaming via WebSocket
- Container health checks
- Automatic rollback on failure

