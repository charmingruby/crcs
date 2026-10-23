CREATE TABLE IF NOT EXISTS cartel_matches (
    id VARCHAR PRIMARY KEY,
    cartel_id VARCHAR NOT NULL REFERENCES cartels(id),
    event_name VARCHAR NOT NULL,
    event_date TIMESTAMPTZ NOT NULL,
    event_rating INTEGER NOT NULL,
    method VARCHAR NOT NULL,
    outcome VARCHAR NOT NULL,
    opponent_name VARCHAR NOT NULL,
    weight_class VARCHAR NOT NULL,
    fight_duration INTEGER NOT NULL,
    scheduled_rounds INTEGER NOT NULL,
    actual_rounds INTEGER NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ DEFAULT NULL,
    deleted_at TIMESTAMPTZ DEFAULT NULL
);

CREATE INDEX IF NOT EXISTS idx_cartel_matches_cartel_id ON cartel_matches (cartel_id);
