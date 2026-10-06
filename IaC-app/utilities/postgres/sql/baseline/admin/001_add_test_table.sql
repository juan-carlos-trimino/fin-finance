
SELECT CONCAT('*** Output from script, run began at: ', NOW(), ' ***') AS msg \gset
\qecho :msg

-- LOCK: Protect Migration Application Sequence
SET lock_timeout = '5s';
-- Using the same global 64-bit ID. This forces concurrent pods to wait.
SELECT pg_advisory_lock(9876543210);

-- Evaluate if this specific migration script version has already been run.
SELECT (SELECT COUNT(*) FROM fin.schema_migrations WHERE version = :'MIGRATION_NAME') > 0 AS migration_exists \gset

\if :migration_exists
  \qecho Migration :'MIGRATION_NAME' already applied. Exiting...
  \q
\else
  /******************************************************************************************************************************
                                      *** APPLY THE NEW UPGRADE HERE (BEGIN) ***
  *******************************************************************************************************************************
  Keep the upgrade immutable: Once a migration script is deployed to production, never modify its contents. If you need to
  make a change (like altering a column), create a brand new upgrade  using this same template.
  ******************************************************************************************************************************/

  CREATE TABLE IF NOT EXISTS fin.schema_migrations(
    migration   TEXT PRIMARY KEY,
    applied_at  TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
  );

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
  /******************************************************************************************************************************
                                       *** APPLY THE NEW UPGRADE HERE (END) ***
  ******************************************************************************************************************************/

  -- Log the upgrade execution into the history table matrix so it never runs again.
  INSERT INTO fin.schema_migrations(migration)
    VALUES(:'MIGRATION_NAME')
    ON CONFLICT(migration) DO NOTHING;

  \qecho Migration :'MIGRATION_NAME' applied successfully.
\endif
-- When psql finishes executing this file, the connection drops, causing Postgres to automatically release the advisory lock.
