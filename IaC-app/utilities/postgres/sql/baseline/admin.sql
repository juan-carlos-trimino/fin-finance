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

-- ************************************************************************************************
-- Create the tables
-- ************************************************************************************************
-- customers-to-customer_contact_details Relationship: One-to-One
-- customers-to-credentials Relationship: One-to-One
-- URL: /register
CREATE TABLE IF NOT EXISTS fin.customers(
  id                 INT PRIMARY KEY GENERATED ALWAYS ASW IDENTITY,
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
  id            INT PRIMARY KEY
                  CONSTRAINT fk_customer_contact_details_to_customers
                    FOREIGN KEY(id)
                    REFERENCES xxxadmin.customers(id)
                    ON DELETE CASCADE,
  birth_date    DATE NOT NULL  --YYYY-MM-DD
                  CONSTRAINT check_birth_date
                    CHECK(birth_date > '1899-12-31'),
  gender        CHAR NOT NULL
                  CONSTRAINT check_gender
                    CHECK(gender IN('F', 'M')),
  address_1     TEXT NOT NULL,
  address_2     TEXT,
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
  id             INT PRIMARY KEY
                   CONSTRAINT fk_credentials_to_customers
                     FOREIGN KEY(id)
                     REFERENCES xxxadmin.customers(id)
                     ON DELETE CASCADE,
  userid         TEXT UNIQUE NOT NULL,
  password_hash  TEXT NOT NULL,
  is_active      BOOLEAN NOT NULL DEFAULT FALSE,
  --  failed_login_attempts
  last_login     TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);
