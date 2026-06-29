CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE gift_items (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(255) NOT NULL,
    price_cents BIGINT NOT NULL CHECK (price_cents > 0),
    image_url TEXT NOT NULL DEFAULT '',
    category VARCHAR(20) NOT NULL CHECK (category IN ('emote', 'badge', 'effect', 'sub')),
    is_active BOOLEAN NOT NULL DEFAULT true,
    sort_order INT NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_gift_items_active ON gift_items(is_active);
CREATE INDEX idx_gift_items_category ON gift_items(category);
