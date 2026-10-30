CREATE TABLE outbox (
    id VARCHAR PRIMARY KEY,
    topic VARCHAR NOT NULL,
    key VARCHAR NOT NULL,
    payload JSONB NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    published_at TIMESTAMPTZ
);

CREATE INDEX idx_outbox_pending ON outbox (created_at, id)
    WHERE published_at IS NULL;

CREATE TABLE inbox (
    listener VARCHAR NOT NULL,
    message_id VARCHAR NOT NULL,
    processed_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    PRIMARY KEY (listener, message_id)
);
