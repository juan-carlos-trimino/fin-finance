#!/bin/bash
BACKUP_DIR="/wsf_data_dir/backups"
printf "Backup directory: $BACKUP_DIR\n"
TIMESTAMP=$(date +%Y%m%d_%H%M%S)


export PGPASSWORD="PG_PASSWORD"

printf "pgpaasword=$PGPASSWORD\n"
printf "pgdb=$POSTGRES_DB\n"
printf "pguser=$POSTGRES_USER\n"


# pg_dump -U "$DB_USER" "$DB_NAME" > "$BACKUP_DIR/$DB_NAME-$TIMESTAMP.sql"

# pg_dump -h remote_host -p 5432 -U db_user db_name

# pg_basebackup -v -h $PGHOST -p 5432 -D $PGDATA -U replication -R -Xs -Fp

# Optional: Compress the backup
# gzip "$BACKUP_DIR/$DB_NAME-$TIMESTAMP.sql"
