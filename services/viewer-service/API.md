# viewer-service — API Contract

## Watch History
- POST /v1/history — Record watch start/progress/completion event
- GET /v1/history — Paginated watch history (filter by type, date range)
- DELETE /v1/history/{id} — Remove single history entry
- DELETE /v1/history — Clear all history (confirmation required)

## Preferences
- GET /v1/preferences — Full preference object (categories, language, autoplay, mature content)
- PATCH /v1/preferences — Partial update of preference fields
- PUT /v1/preferences — Full replace of preference object

## Bookmarks
- POST /v1/bookmarks — Bookmark content with optional note
- GET /v1/bookmarks — List with folder/collection grouping
- DELETE /v1/bookmarks/{id} — Remove bookmark
- POST /v1/watch-later — Add to watch later queue (ordered)
- GET /v1/watch-later — List queue with reorder support
- DELETE /v1/watch-later/{id} — Remove from queue

## Engagement
- POST /v1/ratings — Rate content 1-5 stars (idempotent per user+content)
- POST /v1/reactions — Like/dislike toggle
- POST /v1/reports — Report content (type: spam, harassment, copyright)

All list endpoints support cursor-based pagination with ?cursor=&limit=100.
