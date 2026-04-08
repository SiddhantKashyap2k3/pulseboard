-- migrations/001_create_users.sql

-- CREATE TABLE defines a new table in the database
-- IF NOT EXISTS means don't error if the table already exists
CREATE TABLE IF NOT EXISTS users (

    -- id is the primary key — a unique number for each row
    -- SERIAL means Postgres automatically assigns the next number (1, 2, 3...)
    -- PRIMARY KEY means this column uniquely identifies each row
    id SERIAL PRIMARY KEY,

    -- email must be unique across all rows — no two users share an email
    -- NOT NULL means this field is required — cannot be empty
    email VARCHAR(255) UNIQUE NOT NULL,

    -- we never store plain passwords — only the bcrypt hash
    -- VARCHAR(255) means a text string up to 255 characters
    password_hash VARCHAR(255) NOT NULL,

    -- TIMESTAMP stores date and time
    -- DEFAULT NOW() automatically sets it to the current time when a row is inserted
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);