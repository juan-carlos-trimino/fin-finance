# PostgreSQL Distributed Migration Engine

A lightweight, concurrent-safe database migration and bootstrapping engine built in **Go** using native **PostgreSQL Advisory Locks** and the `psql` CLI utility.

## 🚀 Overview

This architecture provides a decentralized migration runner capable of running across multiple parallel application nodes (e.g., Kubernetes pods) simultaneously without race conditions, duplicate database creations, or schema corruption.

### Key Features
* **Distributed Locking:** Uses Postgres 64-bit advisory locks to serialize initialization across auto-scaling application instances.
* **Two-Stage Isolation:** Protects cluster-level definitions (roles, database creation) before executing internal schema modifications (tables, indexes, extensions).
* **Automated `search_path` Scoping:** Isolates execution automatically into application-specific namespaces.
* **Fail-Fast Loop:** Halts structural updates immediately upon encountering an error, keeping log aggregators clean and predictable.

---

## 📁 Directory Structure & Naming Conventions

All migration files must reside in your designated migrations directory and utilize a strict **three-digit chronological lexical sorting prefix** (e.g., `000`, `001`, `002`).

```text
migrations/
├── 000_admin.sql          # Cluster Bootstrapper (Creates DB, Users, Roles)
├── 001_add_vault.sql      # Schema Migration (Creates tables, indexes, tracking)
└── 002_audit_logs.sql     # Subsequent Migration Step
```

### Script Ordering Rules
1. **`000_admin.sql`:** Must always be reserved exclusively for setting up the cluster roles, permissions, and running `CREATE DATABASE`.
2. **`001_xxx.sql`+:** Reserved for creating structural tables, columns, indexes, and user permissions.

---

## 🛠️ Writing Script Files

### 1. The Bootstrap Template (`000_admin.sql`)
This file runs against the system's default administrative database (e.g., `postgres`) to guarantee structural setup operations execute before shifting dependencies to the localized app environment.

* **Important Requirement:** Must utilize `\c <target_db>` at the bottom of the bootstrap phase to establish a safe hand-off sequence for Lock 2 initialization.

```sql
SELECT CONCAT('*** Run began at: ', NOW(), ' ***') AS msg \gset
\qecho :msg

-- LOCK 1: Protect Cluster Actions (Roles & Database Creation)
SET lock_timeout = '5s';
SELECT pg_advisory_lock(9876543210);

-- [... Conditional check and Database/Role creation logic ...]

-- Connect to the target database
\c :DB_NAME

-- LOCK 2: Protect Schema Actions inside the newly minted database
SET lock_timeout = '5s';
SELECT pg_advisory_lock(9876543210);

-- Initialize foundational tables, extensions, and schema history here
```

### 2. Subsequent Migration Template (`001_xxx.sql` onwards)
Subsequent migrations are executed directly against the target application database natively by the Go runner.

* **CRUCIAL RULE:** Do **NOT** use `\c <database_name>` or `\connect` at the top of these subsequent scripts. Doing so will drop the Go runner's secure connection context, instantly wiping out your advisory lock safety net.

```sql
SELECT CONCAT('*** Run began at: ', NOW(), ' ***') AS msg \gset
\qecho :msg

-- LOCK: Protect Migration Sequence
SET lock_timeout = '5s';
SELECT pg_advisory_lock(9876543210);

-- NAMESPACE SCOPING: Direct variables to target app schema space
CREATE SCHEMA IF NOT EXISTS fin;
SET search_path = fin, public;

-- Evaluate if migration was already applied using Go context parameters
SELECT (NOT EXISTS (SELECT 1 FROM schema_migrations WHERE version = :'MIGRATION_NAME')) AS run_script \gset

\if :run_script
  -- =========================================================================
  -- APPLY MIGRATIONS HERE
  -- =========================================================================
  CREATE TABLE user_data_vault ( ... );

  -- Record migration history
  INSERT INTO schema_migrations (version) VALUES (:'MIGRATION_NAME');
  \qecho Migration :'MIGRATION_NAME' applied successfully.
\else
  \qecho Migration :'MIGRATION_NAME' already applied. Exiting cleanly.
  \q
\endif
```

---

## 💻 Go Implementation Reference

The main orchestration entrypoint is driven by the `ExecuteSqlScript` function. It executes the directory collection scanner, performs rigid validation filtering to eliminate string-injection exploits, and shifts target credentials seamlessly after indexing past the initial bootstrapping file.

### Safe Execution Logic Loop
1. Scans directory and isolates files ending with `.sql`.
2. Validates if `000_admin.sql` is present. If missing, it executes a passive query using `-tA` tuples-only flags to verify if the server environment was already initialized. If missing and no server instance is found, it safely invokes a hard fail-fast panic.
3. Alphabetically sorts the file slide arrays to guarantee uniform cross-platform behavior.
4. Spawns automated sub-processes using `exec.Command("psql", ...)` using safe parameters (`-v`).

---

## ⚠️ Important Rules for Production

1. **Immutability:** Once a migration script has been merged or deployed to an environment, **never modify its contents**. If structural deviations are required (e.g., adding an extra index or altering an active column constraint), generate an isolated sequential file step (e.g., `002_alter_vault.sql`).
2. **Variable Substitution Rules:**
   * Use **`:'VAR_NAME'`** (with single quotes) for string values inside `WHERE` filters, log printing, or string insertions.
   * Use **`:VAR_NAME`** (without single quotes) for boolean flags or raw SQL identifiers like database and table names.
