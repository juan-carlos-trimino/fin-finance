
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

  /***
  It allows PostgreSQL to compare and sort text while completely ignoring both uppercase/lowercase differences and accent marks
  (like e, é, è, and Ê will all be treated as the exact same letter).

  -------------------
  Parameter Breakdown
  -------------------
  * PROVIDER = icu
  This tells Postgres to use the ICU (International Components for Unicode) library instead of the standard operating system
  (libc). ICU is required here because standard OS collations usually cannot handle complex rules like advanced case-and-accent
  insensitivity.
  * LOCALE = 'und-u-ks-level1'
  This is a special Unicode locale string that defines the sorting strength:
    * und: Short for "undetermined" (or root). It means it uses generic, universal Unicode sorting rules rather than rules tied
      to a specific country or language (like French or German).
    * ks-level1: Sets the Collation Strength to Level 1 (Primary). At Level 1, ICU only looks at the base alphabet letters. It
      ignores Level 2 (accents) and Level 3 (case).
  * DETERMINISTIC = false
  By default, Postgres collations are deterministic, meaning strings are only equal if their exact byte-by-byte data is identical.
  Setting this to false allows Postgres to say that 'Cafe' and 'café' are deterministically equal, even though their underlying
  bytes are different. This is what enables you to use standard operators like = or LIKE for case/accent-insensitive matching.
  ***/
  -- 1. Drop the dependent index first.
  DROP INDEX IF EXISTS fin.idx_users_slider_pagination;
  -- 2. Drop the old collation cleanly from the schema.
  DROP COLLATION IF EXISTS fin.euro_case_insensitive;
  -- 3. Create the non-deterministic ICU collation cleanly.
  CREATE COLLATION fin.euro_case_insensitive(
    PROVIDER = icu,
    LOCALE = 'und-u-ks-level1',  -- Ignore both case and accents.
    DETERMINISTIC = false
  );
  -- 4. Rebuild the index matching the Go query exactly.
  --    This includes the raw 'user_name' to handle the case/accent ties natively.
  CREATE INDEX fin.idx_users_slider_pagination
    ON fin.customers_credentials
    USING btree(user_name COLLATE fin.euro_case_insensitive ASC,
                user_name ASC,
                id ASC);

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
