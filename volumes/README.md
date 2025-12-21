# Persistent Volumes

This directory contains persistent data for Stackyn services.

## Structure

```
volumes/
├── db/          # PostgreSQL database files
├── redis/       # Redis data (if used)
└── logs/        # Application logs
```

## Volumes

### db/

PostgreSQL database files:
- Database data files
- Transaction logs
- Configuration

**Backup:**
```bash
# Backup database
docker compose exec postgres pg_dump -U stackyn stackyn > backup.sql

# Or backup volume directly
docker run --rm -v stackyn_db-data:/data -v $(pwd):/backup alpine tar czf /backup/db-backup.tar.gz /data
```

**Restore:**
```bash
# Restore from SQL dump
docker compose exec -T postgres psql -U stackyn stackyn < backup.sql
```

### redis/

Redis data files (if Redis is used):
- RDB snapshots
- AOF files

### logs/

Application logs:
- Backend logs
- Worker logs
- Deployment logs

**Log Rotation:**
Consider setting up log rotation to prevent disk space issues:
```bash
# Example logrotate config
/opt/stackyn/volumes/logs/*.log {
    daily
    rotate 7
    compress
    missingok
    notifempty
}
```

## Permissions

Ensure proper permissions:
```bash
sudo chown -R 999:999 volumes/db      # PostgreSQL user
sudo chown -R 999:999 volumes/redis    # Redis user (if used)
sudo chown -R $USER:$USER volumes/logs  # Your user
```

## Backup Strategy

1. **Database**: Daily backups of `volumes/db/`
2. **Logs**: Weekly rotation, monthly archival
3. **Redis**: If used, backup before major updates

## Disk Space

Monitor disk usage:
```bash
du -sh volumes/*
df -h
```

## Security

- Do not commit sensitive data
- Use proper file permissions
- Encrypt backups if containing sensitive data
- Regularly audit access logs

