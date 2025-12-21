#!/bin/bash

# Deployment script for Stackyn Frontend
# Usage: ./deploy.sh

set -e

echo "🚀 Starting deployment..."

# Check if Node.js is installed
if ! command -v node &> /dev/null; then
    echo "❌ Node.js is not installed. Please install Node.js 20+ first."
    exit 1
fi

# Check Node.js version
NODE_VERSION=$(node -v | cut -d'v' -f2 | cut -d'.' -f1)
if [ "$NODE_VERSION" -lt 20 ]; then
    echo "❌ Node.js version 20+ is required. Current version: $(node -v)"
    exit 1
fi

# Set environment variable
export NEXT_PUBLIC_API_BASE_URL=https://staging.stackyn.com

# Install dependencies
echo "📦 Installing dependencies..."
npm ci

# Build the application
echo "🔨 Building application..."
npm run build

# Check if PM2 is installed
if command -v pm2 &> /dev/null; then
    echo "🔄 Restarting application with PM2..."
    pm2 restart stackyn-frontend || pm2 start npm --name "stackyn-frontend" -- start
    pm2 save
    echo "✅ Application restarted successfully!"
else
    echo "⚠️  PM2 is not installed. Install it with: npm install -g pm2"
    echo "📝 To start the application manually, run: npm start"
fi

echo "✅ Deployment completed!"

