CREATE TABLE IF NOT EXISTS contenders (
    id VARCHAR PRIMARY KEY,
    first_name VARCHAR NOT NULL,
    last_name VARCHAR NOT NULL,
    nickname VARCHAR NOT NULL,
    birthdate TIMESTAMPTZ NOT NULL,
    place_of_birth VARCHAR NOT NULL,
    current_weight_class VARCHAR NOT NULL,
    reach INTEGER NOT NULL,
    height INTEGER NOT NULL,
    weight INTEGER NOT NULL,
    gender VARCHAR NOT NULL,
    season_id VARCHAR NOT NULL REFERENCES seasons(id),
    status VARCHAR NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ DEFAULT NULL,
    deleted_at TIMESTAMPTZ DEFAULT NULL
);

CREATE INDEX IF NOT EXISTS idx_contenders_season_id ON contenders (season_id);

CREATE UNIQUE INDEX IF NOT EXISTS uidx_contenders_name_season
    ON contenders (first_name, last_name, season_id)
    WHERE deleted_at IS NULL;
