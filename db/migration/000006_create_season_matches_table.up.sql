CREATE TABLE IF NOT EXISTS season_matches (
    id VARCHAR PRIMARY KEY,
    season_id VARCHAR NOT NULL REFERENCES seasons(id),
    favorite_contender_id VARCHAR NOT NULL REFERENCES contenders(id),
    underdog_contender_id VARCHAR NOT NULL REFERENCES contenders(id),
    week INTEGER NOT NULL,
    position INTEGER NOT NULL,
    weight_class VARCHAR NOT NULL,
    gender VARCHAR NOT NULL,
    is_main_event BOOLEAN NOT NULL DEFAULT FALSE,
    is_co_main_event BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ DEFAULT NULL,
    deleted_at TIMESTAMPTZ DEFAULT NULL
);

CREATE INDEX IF NOT EXISTS idx_season_matches_season_id ON season_matches (season_id);
CREATE INDEX IF NOT EXISTS idx_season_matches_favorite_contender_id ON season_matches (favorite_contender_id);
CREATE INDEX IF NOT EXISTS idx_season_matches_underdog_contender_id ON season_matches (underdog_contender_id);
