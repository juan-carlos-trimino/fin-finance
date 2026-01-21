-- Online Banking System.
SELECT 'Output from script, run began at: ' AS "Script Information",
  NOW() AS "Date and Time Executed";

-- *******************
-- Create the database
-- *******************
-- Important - You should back up the master database each time that you create, modify, or drop a
-- database.
CREATE DATABASE finances
WITH
  ALLOW_CONNECTIONS = TRUE
  CONNECTION_LIMIT = -1  -- Unlimited connections.
  ENCODING = 'UTF8'
  LC_COLLATE = 'C.UTF8'  -- Determine the sort order of strings.
  LC_CTYPE = 'C.UTF8'  -- Define character classification rules.
  IS_TEMPLATE = FALSE  -- Only superusers or the database owner can clone the database.
  TEMPLATE = 'template0';

-- Connect to the database.
\c finances

\qecho 'Current database version:'
SELECT version();

-- ************************************************************************************************
-- Create the schemas
-- ************************************************************************************************
CREATE SCHEMA IF NOT EXISTS fin;

-- Once connected, set the search path to look for objects in your schema first, and if not found,
-- to fall back to the default public schema.
SET search_path TO fin, public;

-- ************************************************************************************************
-- Create the tables
-- ************************************************************************************************
-- customers-to-customer_contact_details Relationship: One-to-One
-- customers-to-credentials Relationship: One-to-One
-- URL: /register
CREATE TABLE IF NOT EXISTS fin.customers(
  id                 INT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
  first_name         TEXT NOT NULL,
  middle_name        TEXT,
  last_name          TEXT NOT NULL,
  marketing_consent  BOOLEAN NOT NULL DEFAULT FALSE,
  -- https://www.postgresql.org/docs/current/datatype-datetime.html
  -- This stores date and time along with time zone information. PostgreSQL automatically converts
  -- the timestamp to UTC for storage and adjusts it back based on the current time zone settings
  -- when queried. 8 bytes in length.
  created_at         TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  updated_at         TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS fin.customer_contact_details(
  -- A column can be both a primary key (PK) and a foreign key (FK) in a database table. This
  --design is used to represent a one-to-one or one-to-zero relationship between two tables,
  --ensuring that for every row in the child table there is exactly one corresponding row in the
  --parent table.
  id            INT PRIMARY KEY,
                CONSTRAINT fk_customer_contact_details_to_customers
                  FOREIGN KEY(id)
                  REFERENCES fin.customers(id)
                  ON DELETE CASCADE,
  birth_date    DATE NOT NULL  --YYYY-MM-DD
                  CONSTRAINT check_birth_date
                    CHECK(birth_date > '1899-12-31'),
  gender        TEXT NOT NULL,
  address1      TEXT NOT NULL,
  address2      TEXT,
  city_name     TEXT NOT NULL,
  state_name    TEXT NOT NULL,
  country_name  TEXT NOT NULL,
  zip_code      TEXT,
  email         TEXT NOT NULL,
  phone         TEXT NOT NULL,
  -- https://www.postgresql.org/docs/current/datatype-datetime.html
  -- This stores date and time along with time zone information. PostgreSQL automatically converts
  -- the timestamp to UTC for storage and adjusts it back based on the current time zone settings
  -- when queried. 8 bytes in length.
  created_at    TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  updated_at    TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS  fin.credentials(
  id             INT PRIMARY KEY,
                 CONSTRAINT fk_credentials_to_customers
                   FOREIGN KEY(id)
                   REFERENCES fin.customers(id)
                   ON DELETE CASCADE,
  user_name      TEXT UNIQUE NOT NULL,
  password_hash  TEXT NOT NULL,
  is_active      BOOLEAN NOT NULL DEFAULT FALSE,
  --  failed_login_attempts
  last_login     TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- ************************************************************************************************
-- Create functions/stored procedures
-- ************************************************************************************************
--In PostgreSQL, stored procedures are not atomic by default but can manage their own transactions
--explicitly to ensure atomicity. This is a key difference from PostgreSQL functions, which are
--always atomic and run within a single, implicit transaction.
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
--This block encloses the executable logic of the stored procedure's body.
BEGIN
  -- First, insert into the 'customers' table and return the new id.
  INSERT INTO fin.customers(first_name, middle_name, last_name, marketing_consent)
    VALUES(p_first_name, p_middle_name, p_last_name, p_marketing)
    RETURNING id INTO c_id;  -- This makes the new 'id' available to the next statements.
  INSERT INTO fin.customer_contact_details(id, birth_date, gender, address1, address2,
    city_name, state_name, country_name, zip_code, email, phone)
    VALUES(c_id, p_birth_date, p_gender, p_address1, p_address2, p_city,
      p_state, p_country, p_zip_code, p_email, p_phone);
  INSERT INTO fin.credentials(id, user_name, password_hash, is_active)
    VALUES(c_id, p_user_name, p_password, TRUE);
EXCEPTION  -- https://www.postgresql.org/docs/current/errcodes-appendix.html
  WHEN unique_violation THEN
    RAISE EXCEPTION 'Username (%) is already taken. Please choose another one.', p_user_name;
  WHEN OTHERS THEN
    RAISE EXCEPTION '% -- %', SQLSTATE, SQLERRM;
END;
-- The final semicolon marks the end of the CREATE PROCEDURE statement.
$$;


/*
After a COMMIT or ROLLBACK is issued inside a procedure, a new transaction is automatically started, so you do not need a separate START TRANSACTION command.
In procedures invoked by the CALL command as well as in anonymous code blocks (DO command), it is possible to end transactions using the commands COMMIT and ROLLBACK. A new transaction is started automatically after a transaction is ended using these commands, so there is no separate START TRANSACTION command.

\set PROMPT1 '\t\t\t\t\t\t>>>%`date +%H:%M:%S`<<<\n%/%R%# '














*/
