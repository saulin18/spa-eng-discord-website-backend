-- Add archived column to podcasts
ALTER TABLE podcasts ADD COLUMN archived BOOLEAN NOT NULL DEFAULT FALSE;

-- Index for faster filtering
CREATE INDEX idx_podcasts_archived ON podcasts(archived);
