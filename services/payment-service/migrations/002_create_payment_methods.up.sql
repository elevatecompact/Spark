CREATE TABLE payment_methods (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL,
    external_id VARCHAR(255) NOT NULL,
    processor VARCHAR(20) NOT NULL CHECK (processor IN ('stripe', 'paypal')),
    type VARCHAR(20) NOT NULL CHECK (type IN ('card', 'paypal', 'bank')),
    fingerprint VARCHAR(255) NOT NULL DEFAULT '',
    last4 VARCHAR(4) NOT NULL DEFAULT '',
    exp_month INT NOT NULL DEFAULT 0,
    exp_year INT NOT NULL DEFAULT 0,
    is_default BOOLEAN NOT NULL DEFAULT false,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_methods_user ON payment_methods(user_id);
CREATE INDEX idx_methods_fingerprint ON payment_methods(fingerprint);
