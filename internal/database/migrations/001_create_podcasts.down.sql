-- Drop indexes
DROP INDEX IF EXISTS idx_podcasts_language;
DROP INDEX IF EXISTS idx_podcasts_level;
DROP INDEX IF EXISTS idx_podcasts_country;
DROP INDEX IF EXISTS idx_podcasts_topic;
DROP INDEX IF EXISTS idx_podcasts_created_at;

-- Drop table
DROP TABLE IF EXISTS podcasts;
