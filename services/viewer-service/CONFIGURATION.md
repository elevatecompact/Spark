# viewer-service — Configuration

## Environment Variables
VIEWER_PORT=4003, VIEWER_DB_URL, VIEWER_REDIS_URL, HISTORY_RETENTION_DAYS=90, WATCH_PROGRESS_INTERVAL=30, MAX_BOOKMARKS=5000, MAX_WATCH_LATER=1000

## Feature Flags
track_watch_progress=true, enable_reactions=true, autoplay_enabled=true, history_export_enabled=false

## Rate Limits
Watch events: 1/s per viewer (batch preferred). Ratings: 1 per content per viewer. Bookmarks: 100/h per viewer.
