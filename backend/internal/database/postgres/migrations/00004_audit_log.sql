CREATE TYPE audit_log_event_type AS ENUM (
    'other',
    'created',
    'updated',
    'archived'
);

CREATE TABLE IF NOT EXISTS audit_log_entries (
    id TEXT NOT NULL PRIMARY KEY,
    resource_type TEXT NOT NULL,
    relevant_id TEXT NOT NULL DEFAULT '',
    event_type audit_log_event_type NOT NULL DEFAULT 'other',
    changes JSONB NOT NULL,
    belongs_to_household TEXT REFERENCES households("id") ON DELETE CASCADE,
    belongs_to_user TEXT REFERENCES users("id") ON DELETE CASCADE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW() NOT NULL
);
