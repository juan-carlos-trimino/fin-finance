/**************************************************************************************************
Notes
 * By default, Postgres ignores case and always turns every identifier into lowercase. Postgres
   preserves case only when the identifier is double quoted (e.g., "id"), and double-quoted
   identifiers, known as 'delimited identifiers', are case sensitive.
 * In Postgres, you can CAST between types with the :: shorthand.
 * In PostgreSQL, stored procedures are not atomic by default but can manage their own transactions
   explicitly to ensure atomicity. This is a key difference from PostgreSQL functions, which are
   always atomic and run within a single, implicit transaction.
 * Postgres executes a function atomically and transactionally; i.e, if the function fails at any
   step during its execution, all previous changes made within that function are rolled back,
   ensuring data integrity and consistency.
 * After a COMMIT or ROLLBACK is issued inside a procedure, a new transaction is automatically
   started; you do not need a separate START TRANSACTION command.
 * END is a Postgres extension to the SQL language that is equivalent to COMMIT.
Notes (Backup)
 To perform a dump.
 $ pg_dump finances > /tmp/finances.dump
 To restore (first create the database).
 $ psql finances < /tmp/finances.dump
 To time how long it takes to restore.
 $ time psql finances < /tmp/finances.dump
Notes (Indexes)
 * Postgres defaults to a B-tree index unless another data strcuture is specified. B-tree supports
   a wide range of data types and can be used for both equality searches (=) and range queries
   using greater than (>, >=), less than (<, <=), and BETWEEN operators.
 * It's good practice to run the ANALYZE command after creating an index.
   When a new index is created, the query planner may not take advantage of it immediately and may
   continue relying on previously collected statistics. Running ANALYZE updates the table
   statistics, allowing the planner to take full advantage of the new index.
 * Postgres automatically creates an index for the primary key column.
 * Indexes can significantly optimize query performance, but they don't come for free. Each time we
   create a new index, Postgres must maintain it by updating its structure whenever the value of an
   indexed column changes in the primary table.
   Apart from the index maintenance aspects, the more indexes Postgres has, the more time it will
   spend in the planning phase while selecting and generating the most efficient execution plan for
   a query.
**************************************************************************************************/
-- Online Banking System.
SELECT 'Output from script, run began at: ' AS "Script Information",
  NOW() AS "Date and Time Executed";

/**************************************************************************************************
                                 *** DATABASE ROLES AND PRIVILEGES ***
**************************************************************************************************/
DROP ROLE IF EXISTS admin_role;
/***
A Postgres cluster refers to the individual server/instance that's running and hosting (a cluster
of) databases. It does not mean that multiple servers are setup in a multi-node environment.

A ROLE exists at the cluster level. By convention, a ROLE that allows login is considered a user,
and a role that is not allowed to login is a group.

In highly concurrent environments, checking for a role's existence right before creating it can
cause a race condition. You can instead try to create the user unconditionally and safely intercept
the duplicate_object error.
***/
DO $$  -- The anonymous block executes procedural logic directly on the server.
BEGIN
  -- Since CREATE ROLE defaults to NOLOGIN, this role can't connect to the database, which is
  -- perfect for a group role.
  CREATE ROLE admin_role WITH NOLOGIN NOINHERIT NOSUPERUSER CREATEROLE NOCREATEDB
                              CONNECTION LIMIT 5;
EXCEPTION
  WHEN duplicate_object THEN
    RAISE NOTICE 'Role "admin_role" already exists, skipping...';
END
$$;

DROP ROLE IF EXISTS admin_user;
-- Create individual users.
DO $$  -- The anonymous block executes procedural logic directly on the server.
BEGIN
  CREATE ROLE admin_user WITH LOGIN INHERIT PASSWORD '12345' NOSUPERUSER NOCREATEROLE NOCREATEDB
                              CONNECTION LIMIT 5 IN ROLE admin_role;
EXCEPTION
  WHEN duplicate_object THEN
    RAISE NOTICE 'Role "admin_user" already exists, skipping...';
END
$$;

/**************************************************************************************************
                                            *** DATABASE ***
**************************************************************************************************/
CREATE DATABASE finances
WITH
  OWNER = admin_role
  ALLOW_CONNECTIONS = TRUE
  CONNECTION_LIMIT = -1  -- Unlimited connections.
  ENCODING = 'UTF8'
  LC_COLLATE = 'C.UTF8'  -- Determine the sort order of strings.
  LC_CTYPE = 'C.UTF8'  -- Define character classification rules.
  IS_TEMPLATE = FALSE  -- Only superusers or the database owner can clone the database.
  TEMPLATE = 'template0';
/**************************************************************************************************
                *** DATABASE ROLES AND PRIVILEGES (Database-level privileges) ***
**************************************************************************************************/
-- GRANT CONNECT ON DATABASE finances TO admin_role;
-- GRANT ALL PRIVILEGES ON DATABASE finances TO admin_role;
/***
Every Postgres cluster has an implicit role called 'public' which cannot be deleted. All other
roles are always granted membership in 'public' by default and inherit whatever privileges are
currently assigned to it. Unless otherwise modified, the privileges granted to the 'public' role
are as follows:
PostgreSQL 14 and below                   PostgreSQL 15 and above
-----------------------------------------------------------------
CONNECT                                   CONNECT
CREATE                                    TEMPORARY
TEMPORARY                                 EXECUTE (functions and procedures)
EXECUTE (functions and procedures)        USAGE (domains, languages, and types)
USAGE (domains, languages, and types)

Notice that the 'public' role always has the CONNECT privilege granted by default, which allows all
roles to connect to a newly created database. Without the privilege to connect to a database, none
of the newly created roles would be able to do much.
***/
-- Revoke all privileges that are granted by default to the 'public' role.
REVOKE ALL ON DATABASE finances FROM public;

/***
Connect to the database.
***/
\c finances

\qecho 'Current database version:'
SELECT version();

/***
Set the session to the new role. The role that is in force at the time of an object creation will
own the object. Essentially, the owner of an object is analogous to a superuser of that object.
***/
SET ROLE admin_role;

/**************************************************************************************************
                                         *** EXTENSION ***
***************************************************************************************************
To use bcrypt in Postgres, you can utilize the pgcrypto extension. This extension provides the
crypt() and gen_salt() functions necessary for secure password hashing and verification within SQL.
Enable the extension in your specific database by running the following SQL command.
***/
CREATE EXTENSION IF NOT EXISTS pgcrypto;

/**************************************************************************************************
                                           *** SCHEMAS ***
**************************************************************************************************/
CREATE SCHEMA IF NOT EXISTS fin;
/**************************************************************************************************
                *** DATABASE ROLES AND PRIVILEGES (Schema-level privileges) ***
**************************************************************************************************/
GRANT USAGE ON SCHEMA fin TO admin_role;
/***
Notice that you removed ALL privileges from 'public' at the database level, but not at the schema
level of the database. If you remove all privileges from the schema as well, then normal Postgres
commands would not work for many users without resetting additional privileges.
Note that starting with Postgres 15 the 'public' role can NO longer create anything by default,
regardless of the schema.
***/
REVOKE CREATE ON SCHEMA public FROM public;

/***
Set the search path to look for objects in your schema.
The role must have been granted USAGE privileges on your schema, otherwise the path will be
 ignored for your schema.
Setting the search path at the role level will override the database-wide defaults set by
 ALTER DATABASE.
If you need all users in a group to share the same search_path, use the following:
 Set at the Database Level: Apply the search_path to the database itself, so every new connection
                            defaults to it.
***/
--ALTER ROLE admin_role SET search_path = fin;
ALTER DATABASE finances SET search_path TO fin, public;

/**************************************************************************************************
                                           *** TABLES ***
***************************************************************************************************
customers-to-customer_contact_details Relationship: One-to-One
customers-to-credentials Relationship: One-to-One
URL: /register
***/
CREATE TABLE IF NOT EXISTS fin.customers(
  id                 INT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
  /***
  A column defined with a DEFAULT value and a CHECK constraint still needs an explicit NOT NULL
  constraint if you want to prevent NULL values.
  Here's why:
  DEFAULT Constraint: This constraint only applies if no value is provided during an INSERT
  operation. If you explicitly INSERT a NULL value into a column, the default value is bypassed,
  and the NULL is inserted (unless NOT NULL is present).
  CHECK Constraint: In SQL, any comparison involving a NULL value evaluates to UNKNOWN, not TRUE
  or FALSE. A CHECK constraint only fails if the condition evaluates to FALSE. If it evaluates to
  TRUE or UNKNOWN (due to a NULL value), the constraint is satisfied and the NULL is accepted.
  NOT NULL Constraint: This is the only constraint specifically designed to enforce the presence of
  data and prohibit NULL values in a column.
  ***/
  first_name         TEXT NOT NULL
  -- Block empty strings ('') and strings with only blanks (' ').
                       CONSTRAINT check_first_name
                         CHECK(TRIM(first_name) <> ''),
  middle_name        TEXT,
  last_name          TEXT NOT NULL
                       CONSTRAINT check_last_name
                         CHECK(TRIM(last_name) <> ''),
  marketing_consent  BOOLEAN NOT NULL DEFAULT FALSE,
  /***
  https://www.postgresql.org/docs/current/datatype-datetime.html
  This stores date and time along with time zone information. PostgreSQL automatically converts the
  timestamp to UTC for storage and adjusts it back based on the current time zone settings when
  queried. 8 bytes in length.
  ***/
  created_at         TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  updated_at         TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_last_name
  ON fin.customers
  USING btree(last_name DESC);
ANALYZE fin.customers;

CREATE TABLE IF NOT EXISTS fin.customers_contact_details(
  /***
  A column can be both a primary key (PK) and a foreign key (FK) in a database table. This design
  is used to represent a one-to-one or one-to-zero relationship between two tables, ensuring that
  for every row in the child table there is exactly one corresponding row in the parent table.
  ***/
  id            INT PRIMARY KEY,
                /***
                A foreign key is a column or a set of columns in a database table (the child table)
                that refers a unique constraint in another table (the parent table) establishing a
                link between the two.
                ***/
                CONSTRAINT fk_customers_contact_details_to_customers
                  FOREIGN KEY(id)
                  REFERENCES fin.customers(id)
                  /***
                  Automatically deletes all the referencing rows in the child table when the
                  referenced rows in the parent table are deleted.
                  ***/
                  ON DELETE CASCADE,
  birth_date    DATE NOT NULL  --YYYY-MM-DD
                  CONSTRAINT check_birth_date
                    CHECK(birth_date > '1899-12-31'),
  gender        TEXT NOT NULL,
  address1      TEXT NOT NULL
                  CONSTRAINT check_address1
                    CHECK(TRIM(address1) <> ''),
  address2      TEXT,
  city_name     TEXT NOT NULL
                  CONSTRAINT check_city_name
                    CHECK(TRIM(city_name) <> ''),
  state_name    TEXT NOT NULL
                  CONSTRAINT check_state_name
                    CHECK(TRIM(state_name) <> ''),
  country_name  TEXT NOT NULL
                  CONSTRAINT check_country_name
                    CHECK(TRIM(country_name) <> ''),
  zip_code      TEXT,
  email         TEXT NOT NULL
                  CONSTRAINT check_email
                    CHECK(TRIM(email) <> ''),
  phone         TEXT NOT NULL
                  CONSTRAINT check_phone
                    CHECK(TRIM(phone) <> ''),
  /***
  https://www.postgresql.org/docs/current/datatype-datetime.html
  This stores date and time along with time zone information. PostgreSQL automatically converts the
  timestamp to UTC for storage and adjusts it back based on the current time zone settings when
  queried. 8 bytes in length.
  ***/
  created_at    TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  updated_at    TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS  fin.customers_credentials(
  id               INT PRIMARY KEY,
                   CONSTRAINT fk_customers_credentials_to_customers
                     FOREIGN KEY(id)
                     REFERENCES fin.customers(id)
                     ON DELETE CASCADE,
  user_name        TEXT UNIQUE NOT NULL
                     CONSTRAINT check_user_name
                       CHECK(TRIM(user_name) <> ''),
  password_hash    TEXT NOT NULL
                     CONSTRAINT check_password_hash
                       CHECK(TRIM(password_hash) <> ''),
  is_admin         BOOLEAN NOT NULL DEFAULT FALSE,
  -- Store only failed attempts to save space and speed up queries. This is effective for simple
  -- throttling.
  failed_attempts  INT NOT NULL DEFAULT 0,
  -- When tracking the last attempt time for a task, the initial value should generally be set to
  -- NULL to represent that an attempt has not yet occurred.
  last_attempt     TIMESTAMP WITH TIME ZONE DEFAULT NULL
);

/**************************************************************************************************
               *** DATABASE ROLES AND PRIVILEGES (Table-level privileges) ***
**************************************************************************************************/
GRANT SELECT, INSERT, UPDATE, DELETE ON ALL TABLES IN SCHEMA fin TO admin_role;
GRANT USAGE, SELECT ON ALL SEQUENCES IN SCHEMA fin TO admin_role;

/***
Privileges are only granted/revoked for objects in existence at the time of the GRANT/REVOKE. ALTER
DEFAULT PRIVILEGES allows you to set the privileges that will be applied to objects created in the
future. (It does not affect privileges assigned to already-existing objects.) Privileges can be set
globally (i.e., for all objects created in the current database), or just for objects created in
specified schemas.
***/
ALTER DEFAULT PRIVILEGES IN SCHEMA fin
GRANT SELECT, INSERT, UPDATE, DELETE ON TABLES TO admin_role;
ALTER DEFAULT PRIVILEGES IN SCHEMA fin
GRANT USAGE, SELECT ON SEQUENCES TO admin_role;

/***
Add a row for the default admin account ONLY when the fin.customers table is empty.
***/
WITH admin_customer AS (  -- Common Table Expression (CTE).
  INSERT INTO fin.customers(
    first_name,
    last_name)
  SELECT
    'n/a',  -- Not applicable.
    'n/a'
  WHERE NOT EXISTS(SELECT 1 FROM fin.customers)
  RETURNING id  -- Capture the auto-generated primary key.
)
INSERT INTO fin.customers_credentials(
  id,
  user_name,
  password_hash,
  is_admin)
SELECT
  id,
  'admin',
  crypt('admin', gen_salt('bf', 10)),
  TRUE
FROM admin_customer;

/**************************************************************************************************
                                 *** FUNCTIONS/STORED PROCEDURES ***
**************************************************************************************************/
-- DROP PROCEDURE IF EXISTS fin.add_customer;
CREATE OR REPLACE PROCEDURE fin.add_customer(
  IN p_user_name TEXT,
  IN p_password TEXT,
  IN p_first_name TEXT,
  IN p_middle_name TEXT,
  IN p_last_name TEXT,
  IN p_marketing BOOLEAN,
  IN p_birth_date DATE,
  IN p_gender CHAR,
  IN p_address1 TEXT,
  IN p_address2 TEXT,
  IN p_city TEXT,
  IN p_state TEXT,
  IN p_country TEXT,
  IN p_zip_code TEXT,
  IN p_email TEXT,
  IN p_phone TEXT
)
LANGUAGE PLPGSQL
/***
This delimits the SP body using dollar-quoted string constants, which avoids the need to escape
single quotes within the code.
***/
AS $$
DECLARE
  c_id INT;
  pwd_hash TEXT;
--This block encloses the executable logic of the stored procedure's body.
BEGIN
  -- First, insert into the 'customers' table and return the new id.
  INSERT INTO fin.customers(
    first_name,
    middle_name,
    last_name,
    marketing_consent)
  VALUES(
    p_first_name,
    p_middle_name,
    p_last_name,
    p_marketing)
  RETURNING id INTO c_id;  -- This makes the new 'id' available to the next statements.
  INSERT INTO fin.customers_contact_details(
    id,
    birth_date,
    gender,
    address1,
    address2,
    city_name,
    state_name,
    country_name,
    zip_code,
    email,
    phone)
  VALUES(
    c_id,
    p_birth_date,
    p_gender,
    p_address1,
    p_address2,
    p_city,
    p_state,
    p_country,
    p_zip_code,
    p_email,
    p_phone);
  pwd_hash := crypt(p_password, gen_salt('bf', 10));
  INSERT INTO fin.customers_credentials(
    id,
    user_name,
    password_hash)
  VALUES(
    c_id,
    p_user_name,
    pwd_hash);
EXCEPTION  -- https://www.postgresql.org/docs/current/errcodes-appendix.html
  WHEN unique_violation THEN
    -- Re-raise the exception to inform the caller and ensure rollback of the transaction.
    RAISE EXCEPTION 'Username (%) is already taken. Please choose another one.', p_user_name;
  WHEN check_violation THEN
    RAISE EXCEPTION '% -- %', SQLSTATE, SQLERRM;
  WHEN OTHERS THEN
    RAISE EXCEPTION '% -- %', SQLSTATE, SQLERRM;
END;
$$;

/***
Return values:
  -2 -- Invalid (unknown) user.
  -1 -- Authentication failed.
   0 -- User authenticated.
***/
CREATE OR REPLACE PROCEDURE fin.authenticate_user(
  IN puser_name TEXT,
  IN ppassword TEXT,
  IN correlation_id TEXT,
  OUT pout INT,
  OUT pis_admin BOOL
)
LANGUAGE PLPGSQL
AS $$
DECLARE
  vmax_attempts CONSTANT INT := 3;
  vbase_delay CONSTANT NUMERIC := 1;  -- Seconds.
  vmin_delay CONSTANT NUMERIC := 0.5;  -- Seconds (500ms).
  vmax_delay CONSTANT NUMERIC := 30;  -- Seconds.
  vunknown_user_delay CONSTANT NUMERIC := 5;  -- Seconds.
  vsleep_time NUMERIC := 0;
  vfailed_attempts INT;
  vhash TEXT;
BEGIN
  pis_admin := false;
  SELECT
    password_hash
  INTO
    vhash
  FROM fin.customers_credentials
  WHERE user_name = puser_name;
  /***
  The SELECT INTO statement in Postgres (without the STRICT keyword) sets a special variable
  called FOUND to TRUE if a row is returned, and FALSE if no row is found. No exception is raised
  automatically in this case, allowing you to check the condition and handle it gracefully.
  ***/
  IF NOT FOUND THEN  -- User not found; password can't be valid.
    pout := -2;
    RAISE NOTICE 'Unknown user (%); sleeping for % seconds', puser_name, vunknown_user_delay;
    -- Pause the current session for the specified number of seconds (can be a decimal).
    PERFORM pg_sleep(vunknown_user_delay);
  /***
  Use the crypt function to compare the provided password against the stored hash;
  the crypt function uses the salt from the stored hash to perform the comparison.
  ***/
  ELSIF vhash = crypt(ppassword, vhash) THEN  -- User authenticated.
    pout := 0;
    UPDATE fin.customers_credentials
    SET
      failed_attempts = 0,
      last_attempt = CURRENT_TIMESTAMP
    WHERE user_name = puser_name
    RETURNING is_admin INTO pis_admin;
  ELSE  -- Authentication failed.
    pout := -1;
    UPDATE fin.customers_credentials
    SET
      failed_attempts = failed_attempts + 1,
      last_attempt = CURRENT_TIMESTAMP
    WHERE user_name = puser_name
    RETURNING failed_attempts INTO vfailed_attempts;
    IF vfailed_attempts > vmax_attempts THEN
      /***
      The basic algorithm increases the wait time between each retry attempt, often calculated as
      base_delay * (2 ^ failed_attempts).
      Exponential backoff: Delay doubles each time (e.g., 1s, 2s, 4s, 8s, 16s etc.), usually up to
      a defined maximum.
      Add random noise (jitter) to the exponential backoff time to prevent multiple clients from
      retrying simultaneously, which can cause a "thundering herd."
      The random() function supports a (min, max) syntax for integers and numeric types, returning
      a random value within the specified INCLUSIVE range.
      ***/
      vsleep_time := RANDOM(vmin_delay,
        LEAST(vmax_delay, CAST(vbase_delay * POWER(2, vfailed_attempts) AS NUMERIC)));
      RAISE NOTICE 'Attempt % failed; blocking for % seconds...', vfailed_attempts, vsleep_time
        USING DETAIL = correlation_id;
      PERFORM pg_sleep(vsleep_time);
    END IF;
  END IF;
END;
$$;

CREATE OR REPLACE PROCEDURE fin.change_password(
  IN p_user_name TEXT,
  IN p_old_password TEXT,
  IN p_new_password TEXT,
  OUT ret BOOL
)
LANGUAGE PLPGSQL
/***
This delimits the SP body using dollar-quoted string constants, which avoids the need to escape single quotes within the code.
***/
AS $$
DECLARE
  vpwd_hash TEXT;
--This block encloses the executable logic of the stored procedure's body.
BEGIN
  ret := false;
  SELECT
    password_hash
  INTO
    vpwd_hash
  FROM fin.customers_credentials
  WHERE user_name = p_user_name;
  IF NOT FOUND THEN  -- User not found; password can't be valid.
    RAISE NOTICE 'Unknown user (%); password not changed.', p_user_name;
  ELSIF vpwd_hash = crypt(p_old_password, vpwd_hash) THEN  -- User authenticated.
    vpwd_hash := crypt(p_new_password, gen_salt('bf', 10));
    UPDATE fin.customers_credentials
    SET
      password_hash = vpwd_hash
    WHERE user_name = p_user_name;
    ret := true;
  END IF;
END;
$$;

/**************************************************************************************************
                            *** TRIGGER FUNCTIONS/STORED PROCEDURES ***
**************************************************************************************************/
CREATE OR REPLACE FUNCTION fin.customers_block_row_deletion()
RETURNS TRIGGER
LANGUAGE PLPGSQL
AS $$
BEGIN
  /***
  In a DELETE trigger, the special variable NEW is always NULL. You must use OLD.column_name to
  grab the data from the row being removed.
  ***/
  IF OLD.id = 1 THEN
    RAISE EXCEPTION 'Deletion not allowed: This row is protected.';
  END IF;
  -- Crucial: Return OLD so Postgres proceeds with the deletion.
  RETURN OLD;
END;
$$;

CREATE OR REPLACE TRIGGER trg_customers_prevent_delete
BEFORE DELETE ON fin.customers
FOR EACH ROW
EXECUTE FUNCTION fin.customers_block_row_deletion();

CREATE OR REPLACE FUNCTION fin.customers_credentials_block_row_deletion()
RETURNS TRIGGER
LANGUAGE PLPGSQL
AS $$
BEGIN
  IF OLD.is_admin = TRUE THEN
    RAISE EXCEPTION 'Deletion not allowed: This row is protected.';
  END IF;
  RETURN OLD;
END;
$$;

CREATE OR REPLACE TRIGGER trg_customers_credentials_prevent_delete
BEFORE DELETE ON fin.customers_credentials
FOR EACH ROW
EXECUTE FUNCTION fin.customers_credentials_block_row_deletion();







/*************************************************************************************************/

-- CREATE OR REPLACE FUNCTION fin.get_customers_contact_details()
-- RETURNS SETOF fin.customers_contact_details
-- LANGUAGE PLPGSQL
-- AS $$
-- BEGIN
--   RETURN QUERY SELECT * FROM fin.customers_contact_details;
-- END;
-- $$;



-- CREATE OR REPLACE PROCEDURE fin.customer_contact_details(
--   IN puser_name TEXT
-- )
-- LANGUAGE PLPGSQL
-- AS $$
-- BEGIN
-- END;
-- $$;

-- CREATE OR REPLACE PROCEDURE fin.customers_credentials()
-- LANGUAGE PLPGSQL
-- AS $$
-- BEGIN
-- END;
-- $$;

-- CREATE OR REPLACE PROCEDURE fin.customer_credentials(
--   IN puser_name TEXT
-- )
-- LANGUAGE PLPGSQL
-- AS $$
-- BEGIN
-- END;
-- $$;
