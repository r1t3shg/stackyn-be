# Testing Guide for Stackyn

This guide explains how to test all components of Stackyn.

## Service URLs

Based on Traefik configuration:

- **Frontend**: `https://staging.stackyn.com`
- **Backend API**: `https://api.staging.stackyn.com`
- **Traefik Dashboard**: `http://localhost:8080` (development only)

## Health Check Endpoint

The health endpoint is on the **backend API**, not the frontend:

✅ **Correct**: `https://api.staging.stackyn.com/health`  
❌ **Wrong**: `https://staging.stackyn.com/health`

## Testing Methods

### 1. Check if Services are Running

```bash
# Check all services
docker compose ps

# Check specific service
docker compose ps backend
docker compose ps frontend
docker compose ps traefik
docker compose ps postgres
```

### 2. Check Service Logs

```bash
# All services
docker compose logs -f

# Specific service
docker compose logs -f backend
docker compose logs -f frontend
docker compose logs -f traefik
docker compose logs -f postgres
```

### 3. Test Backend API Directly (Inside Container)

```bash
# Test health endpoint directly
docker compose exec backend curl http://localhost:8080/health

# Test API endpoint
docker compose exec backend curl http://localhost:8080/api/v1/apps
```

### 4. Test from Host Machine

#### Using curl

```bash
# Health check (should use api subdomain)
curl https://api.staging.stackyn.com/health

# If SSL not ready, test HTTP first
curl http://api.staging.stackyn.com/health

# Or test via IP (if you know it)
curl http://YOUR_VPS_IP/health
```

#### Using Browser

- **Frontend**: `https://staging.stackyn.com`
- **Backend Health**: `https://api.staging.stackyn.com/health`
- **Backend API**: `https://api.staging.stackyn.com/api/v1/apps`

### 5. Test Traefik Routing

```bash
# Check Traefik logs for routing
docker compose logs traefik | grep -i "backend\|frontend\|routing"

# Check Traefik dashboard (if enabled)
# Access: http://localhost:8080/dashboard/
```

### 6. Test Database Connection

```bash
# Test PostgreSQL connection
docker compose exec postgres psql -U stackyn -d stackyn -c "SELECT 1;"

# Check if backend can connect
docker compose logs backend | grep -i "database\|postgres"
```

## Common Issues and Solutions

### Issue: Health endpoint not accessible

**Symptoms**: `https://api.staging.stackyn.com/health` returns connection error

**Solutions**:

1. **Check if backend is running:**
   ```bash
   docker compose ps backend
   # Should show "Up"
   ```

2. **Check backend logs:**
   ```bash
   docker compose logs backend
   # Look for "Server starting on port 8080"
   ```

3. **Test backend directly:**
   ```bash
   docker compose exec backend curl http://localhost:8080/health
   # Should return: {"status":"ok"}
   ```

4. **Check Traefik routing:**
   ```bash
   docker compose logs traefik | grep backend
   # Should show backend service registered
   ```

5. **Verify DNS:**
   ```bash
   # Check if DNS is pointing to your server
   nslookup api.staging.stackyn.com
   dig api.staging.stackyn.com
   ```

6. **Check SSL certificates:**
   ```bash
   # Check if certificates are being generated
   docker compose logs traefik | grep -i "acme\|certificate"
   
   # Check acme.json
   ls -la traefik/acme/acme.json
   ```

### Issue: SSL Certificate Errors

**Symptoms**: Browser shows SSL certificate error

**Solutions**:

1. **Check Let's Encrypt logs:**
   ```bash
   docker compose logs traefik | grep -i "acme\|letsencrypt"
   ```

2. **Verify DNS is correct:**
   ```bash
   dig staging.stackyn.com
   dig api.staging.stackyn.com
   ```

3. **Check acme.json permissions:**
   ```bash
   ls -la traefik/acme/acme.json
   # Should be: -rw------- (600)
   chmod 600 traefik/acme/acme.json
   ```

4. **For development, use HTTP:**
   - Update Traefik config to allow HTTP
   - Test with `http://api.staging.stackyn.com/health`

### Issue: 502 Bad Gateway

**Symptoms**: Traefik returns 502 when accessing backend

**Solutions**:

1. **Backend not running:**
   ```bash
   docker compose restart backend
   ```

2. **Backend not on correct port:**
   ```bash
   # Check backend is listening on 8080
   docker compose exec backend netstat -tlnp | grep 8080
   ```

3. **Network issue:**
   ```bash
   # Check if backend is on same network
   docker network inspect stackyn-network
   ```

### Issue: Frontend not loading

**Symptoms**: `https://staging.stackyn.com` shows error

**Solutions**:

1. **Check frontend container:**
   ```bash
   docker compose ps frontend
   docker compose logs frontend
   ```

2. **Test frontend directly:**
   ```bash
   docker compose exec frontend curl http://localhost:3000
   ```

3. **Check if dist/ was built:**
   ```bash
   docker compose exec frontend ls -la /app/dist
   ```

## Quick Test Script

Create a test script `test-services.sh`:

```bash
#!/bin/bash

echo "=== Testing Stackyn Services ==="

echo ""
echo "1. Checking services..."
docker compose ps

echo ""
echo "2. Testing backend health (direct)..."
docker compose exec -T backend curl -s http://localhost:8080/health

echo ""
echo "3. Testing backend health (via Traefik)..."
curl -s https://api.staging.stackyn.com/health || echo "Failed - check DNS/SSL"

echo ""
echo "4. Testing frontend (direct)..."
docker compose exec -T frontend curl -s http://localhost:3000 | head -n 5

echo ""
echo "5. Testing database..."
docker compose exec -T postgres psql -U stackyn -d stackyn -c "SELECT 1;" 2>&1 | grep -q "1 row" && echo "Database OK" || echo "Database connection failed"

echo ""
echo "=== Test Complete ==="
```

Make it executable and run:
```bash
chmod +x test-services.sh
./test-services.sh
```

## Expected Responses

### Health Endpoint
```json
{"status":"ok"}
```

### API Endpoints
```json
// GET /api/v1/apps
[]

// GET /health
{"status":"ok"}
```

## Debugging Tips

1. **Always check logs first:**
   ```bash
   docker compose logs -f [service]
   ```

2. **Test services directly before testing through Traefik:**
   ```bash
   docker compose exec [service] curl http://localhost:[port]/[endpoint]
   ```

3. **Verify Traefik configuration:**
   ```bash
   docker compose exec traefik traefik version
   docker compose exec traefik cat /etc/traefik/traefik.yml
   ```

4. **Check network connectivity:**
   ```bash
   docker compose exec backend ping postgres
   docker compose exec frontend ping backend
   ```

5. **Monitor Traefik in real-time:**
   ```bash
   docker compose logs -f traefik
   # Then try accessing endpoints in another terminal
   ```

