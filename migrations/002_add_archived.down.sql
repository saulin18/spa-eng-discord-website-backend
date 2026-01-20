-- Remove archived column
DROP INDEX IF EXISTS idx_podcasts_archived;
ALTER TABLE podcasts DROP COLUMN IF EXISTS archived;
