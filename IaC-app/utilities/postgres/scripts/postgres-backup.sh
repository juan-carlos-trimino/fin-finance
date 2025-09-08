#!/bin/bash
BACKUP_DIR="/wsf_data_dir/postgres/backups"
printf "Backup directory: $BACKUP_DIR\n"
TIMESTAMP=$(date +%Y%m%d_%H%M%S)

mkdir -p /wsf_dir/postgres/backups

export PGPASSWORD=$POSTGRES_PASSWORD
# Backup.
pg_dump -h "$PGHOST" -U "$POSTGRES_USER" "$POSTGRES_DB" > "$BACKUP_DIR/$TIMESTAMP_$POSTGRES_DB.sql"
# sleep 120s
printf "Backup created at $BACKUP_DIR"

# Optional: Compress the backup
# gzip "$BACKUP_DIR/$DB_NAME-$TIMESTAMP.sql"
