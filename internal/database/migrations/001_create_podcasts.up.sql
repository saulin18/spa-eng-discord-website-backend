-- Create podcasts table
CREATE TABLE IF NOT EXISTS podcasts (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    title VARCHAR(255) NOT NULL,
    description TEXT NOT NULL,
    image_url TEXT NOT NULL,
    language VARCHAR(10) NOT NULL CHECK (language IN ('en', 'es', 'both')),
    level VARCHAR(20) NOT NULL CHECK (level IN ('beginner', 'intermediate', 'advanced')),
    country VARCHAR(100) NOT NULL,
    topic VARCHAR(100) NOT NULL,
    url TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Create indexes for common filters
CREATE INDEX idx_podcasts_language ON podcasts(language);
CREATE INDEX idx_podcasts_level ON podcasts(level);
CREATE INDEX idx_podcasts_country ON podcasts(country);
CREATE INDEX idx_podcasts_topic ON podcasts(topic);
CREATE INDEX idx_podcasts_created_at ON podcasts(created_at DESC);

-- Seed some initial data
INSERT INTO podcasts (title, description, image_url, language, level, country, topic, url) VALUES
(
    'SpanishPod101',
    'Learn Spanish with fun and effective lessons. Perfect for beginners to advanced learners with native speaker hosts.',
    'https://images.unsplash.com/photo-1478737270239-2f02b77fc618?w=200&h=200&fit=crop',
    'es',
    'beginner',
    'Mexico',
    'General Learning',
    'https://www.spanishpod101.com/'
),
(
    'Notes in Spanish',
    'Real Spanish conversation between a native speaker and an English learner. Great for intermediate listeners.',
    'https://images.unsplash.com/photo-1543946207-39bd91e70ca7?w=200&h=200&fit=crop',
    'both',
    'intermediate',
    'Spain',
    'Conversation',
    'https://www.notesinspanish.com/'
),
(
    'Coffee Break Spanish',
    'Learn Spanish in coffee-break sized lessons. Structured courses from beginner to advanced.',
    'https://images.unsplash.com/photo-1509042239860-f550ce710b93?w=200&h=200&fit=crop',
    'es',
    'beginner',
    'Spain',
    'General Learning',
    'https://coffeebreaklanguages.com/coffeebreakspanish/'
),
(
    'Españolistos',
    'Colombian Spanish lessons covering grammar, vocabulary, and culture with a native Colombian teacher.',
    'https://images.unsplash.com/photo-1516483638261-f4dbaf036963?w=200&h=200&fit=crop',
    'both',
    'intermediate',
    'Colombia',
    'Grammar & Culture',
    'https://espanolistos.com/'
),
(
    'Duolingo Spanish Podcast',
    'True stories in easy-to-understand Spanish with English narration. Great for intermediate learners.',
    'https://images.unsplash.com/photo-1434030216411-0b793f4b4173?w=200&h=200&fit=crop',
    'both',
    'intermediate',
    'Various',
    'Stories',
    'https://podcast.duolingo.com/spanish'
),
(
    'Radio Ambulante',
    'Award-winning storytelling podcast in Spanish covering Latin American stories and culture.',
    'https://images.unsplash.com/photo-1598488035139-bdbb2231ce04?w=200&h=200&fit=crop',
    'es',
    'advanced',
    'Various',
    'Stories & Culture',
    'https://radioambulante.org/'
);
