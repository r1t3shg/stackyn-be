# Frontend Deployment Guide for VPS

This guide explains how to deploy the Next.js frontend application to a VPS.

## Prerequisites

- A VPS with Ubuntu/Debian (or similar Linux distribution)
- Node.js 20+ installed
- Nginx (for reverse proxy)
- Domain name pointing to your VPS (optional but recommended)
- SSL certificate (Let's Encrypt recommended)

## Option 1: Direct Node.js Deployment

### Step 1: Prepare Your VPS

```bash
# Update system packages
sudo apt update && sudo apt upgrade -y

# Install Node.js 20
curl -fsSL https://deb.nodesource.com/setup_20.x | sudo -E bash -
sudo apt install -y nodejs

# Install PM2 for process management
sudo npm install -g pm2

# Install Nginx
sudo apt install -y nginx
```

### Step 2: Clone and Build the Application

```bash
# Navigate to your deployment directory
cd /var/www

# Clone your repository (or upload files)
git clone <your-repo-url> stackyn-frontend
cd stackyn-frontend/frontend

# Install dependencies
npm ci

# Set environment variable
export NEXT_PUBLIC_API_BASE_URL=https://staging.stackyn.com

# Build the application
npm run build
```

### Step 3: Create Environment File

Create a `.env.production` file in the frontend directory:

```bash
echo "NEXT_PUBLIC_API_BASE_URL=https://staging.stackyn.com" > .env.production
```

### Step 4: Start with PM2

**Option A: Using the deployment script (recommended)**

```bash
# Make the script executable
chmod +x deploy.sh

# Run the deployment script
./deploy.sh

# Setup PM2 to start on boot (first time only)
pm2 startup
```

**Option B: Manual deployment**

```bash
# Start the application with PM2
pm2 start npm --name "stackyn-frontend" -- start

# Save PM2 configuration
pm2 save

# Setup PM2 to start on boot
pm2 startup
```

### Step 5: Configure Nginx Reverse Proxy

Create an Nginx configuration file:

```bash
sudo nano /etc/nginx/sites-available/stackyn-frontend
```

Add the following configuration (replace `your-domain.com` with your domain):

```nginx
server {
    listen 80;
    server_name your-domain.com;

    # Redirect HTTP to HTTPS (if using SSL)
    # return 301 https://$server_name$request_uri;

    location / {
        proxy_pass http://localhost:3000;
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection 'upgrade';
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
        proxy_cache_bypass $http_upgrade;
    }
}

# If using SSL (recommended)
# server {
#     listen 443 ssl http2;
#     server_name your-domain.com;
#
#     ssl_certificate /etc/letsencrypt/live/your-domain.com/fullchain.pem;
#     ssl_certificate_key /etc/letsencrypt/live/your-domain.com/privkey.pem;
#
#     location / {
#         proxy_pass http://localhost:3000;
#         proxy_http_version 1.1;
#         proxy_set_header Upgrade $http_upgrade;
#         proxy_set_header Connection 'upgrade';
#         proxy_set_header Host $host;
#         proxy_set_header X-Real-IP $remote_addr;
#         proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
#         proxy_set_header X-Forwarded-Proto $scheme;
#         proxy_cache_bypass $http_upgrade;
#     }
# }
```

Enable the site:

```bash
sudo ln -s /etc/nginx/sites-available/stackyn-frontend /etc/nginx/sites-enabled/
sudo nginx -t
sudo systemctl reload nginx
```

### Step 6: Setup SSL (Optional but Recommended)

```bash
# Install Certbot
sudo apt install -y certbot python3-certbot-nginx

# Obtain SSL certificate
sudo certbot --nginx -d your-domain.com

# Certbot will automatically configure Nginx and set up auto-renewal
```

## Option 2: Docker Deployment

### Step 1: Install Docker

```bash
# Install Docker
curl -fsSL https://get.docker.com -o get-docker.sh
sudo sh get-docker.sh

# Install Docker Compose
sudo apt install -y docker-compose
```

### Step 2: Build and Run Docker Container

```bash
cd /var/www/stackyn-frontend/frontend

# Build the Docker image
docker build -t stackyn-frontend .

# Run the container
docker run -d \
  --name stackyn-frontend \
  -p 3000:3000 \
  -e NEXT_PUBLIC_API_BASE_URL=https://staging.stackyn.com \
  --restart unless-stopped \
  stackyn-frontend
```

### Step 3: Configure Nginx

Follow Step 5 from Option 1 to configure Nginx as a reverse proxy.

## Environment Variables

The application uses the following environment variable:

- `NEXT_PUBLIC_API_BASE_URL`: The base URL for the API (default: `http://localhost:8080`)
  - For production: `https://staging.stackyn.com`

**Important**: In Next.js, environment variables prefixed with `NEXT_PUBLIC_` are embedded into the JavaScript bundle at build time. Make sure to set this variable before running `npm run build`.

## Updating the Application

### For Direct Node.js Deployment:

**Using the deployment script:**

```bash
cd /var/www/stackyn-frontend/frontend
git pull
./deploy.sh
```

**Manual update:**

```bash
cd /var/www/stackyn-frontend/frontend
git pull
npm ci
export NEXT_PUBLIC_API_BASE_URL=https://staging.stackyn.com
npm run build
pm2 restart stackyn-frontend
```

### For Docker Deployment:

```bash
cd /var/www/stackyn-frontend/frontend
git pull
docker build -t stackyn-frontend .
docker stop stackyn-frontend
docker rm stackyn-frontend
docker run -d \
  --name stackyn-frontend \
  -p 3000:3000 \
  -e NEXT_PUBLIC_API_BASE_URL=https://staging.stackyn.com \
  --restart unless-stopped \
  stackyn-frontend
```

## Troubleshooting

### Check Application Logs

```bash
# PM2 logs
pm2 logs stackyn-frontend

# Docker logs
docker logs stackyn-frontend
```

### Check Nginx Status

```bash
sudo systemctl status nginx
sudo nginx -t
```

### Verify Environment Variable

```bash
# Check if the environment variable is set correctly
echo $NEXT_PUBLIC_API_BASE_URL
```

### Test API Connection

```bash
curl https://staging.stackyn.com/health
```

## Security Considerations

1. **Firewall**: Configure UFW or iptables to only allow necessary ports (80, 443, 22)
2. **SSL**: Always use HTTPS in production
3. **Environment Variables**: Never commit `.env.production` to version control
4. **Updates**: Keep Node.js, Nginx, and system packages updated

## Monitoring

Consider setting up monitoring with:
- PM2 monitoring (built-in)
- Nginx access/error logs
- Application health checks
- Uptime monitoring services

