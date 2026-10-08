
SELECT CONCAT('*** Output from script, run began at: ', NOW(), ' ***') AS msg \gset
\qecho :msg

-- LOCK: Protect Migration Application Sequence
SET lock_timeout = '5s';
-- Using the same global 64-bit ID. This forces concurrent pods to wait.
SELECT pg_advisory_lock(9876543210);

-- SCHEMA SCOPING: Route all default queries directly into the 'fin' schema.
-- Ensure the custom schema physically exists first.
CREATE SCHEMA IF NOT EXISTS fin;
-- Point this active database connection session to search 'fin' before 'public'.
SET search_path = fin, public;

SELECT EXISTS(
  SELECT 1 FROM fin.migration_history WHERE migration_name = :'MIGRATION_NAME'
) AS migration_exists \gset

\if :migration_exists
  \qecho Migration :'MIGRATION_NAME' already applied, skipping...
  \q
\else
  \qecho Applying migration :'MIGRATION_NAME'...
  /******************************************************************************************************************************
                                      *** APPLY THE NEW UPGRADE HERE (BEGIN) ***
  *******************************************************************************************************************************
  Keep the upgrade immutable: Once a migration script is deployed to production, never modify its contents. If you need to
  make a change (like altering a column), create a brand new upgrade  using this same template.
  ******************************************************************************************************************************/

  CREATE TABLE IF NOT EXISTS fin.user_data_vault(
    user_name       TEXT NOT NULL,
    partition_name  TEXT NOT NULL,
    partition_data  JSONB NOT NULL,
    created_at      TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at      TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    --PK ensures the uniqueness of user & partition.
    PRIMARY KEY(user_name, partition_name) -- Composite Primary Key
  );

/***
Indexing JSON data
The more indexes Postgres has, the more time it takes to maintain and update them, which can affect the latency of certain queries and overall application performance.

The B-tree data structure is used by Postgres by default to store and access indexed data internally.

Using GIN Indexes
By using GIN as the underlying data structure, an application can create a regular single-column index on JSON data to enable efficient searches across the nested JSON structure. Postgres supports two types of GIN indexes for JSON objects, each storing indexed data differently and supporting distinct operator classes for data access. The first type is the default GIN index, which extracts all the keys, values, and array elements from the original JSON object and adds them as distinct items to the index structure. The second type indexes only paths from the root of the JSON document down to each value and array element.

Using the default GIN index

***/

  -- Index to optimize querying inside the JSONB payload
  CREATE INDEX IF NOT EXISTS idx_user_data_vault
    ON fin.user_data_vault
    USING gin(partition_data);
  ANALYZE fin.user_data_vault;


ALTER DEFAULT PRIVILEGES IN SCHEMA fin
GRANT SELECT, INSERT, UPDATE, DELETE ON TABLES TO admin_role;
ALTER DEFAULT PRIVILEGES IN SCHEMA fin
GRANT USAGE, SELECT ON SEQUENCES TO admin_role;



  /******************************************************************************************************************************
                                       *** APPLY THE NEW UPGRADE HERE (END) ***
  ******************************************************************************************************************************/

  -- Using the conflict safety clause so these commands can be safely re-run without throwing an error.
  -- Log the upgrade execution into the history table matrix so it never runs again.
  INSERT INTO fin.migration_history(migration_name)
    VALUES(:'MIGRATION_NAME')
    ON CONFLICT(migration_name) DO NOTHING;

  \qecho Migration :'MIGRATION_NAME' applied successfully.
\endif
-- When psql finishes executing this file, the connection drops, causing Postgres to automatically release the advisory lock.
