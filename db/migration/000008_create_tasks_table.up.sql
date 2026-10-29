CREATE TABLE tasks (
    id VARCHAR PRIMARY KEY,
    type VARCHAR NOT NULL,
    metadata JSONB DEFAULT NULL,
    status VARCHAR NOT NULL,
    total_records INTEGER NOT NULL,
    processed_records INTEGER NOT NULL,
    failed_records INTEGER NOT NULL,
    attempt_count INTEGER NOT NULL,
    error_message TEXT,
    started_at TIMESTAMPTZ NOT NULL,
    finished_at TIMESTAMPTZ DEFAULT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ DEFAULT NULL,
    deleted_at TIMESTAMPTZ DEFAULT NULL
);

CREATE INDEX idx_tasks_metadata_season_id ON tasks ((metadata->>'season_id'));
