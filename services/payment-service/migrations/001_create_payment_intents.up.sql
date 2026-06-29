CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE payment_intents (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL,
    external_id VARCHAR(255) NOT NULL DEFAULT '',
    processor VARCHAR(20) NOT NULL DEFAULT 'stripe' CHECK (processor IN ('stripe', 'paypal')),
    amount_cents BIGINT NOT NULL CHECK (amount_cents > 0),
    currency VARCHAR(10) NOT NULL DEFAULT 'USD',
    status VARCHAR(30) NOT NULL DEFAULT 'requires_payment_method' CHECK (status IN ('requires_payment_method', 'processing', 'succeeded', 'failed', 'canceled')),
    idempotency_key VARCHAR(255) NOT NULL DEFAULT '',
    metadata JSONB NOT NULL DEFAULT '{}',
    payment_method_id UUID,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE UNIQUE INDEX idx_intents_idempotency ON payment_intents(idempotency_key) WHERE idempotency_key != '';
CREATE INDEX idx_intents_user ON payment_intents(user_id);
CREATE INDEX idx_intents_external ON payment_intents(external_id);
CREATE INDEX idx_intents_status ON payment_intents(status);
