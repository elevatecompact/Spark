CREATE TABLE IF NOT EXISTS passkeys (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    credential_id VARCHAR(500) NOT NULL UNIQUE,
    public_key BYTEA NOT NULL,
    attestation_type VARCHAR(50),
    transports TEXT[] DEFAULT '{}',
    aaguid UUID,
    sign_count BIGINT DEFAULT 0,
    name VARCHAR(100),
    device_type VARCHAR(50),
    backed_up BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    last_used_at TIMESTAMPTZ
);

CREATE INDEX IF NOT EXISTS idx_passkeys_user_id ON passkeys(user_id);
CREATE INDEX IF NOT EXISTS idx_passkeys_credential_id ON passkeys(credential_id);
