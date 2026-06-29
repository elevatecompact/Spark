CREATE TABLE IF NOT EXISTS products (
    id UUID PRIMARY KEY,
    creator_id UUID NOT NULL,
    name VARCHAR(256) NOT NULL,
    description TEXT NOT NULL DEFAULT '',
    type VARCHAR(16) NOT NULL DEFAULT 'digital',
    price_cents BIGINT NOT NULL DEFAULT 0,
    currency VARCHAR(4) NOT NULL DEFAULT 'USD',
    category VARCHAR(64) NOT NULL DEFAULT '',
    tags JSONB DEFAULT '[]'::jsonb,
    media_urls JSONB DEFAULT '[]'::jsonb,
    inventory_count BIGINT NOT NULL DEFAULT -1,
    is_active BOOLEAN NOT NULL DEFAULT true,
    is_featured BOOLEAN NOT NULL DEFAULT false,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS product_variants (
    id UUID PRIMARY KEY,
    product_id UUID NOT NULL REFERENCES products(id),
    name VARCHAR(128) NOT NULL,
    price_cents BIGINT,
    inventory_count BIGINT NOT NULL DEFAULT -1,
    sort_order INT NOT NULL DEFAULT 0
);

CREATE TABLE IF NOT EXISTS cart_items (
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL,
    product_id UUID NOT NULL,
    variant_id UUID,
    quantity INT NOT NULL DEFAULT 1,
    added_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE (user_id, product_id, COALESCE(variant_id, '00000000-0000-0000-0000-000000000000'))
);

CREATE TABLE IF NOT EXISTS orders (
    id UUID PRIMARY KEY,
    buyer_id UUID NOT NULL,
    merchant_id UUID NOT NULL,
    status VARCHAR(16) NOT NULL DEFAULT 'pending',
    total_cents BIGINT NOT NULL DEFAULT 0,
    currency VARCHAR(4) NOT NULL DEFAULT 'USD',
    payment_intent_id UUID,
    shipping_address JSONB,
    placed_at TIMESTAMPTZ,
    fulfilled_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS order_items (
    id UUID PRIMARY KEY,
    order_id UUID NOT NULL REFERENCES orders(id),
    product_id UUID NOT NULL,
    variant_id UUID,
    quantity INT NOT NULL DEFAULT 1,
    unit_price_cents BIGINT NOT NULL DEFAULT 0,
    fulfillment_status VARCHAR(16) NOT NULL DEFAULT 'pending',
    download_url TEXT NOT NULL DEFAULT '',
    fulfilled_at TIMESTAMPTZ
);

CREATE TABLE IF NOT EXISTS reviews (
    id UUID PRIMARY KEY,
    product_id UUID NOT NULL REFERENCES products(id),
    user_id UUID NOT NULL,
    rating INT NOT NULL CHECK (rating >= 1 AND rating <= 5),
    title VARCHAR(256) NOT NULL DEFAULT '',
    body TEXT NOT NULL DEFAULT '',
    is_verified_purchase BOOLEAN NOT NULL DEFAULT false,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE (product_id, user_id)
);

CREATE TABLE IF NOT EXISTS payouts (
    id UUID PRIMARY KEY,
    merchant_id UUID NOT NULL,
    amount_cents BIGINT NOT NULL DEFAULT 0,
    currency VARCHAR(4) NOT NULL DEFAULT 'USD',
    status VARCHAR(16) NOT NULL DEFAULT 'pending',
    period_start TIMESTAMPTZ NOT NULL,
    period_end TIMESTAMPTZ NOT NULL,
    paid_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
