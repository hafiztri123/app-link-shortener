CREATE TABLE archived_urls (
    id BIGINT NOT NULL,
    short_code VARCHAR(20),
    long_url TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL,
    expires_at TIMESTAMPTZ,
    archived_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY (id)
)
