# chat-service — Database Schema
## chat_rooms: id UUID PK, name, type(stream,channel,dm,system), owner_id FK, slow_mode_seconds INT, is_active BOOLEAN
## chat_messages: id UUID PK, room_id FK, user_id FK, content TEXT, content_type(text,emote,media,system), moderation_status(pending,approved,rejected), edited_at, deleted_at (soft), created_at. Partitioned by month, retained 90 days.
## Redis: Active room members (SET), rate limit counters, slow mode state, user mute/bans (cached)
## ClickHouse: Message volume/min, active chatters, emote usage stats
