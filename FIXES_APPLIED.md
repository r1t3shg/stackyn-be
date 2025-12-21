# Fixes Applied

## Issues Fixed

### 1. Database Authentication Failure
**Problem**: Backend was using wrong database credentials (`user:password` instead of `stackyn:changeme`)

**Fix**: Updated `docker-compose.yml` to use correct credentials matching PostgreSQL service:
```yaml
DATABASE_URL=postgres://stackyn:changeme@postgres:5432/stackyn?sslmode=disable
```

### 2. Traefik ACME File Read-Only
**Problem**: Traefik couldn't write SSL certificates because acme volume was mounted read-only

**Fix**: Changed volume mount from `:ro` to `:rw`:
```yaml
- ./traefik/acme:/etc/traefik/acme:rw
```

### 3. Traefik Docker API Version Mismatch
**Problem**: Traefik v3.0 has Docker API compatibility issues with older Docker clients

**Fix**: Downgraded to Traefik v2.11 which is more stable and compatible:
```yaml
image: traefik:v2.11
```

## Next Steps

1. **Restart services:**
   ```bash
   docker compose down
   docker compose up -d
   ```

2. **Check backend logs:**
   ```bash
   docker compose logs -f backend
   ```
   Should see: "Server starting on port 8080"

3. **Test health endpoint:**
   ```bash
   # Direct test
   docker compose exec backend curl http://localhost:8080/health
   
   # Via Traefik (may need to wait for SSL)
   curl -k https://api.staging.stackyn.com/health
   # Or HTTP
   curl http://api.staging.stackyn.com/health
   ```

4. **Check Traefik logs:**
   ```bash
   docker compose logs traefik
   ```
   Should not see Docker API version errors

5. **Verify database connection:**
   ```bash
   docker compose exec postgres psql -U stackyn -d stackyn -c "SELECT 1;"
   ```

## SSL Certificate Note

If you still see self-signed certificate errors:
- Traefik is generating Let's Encrypt certificates (can take a few minutes)
- For testing, use `curl -k` to skip SSL verification
- Or test via HTTP: `curl http://api.staging.stackyn.com/health`
- Ensure DNS is pointing to your server for Let's Encrypt to work

## If Backend Still Restarts

1. **Check if prod.env exists:**
   ```bash
   ls -la backend/config/prod.env
   ```

2. **Create it if missing:**
   ```bash
   cp backend/config/prod.env.example backend/config/prod.env
   ```

3. **Verify database is ready:**
   ```bash
   docker compose ps postgres
   docker compose logs postgres
   ```

4. **Wait for database to be healthy:**
   ```bash
   # Check health
   docker compose exec postgres pg_isready -U stackyn
   ```

