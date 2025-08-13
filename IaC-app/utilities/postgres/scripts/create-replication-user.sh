#!/bin/bash
# The '-e' option tells the shell to exit immediately if any invoked command exits with a non-zero
# status, which is the convention for indicating errors in Unix-like operating systems. It's a
# handy option for scripting, as it helps to catch errors and bugs early.
set -e

# The key difference with <<- (compared to just <<) is that it strips all leading tab characters
# from the input lines within the here document (heredoc). This is useful for indentation, allowing
# you to format your code nicely without affecting the actual input data.
# https://www.postgresql.org/docs/current/warm-standby.html#STREAMING-REPLICATION
# https://www.postgresql.org/docs/current/app-psql.html
psql -v ON_ERROR_STOP=1 --username "$REPLICATION_USER" --dbname "$POSTGRES_DB" <<-EOSQL
  CREATE ROLE '$REPLICATION_USER' WITH REPLICATION PASSWORD '$REPLICATION_PASSWORD' LOGIN;
EOSQL
