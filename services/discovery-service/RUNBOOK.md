# discovery-service — Runbook
## Alerts: FeedLatencyP99 > 500ms, TrendingStale > 5min (refresh not running), FeedDiversityScore < 0.2 (too many same category), CategoryContentCountStale, EditorialPickNotShowing
## Warm feeds: POST /v1/admin/feeds/cache/warm {feedTypes: ["home","trending"]}
## Refresh trending: POST /v1/admin/trending/refresh — forces immediate recalculation.
## Override feed: POST /v1/admin/feeds/override {userId, feedContents[]} — for testing or special cases.
## Rebuild category tree: POST /v1/admin/categories/rebuild-tree
