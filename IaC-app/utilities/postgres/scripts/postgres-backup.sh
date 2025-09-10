#!/bin/bash
BACKUP_DIR="/wsf_data_dir/backups"
printf "Backup directory: $BACKUP_DIR\n"
printf "Creating the backup directory...\n"
mkdir -p /wsf_data_dir/backups
export PGPASSWORD=$POSTGRES_PASSWORD
TIMESTAMP=$(date +%Y%m%d_%H%M%S)
# Backup.
pg_dump -v -h "$PGHOST" -p 5432 -U "$POSTGRES_USER" "$POSTGRES_DB" > "$BACKUP_DIR/${TIMESTAMP}_$POSTGRES_DB.sql"
sleep 300s
printf "Backup created at $BACKUP_DIR"
# Optional: Compress the backup
# gzip "$BACKUP_DIR/$DB_NAME-$TIMESTAMP.sql"
