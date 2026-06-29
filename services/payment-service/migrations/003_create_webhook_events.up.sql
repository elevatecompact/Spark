CREATE TABLE webhook_events (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    processor VARCHAR(20) NOT NULL CHECK (processor IN ('stripe', 'paypal')),
    external_event_id VARCHAR(255) NOT NULL,
    type VARCHAR(255) NOT NULL DEFAULT '',
    body JSONB NOT NULL DEFAULT '{}',
    status VARCHAR(20) NOT NULL DEFAULT 'received' CHECK (status IN ('received', 'processed', 'failed')),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    processed_at TIMESTAMPTZ
);

CREATE UNIQUE INDEX idx_webhook_external ON webhook_events(processor, external_event_id);
CREATE INDEX idx_webhook_status ON webhook_events(status);
