# User Deployed Apps

This directory contains applications deployed by users through the Stackyn platform.

## Structure

Each deployed app gets its own subdirectory:

```
apps/
├── app-001/              # Example app
│   ├── docker-compose.yml
│   ├── env/
│   │   └── .env
│   └── data/            # App-specific data
├── app-002/
└── README.md
```

## App Directory Contents

### docker-compose.yml

Each app has its own `docker-compose.yml` file that defines:
- Container configuration
- Environment variables
- Volumes
- Networks
- Ports

### env/

Environment variables specific to the app:
- Database credentials
- API keys
- Configuration values

### data/

Persistent data for the app:
- Database files
- Uploads
- Cache
- Logs

## App Naming Convention

Apps are named using the pattern: `app-{id}` where `{id}` is the app ID from the database.

## Management

Apps are managed through:
- Stackyn frontend UI
- Backend API
- Direct Docker Compose commands (advanced)

## Security

- Each app runs in its own directory
- Environment variables are isolated
- Volumes are scoped to the app
- Network isolation via Docker networks

## Backup

To backup an app:
```bash
# Backup app directory
tar -czf app-001-backup.tar.gz apps/app-001/

# Backup volumes
docker volume ls
docker run --rm -v app-001_data:/data -v $(pwd):/backup alpine tar czf /backup/app-001-data.tar.gz /data
```

## Cleanup

To remove an app:
```bash
# Stop and remove containers
cd apps/app-001
docker-compose down -v

# Remove directory
cd ../..
rm -rf apps/app-001
```

