
-- FORCE psql to immediately exit and close the session if ANY error occurs.
-- This ensures the session lock is never left hanging on a dead connection!
\set ON_ERROR_STOP on

SELECT CONCAT('*** Output from script, run began at: ', NOW(), ' ***') AS msg \gset
\qecho :msg

SELECT CONCAT('*** PostgreSQL version: ', (SELECT version()), ' ***') AS msg \gset
\qecho :msg

-- LOCK: Protect Migration Application Sequence
SET lock_timeout = '5s';
SELECT 'Confirmation: lock_timeout is set to ' || current_setting('lock_timeout') || '.';
-- Using the same global 64-bit ID. This forces concurrent pods to wait.
SELECT pg_advisory_lock(:'LOCK_ID'::bigint);

\qecho Lock (:LOCK_ID) successfully acquired.

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
