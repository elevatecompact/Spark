CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS watch_history (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    viewer_id UUID NOT NULL,
    content_id UUID NOT NULL,
    content_type VARCHAR(20) NOT NULL,
    progress FLOAT DEFAULT 0.0,
    watch_duration_seconds INTEGER DEFAULT 0,
    completed BOOLEAN DEFAULT FALSE,
    watched_at TIMESTAMPTZ DEFAULT NOW(),
    created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_watch_history_viewer_id ON watch_history(viewer_id);
CREATE INDEX IF NOT EXISTS idx_watch_history_content_id ON watch_history(content_id);
CREATE INDEX IF NOT EXISTS idx_watch_history_viewer_content ON watch_history(viewer_id, content_id);
CREATE INDEX IF NOT EXISTS idx_watch_history_watched_at ON watch_history(watched_at DESC);

CREATE TABLE IF NOT EXISTS viewer_preferences (
    viewer_id UUID PRIMARY KEY,
    preferred_categories UUID[] DEFAULT '{}',
    content_language VARCHAR(10) DEFAULT 'en',
    autoplay BOOLEAN DEFAULT TRUE,
    mature_content_allowed BOOLEAN DEFAULT FALSE,
    notification_prefs JSONB DEFAULT '{}',
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS bookmarks (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    viewer_id UUID NOT NULL,
    content_id UUID NOT NULL,
    note TEXT DEFAULT '',
    folder VARCHAR(100) DEFAULT '',
    created_at TIMESTAMPTZ DEFAULT NOW(),
    UNIQUE(viewer_id, content_id)
);

CREATE INDEX IF NOT EXISTS idx_bookmarks_viewer_id ON bookmarks(viewer_id);
CREATE INDEX IF NOT EXISTS idx_bookmarks_folder ON bookmarks(folder);

CREATE TABLE IF NOT EXISTS watch_later (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    viewer_id UUID NOT NULL,
    content_id UUID NOT NULL,
    position INTEGER DEFAULT 0,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    UNIQUE(viewer_id, content_id)
);

CREATE INDEX IF NOT EXISTS idx_watch_later_viewer_id ON watch_later(viewer_id);

CREATE TABLE IF NOT EXISTS ratings (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    viewer_id UUID NOT NULL,
    content_id UUID NOT NULL,
    score INTEGER NOT NULL CHECK (score >= 1 AND score <= 5),
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    UNIQUE(viewer_id, content_id)
);

CREATE INDEX IF NOT EXISTS idx_ratings_content_id ON ratings(content_id);

CREATE TABLE IF NOT EXISTS reactions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    viewer_id UUID NOT NULL,
    content_id UUID NOT NULL,
    type VARCHAR(10) NOT NULL CHECK (type IN ('like', 'dislike')),
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    UNIQUE(viewer_id, content_id)
);

CREATE INDEX IF NOT EXISTS idx_reactions_content_id ON reactions(content_id);

CREATE TABLE IF NOT EXISTS reports (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    viewer_id UUID NOT NULL,
    content_id UUID NOT NULL,
    type VARCHAR(20) NOT NULL CHECK (type IN ('spam', 'harassment', 'copyright', 'other')),
    description TEXT DEFAULT '',
    created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_reports_content_id ON reports(content_id);

CREATE OR REPLACE FUNCTION update_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

DROP TRIGGER IF EXISTS viewer_preferences_updated_at ON viewer_preferences;
CREATE TRIGGER viewer_preferences_updated_at
    BEFORE UPDATE ON viewer_preferences
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at();

DROP TRIGGER IF EXISTS ratings_updated_at ON ratings;
CREATE TRIGGER ratings_updated_at
    BEFORE UPDATE ON ratings
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at();

DROP TRIGGER IF EXISTS reactions_updated_at ON reactions;
CREATE TRIGGER reactions_updated_at
    BEFORE UPDATE ON reactions
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at();
