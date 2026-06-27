# advertising-service — Database Schema
## campaigns: id UUID PK, advertiser_id FK, name, budget_cents, spent_cents, status(draft,active,paused,ended), start_at, end_at, targeting JSONB {categories[],creators[],geo[],demographics{}}, daily_budget_cents, bid_strategy(cpm, cpc)
## d_units: id UUID PK, campaign_id FK, type(preroll,midroll,display,sponsored), format(video,image,text), creative_url, destination_url, width, height, duration_seconds, status(pending,approved,rejected), approval_note
## impressions: id UUID PK, campaign_id FK, ad_unit_id FK, placement_id, user_id FK nullable, cost_micro_cents INT (CPM in micro), device_type, geo, served_at. Append-only, partitioned daily.
## clicks: id UUID PK, impression_id FK, clicked_at. Append-only.
## d_inventory: placement_id PK, content_type, available_from, available_to, floor_price_micro_cents, is_active
## Redis: Ad placement cache (TTL 30s), user targeting profile cache (TTL 1h), campaign budget counters (real-time deduction)
## ClickHouse: Impression/click analytics, fraud detection queries
