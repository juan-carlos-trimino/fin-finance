-- Online Banking System.
SELECT 'Output from script, run began at: ' AS "Script Information",
  NOW() AS "Date and Time Executed";

-- *******************
-- Create the database
-- *******************
-- Important - You should back up the master database each time that you create, modify, or drop a
-- database.
CREATE DATABASE db_check_register
WITH
  ALLOW_CONNECTIONS = TRUE
  CONNECTION_LIMIT = -1  -- Unlimited connections.
  ENCODING = 'UTF8'
  LC_COLLATE = 'C.UTF8'  -- Determine the sort order of strings.
  LC_CTYPE = 'C.UTF8'  -- Define character classification rules.
  IS_TEMPLATE = FALSE  -- Only superusers or the database owner can clone the database.
  TEMPLATE = 'template0';

-- Connect to the database.
\c db_check_register

\qecho 'Current database version:'
SELECT version();

-- ************************************************************************************************
-- Create the schemas
-- ************************************************************************************************
CREATE SCHEMA IF NOT EXISTS customers;
CREATE SCHEMA IF NOT EXISTS accounts;

-- ********
-- Security - Create an admin-level role for the database.
-- ********
-- Any role (user) with the explicitly granted LOGIN attribute can connect to the database.
-- Source - https://stackoverflow.com/a/55954480
-- Posted by Pali
-- Retrieved 2026-01-03, License - CC BY-SA 4.0
--
-- Dollar quoting in PostgreSQL is a method for creating string constants without the need to
-- escape single quotes, backslashes, or other special characters, significantly improving code
-- readability, especially for complex queries or function bodies. Dollar quoting is not part of
-- the SQL standard.
-- The PostgreSQL DO statement executes an anonymous code block in a procedural language, such as
-- PL/pgSQL. This is useful for running one-off, transient code blocks that do not need to be
-- stored permanently as a function or stored procedure.
DO $$
BEGIN
  CREATE ROLE trimino WITH LOGIN PASSWORD 'trimino';
  EXCEPTION WHEN duplicate_object THEN
    RAISE NOTICE '%, skipping', SQLERRM USING ERRCODE = SQLSTATE;
END
$$;
--CREATE ROLE trimino WITH LOGIN PASSWORD 'trimino';
-- It allows the role to connect to the specified database.
GRANT CONNECT ON DATABASE db_check_register TO trimino;
-- It revokes the connect privilege for all other roles to the database. In Postgres, PUBLIC is a
-- default role that includes all users.
REVOKE CONNECT ON DATABASE db_check_register FROM PUBLIC;
--
GRANT USAGE ON SCHEMA public TO trimino;
GRANT USAGE ON SCHEMA customers TO trimino;
GRANT USAGE ON SCHEMA accounts TO trimino;
--
GRANT SELECT, INSERT, UPDATE, DELETE ON ALL TABLES
  IN SCHEMA public TO trimino;
GRANT SELECT, INSERT, UPDATE, DELETE ON ALL TABLES
  IN SCHEMA customers TO trimino;
GRANT SELECT, INSERT, UPDATE, DELETE ON ALL TABLES
  IN SCHEMA accounts TO trimino;

-- Once connected, set the search path to look for objects in your schema first, and if not found,
-- to fall back to the default public schema.
SET search_path TO customers, public;

-- ************************************************************************************************
-- Create the tables
-- ************************************************************************************************
-- Customers - Accounts Relationship
-- One-to-Many Relationship - One customer is allowed to create many accounts. This relationship is
-- established using foreign keys, where the primary key of the "one" table is referenced as a
-- foreign key in the "many" table.
CREATE TABLE IF NOT EXISTS customers.tbl_customers(
  -- Version UUIDv7 includes a Unix timestamp in its value, making it more index-friendly because
  -- UUIDs generated close in time are stored nearer to each other in the index and can be accessed
  -- faster.
  id               UUID PRIMARY KEY DEFAULT uuidv7(),
  username         VARCHAR(64) UNIQUE NOT NULL,
  password_hash    VARCHAR(256) NOT NULL,
  first_name       VARCHAR(64) NOT NULL,
  middle_name      VARCHAR(64),
  last_name        VARCHAR(64) NOT NULL,
  birth_date       DATE NOT NULL  --YYYY-MM-DD
                     CONSTRAINT check_birth_date
                       CHECK(birth_date > '1899-12-31'),
  gender           CHAR NOT NULL
                     CONSTRAINT check_gender
                       CHECK(gender IN('F', 'M')),
  address_1        VARCHAR(128) NOT NULL,
  address_2        VARCHAR(128),
  city_name        VARCHAR(128) NOT NULL,
  state_name       VARCHAR(128) NOT NULL,
  country_name     VARCHAR(64) NOT NULL,
  zip_code         VARCHAR(16),
  primary_email    VARCHAR(128),
  secondary_email  VARCHAR(128),
  primary_phone    VARCHAR(16) NOT NULL,
  secondary_phone  VARCHAR(16),
  -- https://www.postgresql.org/docs/current/datatype-datetime.html
  -- This stores date and time along with time zone information. PostgreSQL automatically converts
  -- the timestamp to UTC for storage and adjusts it back based on the current time zone settings
  -- when queried. 8 bytes in length.
  created_at       TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  updated_at       TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

SET search_path TO accounts, public;

-- Accounts - Registers Relationship
-- Many-to-Many Relationship - An account can have multiple transactions and a transaction can
--                             involve multiple account numbers.
CREATE TABLE IF NOT EXISTS accounts.tbl_accounts(
  id           UUID PRIMARY KEY DEFAULT uuidv7(),
  -- Crucially, this foreign key column must also have a UNIQUE constraint applied to it. This
  -- UNIQUE constraint ensures that no two records in the child table (tbl_accounts) can point
  -- to the same record in the parent table (tbl_customers), thus enforcing the "one-to-one" aspect.
  customer_id  UUID UNIQUE NOT NULL,
  acct_name    VARCHAR(64) UNIQUE NOT NULL,
  -- A column defined with a DEFAULT value and a CHECK constraint still needs an explicit NOT NULL
  -- constraint if you want to prevent NULL values.
  -- Here's why:
  -- DEFAULT Constraint: This constraint only applies if no value is provided during an INSERT
  -- operation. If you explicitly INSERT a NULL value into a column, the default value is bypassed,
  -- and the NULL is inserted (unless NOT NULL is present).
  -- CHECK Constraint: In SQL, any comparison involving a NULL value evaluates to UNKNOWN, not TRUE
  -- or FALSE. A CHECK constraint only fails if the condition evaluates to FALSE. If it evaluates
  -- to TRUE or UNKNOWN (due to a NULL value), the constraint is satisfied and the NULL is
  -- accepted.
  -- NOT NULL Constraint: This is the only constraint specifically designed to enforce the presence
  -- of data and prohibit NULL values in a column.
  acct_type    VARCHAR(16) NOT NULL DEFAULT 'checking'
                 CONSTRAINT check_account_type
                   CHECK(acct_type IN('checking', 'savings')),
  acct_status  VARCHAR(16) NOT NULL DEFAULT 'active'
                 CONSTRAINT check_account_status
                   CHECK(acct_status IN('active', 'closed', 'suspended')),
  created_at   TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  updated_at   TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
               -- A foreign key is a column or a set of columns in a database table (the
               -- "child table") that refers a unique constraint in another table (the
               -- "parent table") establishing a link between the two.
               CONSTRAINT fk_accounts_to_customers
                 FOREIGN KEY(customer_id)
                 REFERENCES customers.tbl_customers(id)
                 -- Automatically deletes all the referencing rows in the child table
                 -- (tbl_accounts) when the referenced rows in the parent table (tbl_customers)
                 -- are deleted.
                 ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS accounts.tbl_register_entries(
  id              UUID PRIMARY KEY DEFAULT uuidv7(),
  acct_id         UUID UNIQUE NOT NULL,
  check_number    INT UNIQUE,
  payment_date    DATE NOT NULL,  --YYYY-MM-DD
  payee           VARCHAR(64) NOT NULL,
  tr_type         VARCHAR(16) NOT NULL
                    CONSTRAINT check_tr_type
                      CHECK(tr_type IN('deposit', 'debit')),
  -- There is no difference between NUMERIC and DECIMAL in PostgreSQL.
  amount          NUMERIC(12, 2) NOT NULL
                    CONSTRAINT check_amount
                      CHECK(amount > 0),
  cleared         BOOLEAN NOT NULL DEFAULT FALSE,
  tr_description  VARCHAR(128),
  created_at      TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  updated_at      TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
                  CONSTRAINT fk_transactions_to_accounts
                    FOREIGN KEY(acct_id)
                    REFERENCES accounts.tbl_accounts(id)
                    ON DELETE CASCADE
);




/**************************************************************************************************
Create the stored procedures
***************************************************************************************************
This statement will create the procedure if it doesn't exist, or replace its definition if it does.
This is a more concise way to update or create a procedure.
**************************************************************************************************/
-- CREATE OR REPLACE PROCEDURE customers.sp_customer(
--  IN pusername VARCHAR(64),
--  IN ppassword_hash VARCHAR(256),
-- --  IN pcustomer_type VARCHAR(16) DEFAULT 'regular',
--  INOUT pcustomer_id INT DEFAULT -1)
--  /***
--  In PostgreSQL, when discussing "procedure language SQL vs PostgreSQL," the distinction lies in
--  the type of functions or procedures you can create within the database.
--  When to choose which:
--   * Use LANGUAGE SQL when your function primarily involves data manipulation or retrieval that can
--     be expressed efficiently through SQL statements, and complex procedural logic is not required.
--   * Use LANGUAGE PL/pgSQL when you need advanced procedural capabilities, such as loops,
--     conditional logic, error handling, or dynamic SQL generation, that are not possible or
--     practical with pure SQL functions.
--  ***/
-- LANGUAGE PLPGSQL
-- /***
-- This defines the procedure's body using dollar-quoted string constants. This avoids the need to
-- escape single quotes within the procedure's code.
-- ***/
-- AS $$
--   /***
--   (Optional) This section within the procedure body is used to declare local variables that are
--   only accessible within the procedure.
--   ***/
--   -- DECLARE
--   /***
--   This block encloses the SQL statements and control structures that constitute the procedure's
--   logic.
--   ***/
-- BEGIN
--   -- INSERT INTO bs.tbl_customer(customer_type, username, password_hash)
--   INSERT INTO customers.tbl_accounts(username, password_hash)
--   --  VALUES(pcustomer_type, pusername, ppassword_hash)
--    VALUES( pusername, ppassword_hash)
--    RETURNING customer_id INTO pcustomer_id;
-- END;
-- /***
-- The final semicolon marks the end of the CREATE PROCEDURE statement.
-- ***/
-- $$;





-- CREATE OR REPLACE PROCEDURE customers.sp_customer_info(
--   IN username VARCHAR(64),
--   IN password_hash VARCHAR(256),
--   -- IN customer_type VARCHAR(16),
--   IN first_name VARCHAR(64),
--   IN middle_name VARCHAR(64),
--   IN last_name VARCHAR(64),
--   IN date_of_birth DATE,
--   IN tax_identifier VARCHAR(16),
--   IN address_1 VARCHAR(128),
--   IN address_2 VARCHAR(128),
--   IN city_name VARCHAR(128),
--   IN state_name VARCHAR(128),
--   IN country_name VARCHAR(64),
--   IN zip_code VARCHAR(16),
--   IN primary_email VARCHAR(256),
--   IN secondary_email VARCHAR(256),
--   IN primary_phone VARCHAR(16),
--   IN secondary_phone VARCHAR(16)
-- )
-- LANGUAGE PLPGSQL
-- AS $$
-- DECLARE
--  customer_id INT = -1;
-- BEGIN
--   -- CALL bs.sp_customer(username, password_hash, customer_type, customer_id);
--   CALL customers.sp_customer(username, password_hash,  customer_id);
--   INSERT INTO customers.tbl_customer_info (customer_id, first_name, middle_name, last_name, date_of_birth, tax_identifier,
--    address_1, address_2, city_name, state_name, country_name, zip_code, primary_email,
--    secondary_email, primary_phone, secondary_phone)
--    VALUES(customer_id, first_name, middle_name, last_name, date_of_birth, tax_identifier, address_1, address_2,
--     city_name, state_name, country_name, zip_code, primary_email, secondary_email, primary_phone,
--     secondary_phone);
--   COMMIT;
-- END;
-- $$;


-- CREATE OR REPLACE FUNCTION customers.fn_get_all_customers()
-- RETURNS TABLE (id INT, name TEXT)
-- LANGUAGE plpgsql
-- AS $$
-- BEGIN
--   RETURN QUERY SELECT customer_id, customer_type, username, password_hash, created_at, updated_at FROM customers.tbl_accounts;
-- END;
-- $$;



SELECT 'Output from script, run ended at: ' AS "Script Information",
  NOW() AS "Date and Time Executed";
