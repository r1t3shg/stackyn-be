# Stackyn - Platform as a Service

Stackyn is a self-hosted PaaS platform that allows you to deploy applications from Git repositories.

## Project Structure

```
/opt/stackyn
├── traefik/                 # Infrastructure (entrypoint)
│   ├── traefik.yml          # Static config
│   ├── dynamic/             # File provider configs
│   │   ├── routers.yml
│   │   ├── middlewares.yml
│   │   └── services.yml
│   └── acme/                # Let's Encrypt certs
│       └── acme.json
│
├── frontend/                # React (Vite) UI
│   ├── Dockerfile
│   ├── src/
│   └── dist/                # Built static files (generated)
│
├── backend/                 # Go API
│   ├── Dockerfile
│   ├── config/
│   │   └── prod.env
│   └── data/                # App data (if any)
│
├── apps/                    # USER DEPLOYED APPS (core of Stackyn)
│   ├── app-001/
│   │   ├── docker-compose.yml
│   │   └── env/
│   ├── app-002/
│   └── README.md
│
├── volumes/                 # Shared persistent volumes
│   ├── db/
│   ├── redis/
│   └── logs/
│
├── docker-compose.yml       # Root orchestrator
└── README.md
```

## Quick Start

### Prerequisites

- Docker and Docker Compose installed
- Domain name pointing to your server (for SSL)
- Ports 80 and 443 open

### Setup

1. **Clone the repository:**
   ```bash
   git clone <your-repo-url> /opt/stackyn
   cd /opt/stackyn
   ```

2. **Configure environment variables:**
   ```bash
   # Edit backend config
   cp backend/config/prod.env.example backend/config/prod.env
   nano backend/config/prod.env
   ```

3. **Create Traefik ACME directory:**
   ```bash
   mkdir -p traefik/acme
   touch traefik/acme/acme.json
   chmod 600 traefik/acme/acme.json
   ```

4. **Update Traefik configuration:**
   - Edit `traefik/traefik.yml` and set your email for Let's Encrypt
   - Update domain names in `traefik/dynamic/routers.yml`

5. **Start the stack:**
   ```bash
   docker-compose up -d
   ```

6. **Access the services:**
   - Frontend: `https://staging.stackyn.com`
   - Backend API: `https://api.staging.stackyn.com`
   - Traefik Dashboard: `https://traefik.staging.stackyn.com` (development only)

## Configuration

### Traefik

- **Static Config**: `traefik/traefik.yml` - Main configuration
- **Dynamic Config**: `traefik/dynamic/` - Runtime routing rules
- **SSL Certificates**: `traefik/acme/acme.json` - Let's Encrypt certificates

### Backend

- **Environment**: `backend/config/prod.env`
- **Required Variables**:
  - `DATABASE_URL`: PostgreSQL connection string
  - `BASE_DOMAIN`: Base domain for deployed apps
  - `DOCKER_HOST`: Docker socket path

### Frontend

- **Build Args**: Set `VITE_API_BASE_URL` in `docker-compose.yml`
- **Environment**: Uses build-time variables

## Deploying User Apps

User apps are deployed in the `apps/` directory. Each app gets its own subdirectory with:
- `docker-compose.yml`: App-specific configuration
- `env/`: Environment variables
- App data and volumes

## Volumes

Persistent data is stored in `volumes/`:
- `db/`: PostgreSQL data
- `logs/`: Application logs
- `redis/`: Redis data (if used)

## Development

### Local Development

1. **Frontend:**
   ```bash
   cd frontend
   npm install
   npm run dev
   ```

2. **Backend:**
   ```bash
   cd backend
   go run cmd/api/main.go
   ```

### Building

```bash
# Build all services
docker-compose build

# Build specific service
docker-compose build frontend
docker-compose build backend
```

## Troubleshooting

### Check Logs

```bash
# All services
docker-compose logs -f

# Specific service
docker-compose logs -f frontend
docker-compose logs -f backend
docker-compose logs -f traefik
```

### Restart Services

```bash
# Restart all
docker-compose restart

# Restart specific service
docker-compose restart frontend
```

### SSL Certificate Issues

- Ensure `traefik/acme/acme.json` has correct permissions (600)
- Check Traefik logs for ACME errors
- Verify domain DNS is pointing to your server

## Security Notes

- Change default passwords in `backend/config/prod.env`
- Secure Traefik dashboard in production
- Use strong database passwords
- Regularly update Docker images
- Review Traefik middleware configurations

## License

[Your License Here]

