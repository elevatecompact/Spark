CREATE TABLE IF NOT EXISTS communities (
    id UUID PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    description TEXT NOT NULL DEFAULT '',
    creator_id UUID NOT NULL,
    type VARCHAR(16) NOT NULL DEFAULT 'public',
    category VARCHAR(50) NOT NULL DEFAULT '',
    avatar_url TEXT NOT NULL DEFAULT '',
    banner_url TEXT NOT NULL DEFAULT '',
    rules TEXT[] NOT NULL DEFAULT '{}',
    member_count INT NOT NULL DEFAULT 0,
    post_count INT NOT NULL DEFAULT 0,
    is_active BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS community_members (
    community_id UUID NOT NULL REFERENCES communities(id),
    user_id UUID NOT NULL,
    role VARCHAR(16) NOT NULL DEFAULT 'member',
    joined_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    last_active_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY (community_id, user_id)
);

CREATE TABLE IF NOT EXISTS community_posts (
    id UUID PRIMARY KEY,
    community_id UUID NOT NULL REFERENCES communities(id),
    author_id UUID NOT NULL,
    title VARCHAR(256) NOT NULL,
    content TEXT NOT NULL,
    is_pinned BOOLEAN NOT NULL DEFAULT false,
    is_announcement BOOLEAN NOT NULL DEFAULT false,
    reaction_counts JSONB NOT NULL DEFAULT '{}'::jsonb,
    comment_count INT NOT NULL DEFAULT 0,
    deleted_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_community_posts_community ON community_posts (community_id, deleted_at);

CREATE TABLE IF NOT EXISTS post_comments (
    id UUID PRIMARY KEY,
    post_id UUID NOT NULL REFERENCES community_posts(id),
    author_id UUID NOT NULL,
    parent_id UUID REFERENCES post_comments(id),
    content TEXT NOT NULL,
    reaction_counts JSONB NOT NULL DEFAULT '{}'::jsonb,
    deleted_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS post_reactions (
    post_id UUID NOT NULL REFERENCES community_posts(id),
    comment_id UUID REFERENCES post_comments(id),
    user_id UUID NOT NULL,
    emoji VARCHAR(32) NOT NULL,
    PRIMARY KEY (post_id, user_id, emoji)
);
