--*************************************************************************************************
-- Notes
-- * By default, Postgres ignores case and always turns every identifier into lowercase. Postgres
--   preserves case only when the identifier is double quoted (e.g., "id"), and double-quoted
--   identifiers, known as 'delimited identifiers', are case sensitive.
-- * In Postgres, you can CAST between types with the :: shorthand.
-- * In PostgreSQL, stored procedures are not atomic by default but can manage their own
--   transactions explicitly to ensure atomicity. This is a key difference from PostgreSQL
--   functions, which are always atomic and run within a single, implicit transaction.
-- * Postgres executes a function atomically and transactionally; i.e, if the function fails at any
--   step during its execution, all previous changes made within that function are rolled back,
--   ensuring data integrity and consistency.
-- Notes (Indexes)
-- * Postgres defaults to a B-tree index unless another data strcuture is specified. B-tree
--   supports a wide range of data types and can be used for both equality searches (=) and range
--   queries using greater than (>, >=), less than (<, <=), and BETWEEN operators.
-- * It's good practice to run the ANALYZE command after creating an index.
--   When a new index is created, the query planner may not take advantage of it immediately and
--   may continue relying on previously collected statistics. Running ANALYZE updates the table
--   statistics, allowing the planner to take full advantage of the new index.
-- * Postgres automatically creates an index for the primary key column.
-- * Indexes can significantly optimize query performance, but they don't come for free. Each time
--   we create a new index, Postgres must maintain it by updating its structure whenever the value
--   of an indexed column changes in the primary table.
--   Apart from the index maintenance aspects, the more indexes Postgres has, the more time it will
--   spend in the planning phase while selecting and generating the most efficient execution plan
--   for a query.
--*************************************************************************************************
-- Online Banking System.
SELECT 'Output from script, run began at: ' AS "Script Information",
  NOW() AS "Date and Time Executed";

-- ************************************************************************************************
-- Create the database
-- ************************************************************************************************
CREATE DATABASE finances
WITH
  ALLOW_CONNECTIONS = TRUE
  CONNECTION_LIMIT = -1  -- Unlimited connections.
  ENCODING = 'UTF8'
  LC_COLLATE = 'C.UTF8'  -- Determine the sort order of strings.
  LC_CTYPE = 'C.UTF8'  -- Define character classification rules.
  IS_TEMPLATE = FALSE  -- Only superusers or the database owner can clone the database.
  TEMPLATE = 'template0';

-- ************************************************************************************************
-- Connect to the database.
-- ************************************************************************************************
\c finances

\qecho 'Current database version:'
SELECT version();

-- ************************************************************************************************
-- To use bcrypt in Postgres, you can utilize the pgcrypto extension. This extension provides the
-- crypt() and gen_salt() functions necessary for secure password hashing and verification within
-- SQL. Enable the extension in your specific database by running the following SQL command.
-- ************************************************************************************************
CREATE EXTENSION IF NOT EXISTS pgcrypto;

-- ************************************************************************************************
-- Create the schemas
-- ************************************************************************************************
CREATE SCHEMA IF NOT EXISTS fin;

-- ************************************************************************************************
-- Once connected, set the search path to look for objects in your schema first, and if not found,
-- to fall back to the default public schema.
-- ************************************************************************************************
SET search_path TO fin, public;

-- ************************************************************************************************
-- Create the tables
-- customers-to-customer_contact_details Relationship: One-to-One
-- customers-to-credentials Relationship: One-to-One
-- URL: /register
-- ************************************************************************************************
CREATE TABLE IF NOT EXISTS fin.customers(
  id                 INT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
  first_name         TEXT NOT NULL
  -- Block empty strings ('') and strings with only blanks (' ').
                       CONSTRAINT check_first_name
                         CHECK(TRIM(first_name) <> ''),
  middle_name        TEXT,
  last_name          TEXT NOT NULL
                       CONSTRAINT check_last_name
                         CHECK(TRIM(last_name) <> ''),
  marketing_consent  BOOLEAN NOT NULL DEFAULT FALSE,
  -- https://www.postgresql.org/docs/current/datatype-datetime.html
  -- This stores date and time along with time zone information. PostgreSQL automatically converts
  -- the timestamp to UTC for storage and adjusts it back based on the current time zone settings
  -- when queried. 8 bytes in length.
  created_at         TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  updated_at         TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_last_name
  ON fin.customers
  USING btree(last_name DESC);
ANALYZE fin.customers;

CREATE TABLE IF NOT EXISTS fin.customers_contact_details(
  -- A column can be both a primary key (PK) and a foreign key (FK) in a database table. This
  --design is used to represent a one-to-one or one-to-zero relationship between two tables,
  --ensuring that for every row in the child table there is exactly one corresponding row in the
  --parent table.
  id            INT PRIMARY KEY,
                CONSTRAINT fk_customers_contact_details_to_customers
                  FOREIGN KEY(id)
                  REFERENCES fin.customers(id)
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
  -- https://www.postgresql.org/docs/current/datatype-datetime.html
  -- This stores date and time along with time zone information. PostgreSQL automatically converts
  -- the timestamp to UTC for storage and adjusts it back based on the current time zone settings
  -- when queried. 8 bytes in length.
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
  -- Store only failed attempts to save space and speed up queries. This is effective for simple
  -- throttling.
  failed_attempts  INT NOT NULL DEFAULT 0,
  -- When tracking the last attempt time for a task, the initial value should generally be set to
  -- NULL to represent that an attempt has not yet occurred.
  last_attempt     TIMESTAMP WITH TIME ZONE DEFAULT NULL
);

-- ************************************************************************************************
-- Create functions/stored procedures
-- ************************************************************************************************
-- DROP PROCEDURE IF EXISTS fin.add_customer;
CREATE OR REPLACE PROCEDURE fin.add_customer( --Require transaction########################
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
  INSERT INTO fin.customers_credentials(
    id,
    user_name,
    password_hash)
  VALUES(
    c_id,
    p_user_name,
    p_password);
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
  OUT pout INT
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
    WHERE user_name = puser_name;
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
















    -- UPDATE fin.credentials
    -- SET failed_attempts = vfailed_attempts, last_attempt = CURRENT_TIMESTAMP
    -- WHERE EXISTS (
    --   -- In SQL, SELECT 1 has two primary meanings depending on its context: as a constant value in
    --   -- the result set, or, more commonly, as an efficient way to check for the existence of a row
    --   -- in a subquery.
    --   SELECT 1
    --   FROM fin.credentials
    --   WHERE user_name = puser_name
    -- );





/*
After a COMMIT or ROLLBACK is issued inside a procedure, a new transaction is automatically started, so you do not need a separate START TRANSACTION command.
In procedures invoked by the CALL command as well as in anonymous code blocks (DO command), it is possible to end transactions using the commands COMMIT and ROLLBACK. A new transaction is started automatically after a transaction is ended using these commands, so there is no separate START TRANSACTION command.


*/
