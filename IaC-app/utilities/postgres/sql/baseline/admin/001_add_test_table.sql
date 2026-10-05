-- Connect straight to the destination database namespace.
\c finances

-- LOCK: Protect Migration Application Sequence
SET lock_timeout = '5s';
SELECT pg_advisory_lock(9876543210);

-- Define a runtime checker wrapper using standard psql blocks.
SELECT (SELECT COUNT(*) FROM schema_migrations WHERE version = '002_add_audit_logs') > 0 AS migration_exists \gset

\if :migration_exists
  \qecho Migration '002_add_audit_logs' already applied. Exiting cleanly.
  \q
\else

  /**************************************************************************************************
                                     *** APPLY YOUR NEW UPGRADES HERE ***
  **************************************************************************************************/
  CREATE TABLE audit_logs (
      id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
      action TEXT NOT NULL,
      performed_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
  );

  -- Log the upgrade execution into your history matrix
  INSERT INTO schema_migrations (version) VALUES ('002_add_audit_logs');

  \qecho Migration '002_add_audit_logs' applied successfully.
\endif
-- Closing connection drops the lock automatically
