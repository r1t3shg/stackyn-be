# Migration from Nginx to Traefik

This document summarizes the migration from Nginx to Traefik for the Stackyn frontend.

## Changes Made

### Frontend Dockerfile
- **Before**: Used nginx:alpine to serve static files
- **After**: Uses `serve` (Node.js static file server) on port 3000
- **Reason**: Traefik handles routing and SSL, frontend just needs to serve files

### Removed Files
- `frontend/nginx.conf` - No longer needed
- `frontend/docker-compose.yml` - Using root docker-compose.yml instead

### New Traefik Configuration

#### Static Configuration (`traefik/traefik.yml`)
- Entry points: HTTP (80) and HTTPS (443)
- Providers: Docker and File
- ACME: Let's Encrypt SSL certificates
- Dashboard: Enabled (insecure for development)

#### Dynamic Configuration (`traefik/dynamic/`)
- **routers.yml**: Routing rules for frontend and backend
- **middlewares.yml**: Security headers, CORS
- **services.yml**: Backend service definitions

### Root Docker Compose

Created `docker-compose.yml` at root level that orchestrates:
- Traefik (reverse proxy)
- Frontend (React app)
- Backend (Go API)
- PostgreSQL (database)

## Benefits of Traefik

1. **Automatic SSL**: Let's Encrypt certificates auto-generated and renewed
2. **Service Discovery**: Automatically discovers Docker containers
3. **Dynamic Configuration**: Update routing without restarting
4. **Unified Management**: Single entry point for all services
5. **Better for Microservices**: Designed for containerized applications

## Configuration

### Frontend Routing
- Domain: `staging.stackyn.com` or `stackyn.local`
- Port: 3000 (internal)
- SSL: Automatic via Let's Encrypt

### Backend Routing
- Domain: `api.staging.stackyn.com` or `api.stackyn.local`
- Port: 8080 (internal)
- SSL: Automatic via Let's Encrypt

### Environment Variables

Frontend build-time variable:
```yaml
VITE_API_BASE_URL=https://api.staging.stackyn.com
```

Set in `docker-compose.yml` build args.

## Deployment

1. **Update Traefik config:**
   - Set your email in `traefik/traefik.yml`
   - Update domains in `traefik/dynamic/routers.yml`

2. **Create ACME directory:**
   ```bash
   mkdir -p traefik/acme
   touch traefik/acme/acme.json
   chmod 600 traefik/acme/acme.json
   ```

3. **Start services:**
   ```bash
   docker compose up -d --build
   ```

4. **Verify:**
   - Frontend: `https://staging.stackyn.com`
   - Backend: `https://api.staging.stackyn.com`
   - Traefik Dashboard: `http://localhost:8080` (development)

## Troubleshooting

### SSL Certificate Issues
- Check `traefik/acme/acme.json` permissions (must be 600)
- Verify DNS is pointing to your server
- Check Traefik logs: `docker compose logs traefik`

### Frontend Not Loading
- Check if frontend is running: `docker compose ps frontend`
- Check Traefik routing: `docker compose logs traefik | grep frontend`
- Verify build: `docker compose exec frontend ls -la /app/dist`

### Backend Not Accessible
- Check backend logs: `docker compose logs backend`
- Verify Traefik labels in docker-compose.yml
- Test backend directly: `docker compose exec backend curl http://localhost:8080/health`

## Next Steps

1. **Production Hardening:**
   - Disable Traefik dashboard or add authentication
   - Set `api.dashboard.insecure: false`
   - Review security headers in middlewares

2. **Monitoring:**
   - Set up log aggregation
   - Monitor SSL certificate expiration
   - Track service health

3. **Scaling:**
   - Traefik supports load balancing
   - Can add multiple backend instances
   - Horizontal scaling ready

