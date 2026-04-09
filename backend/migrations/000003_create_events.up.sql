-- An event is one analytics data point sent by a client app
CREATE TABLE IF NOT EXISTS events (
    id           SERIAL PRIMARY KEY,

    -- which workspace this event belongs to
    workspace_id INTEGER REFERENCES workspaces(id) ON DELETE CASCADE,

    -- the event name e.g. "user_signed_up", "page_view", "api_call"
    name         VARCHAR(255) NOT NULL,

    -- properties is flexible JSON data — any extra info the client sends
    -- JSONB is Postgres's binary JSON type — faster to query than plain JSON
    properties   JSONB DEFAULT '{}',

    -- when the event actually happened (set by the client)
    occurred_at  TIMESTAMP DEFAULT NOW(),

    created_at   TIMESTAMP DEFAULT NOW()
);

-- An index speeds up queries that filter by workspace_id
-- Without this, Postgres scans every row to find matching events
-- With this, it jumps directly to the matching rows
-- We'll query "all events for workspace X" constantly so this is important
CREATE INDEX IF NOT EXISTS idx_events_workspace_id ON events(workspace_id);
CREATE INDEX IF NOT EXISTS idx_events_occurred_at  ON events(occurred_at);