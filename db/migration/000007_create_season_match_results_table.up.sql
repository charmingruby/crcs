CREATE TABLE IF NOT EXISTS season_match_results (
    id VARCHAR PRIMARY KEY,
    match_id VARCHAR NOT NULL UNIQUE REFERENCES season_matches(id),
    winner_id VARCHAR DEFAULT NULL REFERENCES contenders(id),
    method VARCHAR NOT NULL,
    fight_duration INTEGER NOT NULL,
    scheduled_rounds INTEGER NOT NULL,
    actual_rounds INTEGER NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ DEFAULT NULL,
    deleted_at TIMESTAMPTZ DEFAULT NULL
);

CREATE INDEX IF NOT EXISTS idx_season_match_results_winner_id ON season_match_results (winner_id);
