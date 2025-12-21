# Fixing Black Screen Issue

## Problem
Frontend shows for 100-200ms then goes black. This is typically caused by:
1. JavaScript error not being caught
2. API connection failure
3. SSL certificate issues
4. CORS errors

## Fixes Applied

1. **Added Error Boundary** - Catches React errors and displays them
2. **Improved Error Logging** - Better console logging for debugging
3. **API Timeout Handling** - Prevents hanging requests

## Next Steps

### 1. Rebuild Frontend

```bash
cd /opt/stackyn
docker compose build frontend
docker compose up -d frontend
```

### 2. Check Browser Console

Open browser DevTools (F12) and check:
- Console tab for errors
- Network tab for failed requests
- Look for:
  - CORS errors
  - SSL certificate errors
  - Failed API calls

### 3. Test API Directly

```bash
# Test if API is accessible
curl -k https://api.staging.stackyn.com/health
# Or HTTP
curl http://api.staging.stackyn.com/health
```

### 4. Check API URL Configuration

The frontend is configured to use: `https://api.staging.stackyn.com`

If SSL isn't ready, you may need to:
- Use HTTP temporarily: `http://api.staging.stackyn.com`
- Or wait for SSL certificates to be generated

### 5. Temporary Fix: Use HTTP for API

If SSL is causing issues, rebuild frontend with HTTP API URL:

Edit `docker-compose.yml`:
```yaml
frontend:
  build:
    args:
      - VITE_API_BASE_URL=http://api.staging.stackyn.com  # Change to HTTP
```

Then rebuild:
```bash
docker compose build frontend
docker compose up -d frontend
```

### 6. Check CORS

Ensure backend allows requests from frontend domain:
- Frontend: `https://staging.stackyn.com`
- Backend should allow CORS from this origin

### 7. Verify Error Boundary is Working

After rebuild, if there's an error, you should see:
- Error message displayed on screen
- "Reload Page" button
- Error details in expandable section

## Debugging

1. **Check frontend logs:**
   ```bash
   docker compose logs -f frontend
   ```

2. **Check browser console:**
   - Open DevTools (F12)
   - Look for red error messages
   - Check Network tab for failed requests

3. **Test API endpoint:**
   ```bash
   curl -v https://api.staging.stackyn.com/health
   ```

4. **Check if API URL is correct:**
   - Open browser console
   - Look for: "API Base URL: ..."
   - Verify it matches your backend URL

## Common Issues

### SSL Certificate Not Ready
- Use HTTP temporarily
- Wait for Let's Encrypt to generate certificates
- Check Traefik logs: `docker compose logs traefik | grep acme`

### CORS Error
- Backend needs to allow frontend origin
- Check backend CORS configuration
- Verify `Access-Control-Allow-Origin` header

### API Not Accessible
- Check if backend is running: `docker compose ps backend`
- Test backend directly: `docker compose exec backend curl http://localhost:8080/health`
- Check Traefik routing: `docker compose logs traefik | grep backend`

