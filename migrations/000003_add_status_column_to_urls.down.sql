
DROP INDEX idx_urls_status_expires;

DROP INDEX urls_active_long_url_unique;

ALTER TABLE urls DROP COLUMN status;