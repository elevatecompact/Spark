CREATE TABLE IF NOT EXISTS creator_profiles (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID UNIQUE NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    category VARCHAR(100),
    languages TEXT[] DEFAULT '{}',
    is_mature BOOLEAN DEFAULT FALSE,
    social_links TEXT[] DEFAULT '{}',
    follower_count BIGINT DEFAULT 0,
    subscriber_count BIGINT DEFAULT 0,
    total_streams BIGINT DEFAULT 0,
    total_viewers BIGINT DEFAULT 0,
    total_hours_streamed BIGINT DEFAULT 0,
    total_earnings BIGINT DEFAULT 0,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_creator_profiles_user ON creator_profiles(user_id);
CREATE INDEX idx_creator_profiles_category ON creator_profiles(category);
