CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS chat_rooms (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(255) NOT NULL,
    type VARCHAR(50) NOT NULL DEFAULT 'public',
    owner_id UUID NOT NULL,
    slow_mode_seconds INTEGER NOT NULL DEFAULT 0,
    is_active BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_chat_rooms_owner_id ON chat_rooms(owner_id);
CREATE INDEX idx_chat_rooms_type ON chat_rooms(type);
CREATE INDEX idx_chat_rooms_is_active ON chat_rooms(is_active);

CREATE TABLE IF NOT EXISTS chat_messages (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    room_id UUID NOT NULL REFERENCES chat_rooms(id) ON DELETE CASCADE,
    user_id UUID NOT NULL,
    username VARCHAR(255) NOT NULL,
    content TEXT NOT NULL,
    content_type VARCHAR(50) NOT NULL DEFAULT 'text',
    moderation_status VARCHAR(50) NOT NULL DEFAULT 'approved',
    emote_codes TEXT[] NOT NULL DEFAULT '{}',
    edited_at TIMESTAMPTZ,
    deleted_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
) PARTITION BY RANGE (created_at);

CREATE INDEX idx_chat_messages_room_id ON chat_messages(room_id);
CREATE INDEX idx_chat_messages_user_id ON chat_messages(user_id);
CREATE INDEX idx_chat_messages_created_at ON chat_messages(created_at);
CREATE INDEX idx_chat_messages_room_id_created_at ON chat_messages(room_id, created_at DESC);

CREATE TABLE IF NOT EXISTS chat_messages_2024_01 PARTITION OF chat_messages
    FOR VALUES FROM ('2024-01-01') TO ('2024-02-01');
CREATE TABLE IF NOT EXISTS chat_messages_2024_02 PARTITION OF chat_messages
    FOR VALUES FROM ('2024-02-01') TO ('2024-03-01');
CREATE TABLE IF NOT EXISTS chat_messages_2024_03 PARTITION OF chat_messages
    FOR VALUES FROM ('2024-03-01') TO ('2024-04-01');
CREATE TABLE IF NOT EXISTS chat_messages_2024_04 PARTITION OF chat_messages
    FOR VALUES FROM ('2024-04-01') TO ('2024-05-01');
CREATE TABLE IF NOT EXISTS chat_messages_2024_05 PARTITION OF chat_messages
    FOR VALUES FROM ('2024-05-01') TO ('2024-06-01');
CREATE TABLE IF NOT EXISTS chat_messages_2024_06 PARTITION OF chat_messages
    FOR VALUES FROM ('2024-06-01') TO ('2024-07-01');
CREATE TABLE IF NOT EXISTS chat_messages_2024_07 PARTITION OF chat_messages
    FOR VALUES FROM ('2024-07-01') TO ('2024-08-01');
CREATE TABLE IF NOT EXISTS chat_messages_2024_08 PARTITION OF chat_messages
    FOR VALUES FROM ('2024-08-01') TO ('2024-09-01');
CREATE TABLE IF NOT EXISTS chat_messages_2024_09 PARTITION OF chat_messages
    FOR VALUES FROM ('2024-09-01') TO ('2024-10-01');
CREATE TABLE IF NOT EXISTS chat_messages_2024_10 PARTITION OF chat_messages
    FOR VALUES FROM ('2024-10-01') TO ('2024-11-01');
CREATE TABLE IF NOT EXISTS chat_messages_2024_11 PARTITION OF chat_messages
    FOR VALUES FROM ('2024-11-01') TO ('2024-12-01');
CREATE TABLE IF NOT EXISTS chat_messages_2024_12 PARTITION OF chat_messages
    FOR VALUES FROM ('2024-12-01') TO ('2025-01-01');

CREATE TABLE IF NOT EXISTS chat_emotes (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    code VARCHAR(100) NOT NULL UNIQUE,
    image_url VARCHAR(500) NOT NULL,
    provider VARCHAR(100) NOT NULL DEFAULT 'system',
    is_global BOOLEAN NOT NULL DEFAULT true,
    room_id UUID REFERENCES chat_rooms(id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT unique_room_emote UNIQUE (code, room_id)
);

CREATE INDEX idx_chat_emotes_code ON chat_emotes(code);
CREATE INDEX idx_chat_emotes_room_id ON chat_emotes(room_id);
CREATE INDEX idx_chat_emotes_is_global ON chat_emotes(is_global);

CREATE OR REPLACE FUNCTION delete_expired_messages() RETURNS void AS $$
BEGIN
    DELETE FROM chat_messages
    WHERE created_at < NOW() - INTERVAL '90 days'
    AND deleted_at IS NULL;
END;
$$ LANGUAGE plpgsql;
