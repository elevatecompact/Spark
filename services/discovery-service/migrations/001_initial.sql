CREATE TABLE IF NOT EXISTS categories (
    id UUID PRIMARY KEY,
    name VARCHAR(128) NOT NULL,
    slug VARCHAR(128) UNIQUE NOT NULL,
    description TEXT NOT NULL DEFAULT '',
    parent_id UUID REFERENCES categories(id),
    icon_url TEXT NOT NULL DEFAULT '',
    sort_order INT NOT NULL DEFAULT 0,
    is_active BOOLEAN NOT NULL DEFAULT true,
    content_count BIGINT NOT NULL DEFAULT 0
);

CREATE TABLE IF NOT EXISTS category_contents (
    category_id UUID NOT NULL REFERENCES categories(id),
    content_id UUID NOT NULL,
    added_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY (category_id, content_id)
);

CREATE TABLE IF NOT EXISTS collections (
    id UUID PRIMARY KEY,
    title VARCHAR(256) NOT NULL,
    description TEXT NOT NULL DEFAULT '',
    type VARCHAR(16) NOT NULL DEFAULT 'editorial',
    cover_image_url TEXT NOT NULL DEFAULT '',
    is_featured BOOLEAN NOT NULL DEFAULT false,
    start_at TIMESTAMPTZ,
    end_at TIMESTAMPTZ,
    curated_by VARCHAR(128) NOT NULL DEFAULT '',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS collection_items (
    collection_id UUID NOT NULL REFERENCES collections(id),
    content_id UUID NOT NULL,
    sort_order INT NOT NULL DEFAULT 0,
    added_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY (collection_id, content_id)
);

CREATE TABLE IF NOT EXISTS editorial_picks (
    content_id UUID NOT NULL PRIMARY KEY,
    pick_type VARCHAR(16) NOT NULL,
    label VARCHAR(256) NOT NULL DEFAULT '',
    reason TEXT NOT NULL DEFAULT '',
    picked_by VARCHAR(128) NOT NULL DEFAULT '',
    start_at TIMESTAMPTZ NOT NULL,
    end_at TIMESTAMPTZ NOT NULL,
    sort_order INT NOT NULL DEFAULT 0
);

CREATE TABLE IF NOT EXISTS home_feed (
    content_id UUID NOT NULL PRIMARY KEY,
    score DOUBLE PRECISION NOT NULL DEFAULT 0
);

CREATE TABLE IF NOT EXISTS new_contents (
    content_id UUID NOT NULL PRIMARY KEY,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS trending_scores (
    content_id UUID NOT NULL PRIMARY KEY,
    trending_score DOUBLE PRECISION NOT NULL DEFAULT 0,
    calculated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS trending_creators (
    creator_id UUID NOT NULL PRIMARY KEY,
    trending_score DOUBLE PRECISION NOT NULL DEFAULT 0
);

CREATE TABLE IF NOT EXISTS related_contents (
    content_id UUID NOT NULL,
    related_content_id UUID NOT NULL,
    score DOUBLE PRECISION NOT NULL DEFAULT 0,
    PRIMARY KEY (content_id, related_content_id)
);
