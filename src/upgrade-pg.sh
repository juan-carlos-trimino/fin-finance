#!/bin/bash
# Stop if any error happens
set -e
#############################################################
# Change the permissions of the script to make it executable:
# chmod +x update-pg.sh
# Execute the script to update PostgreSQL:
# sudo ./upgrade-pg.sh postgresql-18.1.2
#############################################################
echo "1. Backing up all databases just in case..."
# postgres -c "pg_dumpall > ./all_dbs_backup.sql"
# 1. Stop the current PostgreSQL service
echo "Stopping PostgreSQL..."
systemctl stop postgresql
# 2. Update Ubuntu package lists
echo "Updating package lists..."
apt-get update
# apt-get install -y apt-transport-https ca-certificates
# 3. Upgrade the specific PostgreSQL packages.
echo "Upgrading PostgreSQL to postgresql-$1..."
apt-get install --only-upgrade -y postgresql-"$1" postgresql-client-"$1" postgresql-contrib-"$1"
# 4. Verify the update
echo "Starting PostgreSQL..."
systemctl start postgresql
# Check PostgreSQL status
service postgresql status
# 5. Check the new version
echo "Current PostgreSQL version:"
psql --version
echo "6. Upgrade completed successfully!"

# make sure there are no errors:tail -n 50 /var/log/postgresql/postgresql-18-main.log
