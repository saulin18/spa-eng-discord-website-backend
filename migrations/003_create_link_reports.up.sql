-- Create link_reports table
CREATE TABLE IF NOT EXISTS link_reports (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    podcast_id UUID NOT NULL REFERENCES podcasts(id) ON DELETE CASCADE,
    reporter_ip VARCHAR(45), -- IPv6 can be up to 45 chars
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Index for efficient lookup by podcast
CREATE INDEX idx_link_reports_podcast_id ON link_reports(podcast_id);

-- Index for deduplication by IP (one report per IP per podcast)
CREATE UNIQUE INDEX idx_link_reports_podcast_ip ON link_reports(podcast_id, reporter_ip) WHERE reporter_ip IS NOT NULL;
