# discovery-service — Troubleshooting
## Feed showing stale content: Cache not invalidated, trending refresh lagging, editorial picks expired. Flush Redis feed cache, verify trending_refresh cron, check editorial_pick end_at dates.
## Category page empty: Category inactive, no content assigned, index not rebuilt. Check categories.is_active, verify content has category assignment, rebuild category content index.
## Trending skewed: Velocity formula imbalance, bot activity inflating metrics, window too short. Review trending score components, check for anomalous viewer spikes, adjust TRENDING_VELOCITY_WINDOW_MINUTES.
## Editorial picks not appearing: Date range not active (check start_at/end_at), cache stale, pick_type mismatch. Verify current date within pick range, flush cache, confirm feed query includes editorial source.
