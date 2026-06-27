# licensing-service — Database Schema
## licenses: id UUID PK, rights_holder_id FK, licensee_id FK, content_id FK nullable, type(exclusive,non_exclusive,sync,performance), scope(platform,territory,global), territory TEXT[], start_date DATE, end_date DATE, auto_renew BOOLEAN, rate_type(flat,fixed_per_use,revenue_share), rate_cents INT, revenue_share_percent DECIMAL(5,2), min_guarantee_cents, status(draft,pending,active,expired,terminated), terms_url (S3), created_at
## content_rights: id UUID PK, content_id FK, rights_holder_id FK, license_id FK, restrictions JSONB {geo_block:[], platforms:[], exclusivity_end}, registered_at
## usage_log: id UUID PK, license_id FK, content_id FK, usage_type(stream,download,performance,sync), context(video_id,stream_id,commerce_order_id), metadata JSONB, recorded_at. Append-only.
## oyalty_statements: id UUID PK, license_id FK, rights_holder_id FK, period_start, period_end, usage_count, rate_applied, total_cents, status(pending,paid,disputed), paid_at, created_at
## Redis: License validity cache (TTL 5min), content rights cache (TTL 1min), usage rate limit counters
