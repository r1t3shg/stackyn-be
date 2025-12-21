# Stackyn Setup Guide

This guide will help you set up Stackyn from scratch.

## Prerequisites

- Linux server (Ubuntu 20.04+ recommended)
- Docker 20.10+
- Docker Compose 2.0+
- Domain name with DNS pointing to your server
- Ports 80, 443, and optionally 8080 open

## Step 1: Server Preparation

```bash
# Update system
sudo apt update && sudo apt upgrade -y

# Install Docker
curl -fsSL https://get.docker.com -o get-docker.sh
sudo sh get-docker.sh

# Install Docker Compose
sudo apt install -y docker-compose-plugin

# Add user to docker group (optional)
sudo usermod -aG docker $USER
newgrp docker

# Verify installation
docker --version
docker compose version
```

## Step 2: Clone and Setup

```bash
# Create directory
sudo mkdir -p /opt/stackyn
sudo chown $USER:$USER /opt/stackyn
cd /opt/stackyn

# Clone repository (or copy files)
git clone <your-repo-url> .

# Create required directories
mkdir -p traefik/acme
mkdir -p volumes/db
mkdir -p volumes/logs
mkdir -p apps
mkdir -p backend/config

# Set permissions
chmod 600 traefik/acme/acme.json
touch traefik/acme/acme.json
chmod 600 traefik/acme/acme.json
```

## Step 3: Configure Backend

```bash
# Create production environment file
cat > backend/config/prod.env << EOF
DATABASE_URL=postgres://stackyn:changeme@postgres:5432/stackyn?sslmode=disable
DOCKER_HOST=unix:///var/run/docker.sock
BASE_DOMAIN=staging.stackyn.com
PORT=8080
EOF

# Edit with your values
nano backend/config/prod.env
```

## Step 4: Configure Traefik

```bash
# Edit Traefik static config
nano traefik/traefik.yml

# Update email for Let's Encrypt (line 30)
# email: your-email@example.com

# Edit dynamic routers
nano traefik/dynamic/routers.yml

# Update domain names:
# - stackyn.local -> your-domain.com
# - staging.stackyn.com -> your-production-domain.com
```

## Step 5: Configure Frontend

The frontend API URL is set at build time in `docker-compose.yml`:

```yaml
frontend:
  build:
    args:
      - VITE_API_BASE_URL=https://api.staging.stackyn.com
```

Update this to match your backend API domain.

## Step 6: Start Services

```bash
# Build and start all services
docker compose up -d --build

# Check status
docker compose ps

# View logs
docker compose logs -f
```

## Step 7: Verify Setup

1. **Check Traefik:**
   ```bash
   docker compose logs traefik
   # Should see "Server configuration reloaded"
   ```

2. **Check Backend:**
   ```bash
   docker compose logs backend
   # Should see "Server started on port 8080"
   ```

3. **Check Frontend:**
   ```bash
   docker compose logs frontend
   # Should see "Accepting connections"
   ```

4. **Access Services:**
   - Frontend: `https://staging.stackyn.com`
   - Backend: `https://api.staging.stackyn.com`
   - Traefik Dashboard: `https://traefik.staging.stackyn.com` (if enabled)

## Step 8: Initial Database Setup

```bash
# Run migrations (if needed)
docker compose exec backend /app/api migrate

# Or if migrations run automatically on startup, check logs
docker compose logs backend | grep -i migration
```

## Troubleshooting

### Port Already in Use

```bash
# Check what's using port 80/443
sudo netstat -tlnp | grep :80
sudo netstat -tlnp | grep :443

# Stop conflicting services
sudo systemctl stop nginx
sudo systemctl stop apache2
```

### SSL Certificate Issues

```bash
# Check Traefik logs
docker compose logs traefik | grep -i acme

# Verify DNS
dig staging.stackyn.com
nslookup staging.stackyn.com

# Check acme.json permissions
ls -la traefik/acme/acme.json
# Should be: -rw------- (600)
```

### Database Connection Issues

```bash
# Check PostgreSQL logs
docker compose logs postgres

# Test connection
docker compose exec postgres psql -U stackyn -d stackyn -c "SELECT 1;"
```

### Frontend Not Loading

```bash
# Check if frontend is built
docker compose exec frontend ls -la /app/dist

# Check Traefik routing
docker compose logs traefik | grep frontend

# Test frontend directly
docker compose exec frontend curl http://localhost:3000
```

## Next Steps

1. **Secure the Setup:**
   - Change default passwords
   - Set up firewall rules
   - Enable Traefik dashboard authentication
   - Review security headers

2. **Monitor:**
   - Set up log aggregation
   - Configure health checks
   - Set up alerts

3. **Backup:**
   - Database backups
   - Volume backups
   - Configuration backups

## Maintenance

### Update Services

```bash
# Pull latest images
docker compose pull

# Rebuild and restart
docker compose up -d --build
```

### Backup Database

```bash
# Create backup
docker compose exec postgres pg_dump -U stackyn stackyn > backup.sql

# Restore backup
docker compose exec -T postgres psql -U stackyn stackyn < backup.sql
```

### View Logs

```bash
# All services
docker compose logs -f

# Specific service
docker compose logs -f backend
docker compose logs -f frontend
docker compose logs -f traefik
```

## Support

For issues and questions:
- Check logs: `docker compose logs`
- Review configuration files
- Check Traefik dashboard
- Review application documentation

