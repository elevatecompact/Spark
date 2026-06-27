CREATE TABLE IF NOT EXISTS portfolio_items (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    creator_id UUID NOT NULL REFERENCES creators(id) ON DELETE CASCADE,
    title VARCHAR(200) NOT NULL,
    description TEXT DEFAULT '',
    media_url TEXT NOT NULL,
    media_type VARCHAR(10) NOT NULL CHECK (media_type IN ('video', 'image', 'audio')),
    thumbnail_url TEXT DEFAULT '',
    featured BOOLEAN DEFAULT false,
    sort_order INTEGER DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_portfolio_creator_id ON portfolio_items(creator_id);
CREATE INDEX IF NOT EXISTS idx_portfolio_featured ON portfolio_items(creator_id, featured) WHERE featured = true;
CREATE INDEX IF NOT EXISTS idx_portfolio_sort_order ON portfolio_items(creator_id, sort_order ASC);
