# Stackyn Frontend

A Next.js frontend application for managing applications and deployments on the Stackyn PaaS platform.

## Features

- 📱 **App Management**: Create, view, and delete applications
- 🚀 **Deployment Tracking**: Monitor deployment status and view logs
- 🎨 **Modern UI**: Beautiful, responsive interface built with Tailwind CSS
- ⚡ **Real-time Updates**: View deployment status and logs in real-time

## Getting Started

### Prerequisites

- Node.js 18+ and npm
- Backend API server running (see backend README)

### Installation

1. Install dependencies:
```bash
npm install
```

2. Configure the API base URL:
   - Copy `.env.example` to `.env.local`
   - Update `NEXT_PUBLIC_API_BASE_URL` to match your backend server URL

3. Run the development server:
```bash
npm run dev
```

4. Open [http://localhost:3000](http://localhost:3000) in your browser

## Project Structure

```
frontend/
├── app/                    # Next.js app directory
│   ├── apps/              # App-related pages
│   │   ├── new/           # Create new app
│   │   └── [id]/          # App details page
│   └── page.tsx           # Home page (apps list)
├── components/            # React components
│   ├── AppCard.tsx        # App card component
│   ├── DeploymentCard.tsx # Deployment card component
│   ├── StatusBadge.tsx    # Status badge component
│   └── LogsViewer.tsx     # Logs viewer component
├── lib/                   # Utility functions
│   ├── api.ts            # API client functions
│   ├── config.ts         # Configuration
│   └── types.ts          # TypeScript type definitions
└── public/               # Static assets
```

## API Integration

The frontend communicates with the backend API at the following endpoints:

- `GET /api/v1/apps` - List all apps
- `POST /api/v1/apps` - Create a new app
- `GET /api/v1/apps/{id}` - Get app by ID
- `DELETE /api/v1/apps/{id}` - Delete an app
- `POST /api/v1/apps/{id}/redeploy` - Redeploy an app
- `GET /api/v1/apps/{id}/deployments` - List deployments for an app
- `GET /api/v1/deployments/{id}` - Get deployment by ID
- `GET /api/v1/deployments/{id}/logs` - Get deployment logs

## Building for Production

```bash
npm run build
npm start
```

## Environment Variables

- `NEXT_PUBLIC_API_BASE_URL`: Base URL for the backend API (default: `http://localhost:8080`)
