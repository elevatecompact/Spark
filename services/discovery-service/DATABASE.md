# discovery-service — Database Schema
## PostgreSQL — Categories, collections, editorial
### categories: id UUID PK, name VARCHAR, slug VARCHAR UNIQUE, description, parent_id UUID nullable (self-referential), icon_url, sort_order, is_active, content_count INT (denormalized)
### collections: id UUID PK, title VARCHAR, description TEXT, type(editorial,holiday,theme), cover_image_url, is_featured, start_at, end_at nullable, curated_by VARCHAR, created_at
### collection_items: collection_id+content_id PK, sort_order INT, added_at
### editorial_picks: content_id PK, pick_type(staff_pick,spotlight,holiday), label VARCHAR, reason TEXT, picked_by VARCHAR, start_at, end_at, sort_order
## Redis — Trending feeds (sorted sets by velocity score, TTL 5min), home feed cache (TTL 2min for logged-in, 5min for logged-out), category content cache (TTL 1min)
## Trending score algorithm: (viewers_last_15min^2 + gift_rate * 100 + chat_rate * 10) / hours_since_stream_start
