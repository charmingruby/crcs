CREATE TABLE IF NOT EXISTS matchmaking_suggestions (
    id VARCHAR PRIMARY KEY,
    season_id VARCHAR NOT NULL REFERENCES seasons(id),
    kind VARCHAR NOT NULL,
    payload JSONB NOT NULL,
    rationale VARCHAR NOT NULL,
    status VARCHAR NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ DEFAULT NULL,
    deleted_at TIMESTAMPTZ DEFAULT NULL
);

CREATE INDEX IF NOT EXISTS idx_matchmaking_suggestions_season_id ON matchmaking_suggestions (season_id);
