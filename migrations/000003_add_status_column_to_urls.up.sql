ALTER TABLE urls ADD COLUMN status VARCHAR(20) DEFAULT 'active';

CREATE UNIQUE INDEX urls_active_long_url_unique 
ON urls(long_url) WHERE status = 'active';

CREATE INDEX idx_urls_status_expires ON urls(status, expires_at);
