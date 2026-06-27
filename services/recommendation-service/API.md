# recommendation-service — API Contract
## Feeds: GET /v1/feeds/home (personalized), GET /v1/feeds/trending (global), GET /v1/feeds/up-next/{contentId} (recommend next), GET /v1/feeds/similar/{contentId} (similar content), GET /v1/feeds/creator/{creatorId} (creator's content)
## Models: GET /v1/models/active (current model version), POST /v1/models/deploy (admin), GET /v1/models/metrics (offline eval metrics)
## Feedback: POST /v1/feedback/click (implicit signal), POST /v1/feedback/dismiss (negative signal), POST /v1/feedback/explain/{recId} (why recommended)
## Admin: POST /v1/admin/refresh-features, GET /v1/admin/feature-importance, POST /v1/admin/invalidate-cache
