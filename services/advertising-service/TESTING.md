# advertising-service — Testing Guide
## Unit: Campaign budget calculation (daily vs total), bid evaluation (highest CPM wins), targeting evaluation (match user profile to campaign targeting), fraud detection heuristics (rapid impressions, same IP).
## Integration: Full campaign lifecycle (create→approve→activate→serve→impression→click→analytics), budget exhaustion handling, inventory availability checks.
## Load: Ad server benchmark — 10000 QPS with p99 < 50ms. Impression recording — 50000 events/s. ClickHouse query performance for analytics dashboard.
