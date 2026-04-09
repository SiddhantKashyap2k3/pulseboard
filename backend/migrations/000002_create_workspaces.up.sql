-- A workspace is one tenant/company in PulseBoard
CREATE TABLE IF NOT EXISTS workspaces (
    id         SERIAL PRIMARY KEY,

    -- name of the company or project e.g. "Acme Corp"
    name       VARCHAR(255) NOT NULL,

    -- api_key is what external apps send to authenticate event ingestion
    -- we generate a random string and store it here
    api_key    VARCHAR(255) UNIQUE NOT NULL,

    -- which user owns this workspace
    -- REFERENCES users(id) is a foreign key —
    -- it means this value must exist in the users table's id column
    -- ON DELETE CASCADE means if the user is deleted, their workspace is too
    user_id    INTEGER REFERENCES users(id) ON DELETE CASCADE,

    created_at TIMESTAMP DEFAULT NOW()
);