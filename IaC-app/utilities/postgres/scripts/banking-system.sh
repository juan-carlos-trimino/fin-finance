#!/bin/bash
#

print_time_elapsed() {
  printf "\n***********************************************************\n"
  printf "%s\n" "$1"
  printf "%s\n" "$(date -u)"
  printf "Time Elapsed: %02d hours %02d minutes and %02d seconds." "$(($2 / 3600))" \
         "$(($2 % 3600 / 60))" "$(($2 % 60))"
  printf "\n***********************************************************\n\n"
  return
}

printf "***************************************\n"
printf "Waiting for Postgres...\n"
printf "%s\n" "$(date -u)"
printf "***************************************\n"
start_time_seconds=$(date +%s)
# This script will repeatedly check the PostgreSQL server's status every 2 seconds until
# pg_isready returns an exit status of 0, indicating that the server is accepting
# connections.
until pg_isready --quiet -U $POSTGRES_USER;
do
  printf ".*."
  sleep 2s
done
duration=$(( $(date +%s) - start_time_seconds ))
print_time_elapsed "PostgreSQL is ready to accept connections!" "$duration"
printf "***************************************\n"
printf "Waiting for the sql script...\n"
printf "%s\n" "$(date -u)"
printf "***************************************\n"
export PGPASSWORD=$POSTGRES_PASSWORD
start_time_seconds=$(date +%s)
psql -v -U $POSTGRES_USER -d template1 -f "$1"
duration=$(( $(date +%s) - start_time_seconds ))
print_time_elapsed "Done running the sql scripts..." "$duration"
