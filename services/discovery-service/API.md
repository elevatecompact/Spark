# discovery-service — API Contract
## Feeds: GET /v1/feeds/home (personalized), GET /v1/feeds/trending (timeframe: today,week,month), GET /v1/feeds/category/{slug} (category feed), GET /v1/feeds/new (recently created), GET /v1/feeds/related/{contentId} (because you watched)
## Categories: GET /v1/categories (tree structure), GET /v1/categories/{slug} (category with subcategories), GET /v1/categories/{slug}/contents (paginated)
## Collections: GET /v1/collections, GET /v1/collections/{id} (curated list), POST /v1/collections (admin), PATCH /v1/collections/{id} (admin), POST /v1/collections/{id}/items (admin), DELETE /v1/collections/{id}/items/{contentId} (admin)
## Trending: GET /v1/trending (global), GET /v1/trending/category/{slug}, GET /v1/trending/creators (trending creators)
## Editorial: GET /v1/editorial/picks (staff picks), GET /v1/editorial/spotlight (hero content), GET /v1/editorial/holiday/{campaign}
## Admin: POST /v1/admin/feeds/cache/warm, POST /v1/admin/trending/refresh, POST /v1/admin/categories/reorder
