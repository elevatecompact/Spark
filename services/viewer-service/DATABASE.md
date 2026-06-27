# viewer-service — Event Contracts

## Published Events
| Topic | Description |
|-------|-------------|
| iewer.watch.started | Viewer began watching content (includes contentId, type, timestamp) |
| iewer.watch.completed | Viewer reached 100% progress or ended stream |
| iewer.watch.progress | Periodic progress update (every 30s while watching) |
| iewer.rating.submitted | Star rating given (1-5) with content reference |
| iewer.reaction.added | Like or dislike recorded |

## Consumed Events
| Topic | Source | Handler |
|-------|--------|---------|
| ecommendation.explanation.ready | recommendation-service | Store rec explanation for display |
| stream.session.ended | stream-service | Record stream as watched if viewer was present |
| content.removed | media-service | Clean up bookmarks pointing to removed content |

## Schema (WatchProgressEvent)
`json
{"viewerId":"uuid","contentId":"uuid","contentType":"live|recorded","progress":0.75,"watchDurationSeconds":2700,"timestamp":"ISO8601"}
`
"@ | Set-Content (Join-Path C:\Users\Dell\Downloads\SPARK\services\viewer-service "EVENTS.md") -Encoding UTF8

@"
# viewer-service — Database Schema

## PostgreSQL — Primary Store
### watch_history table
Partitioned by watched_at monthly. Retained for 90 days.
| Column | Type | Notes |
|--------|------|-------|
| id | UUID | PK |
| viewer_id | UUID | FK, indexed |
| content_id | UUID | Polymorphic reference |
| content_type | VARCHAR(20) | live, recorded, clip |
| progress | FLOAT | 0.0 to 1.0 |
| watch_duration_seconds | INTEGER | Cumulative |
| completed | BOOLEAN | DEFAULT false |
| watched_at | TIMESTAMPTZ | Partition key |

### iewer_preferences table
| viewer_id | UUID | PK, FK |
| preferred_categories | UUID[] | Ordered by preference |
| content_language | VARCHAR(10) | ISO 639-1 |
| autoplay | BOOLEAN | DEFAULT true |
| mature_content_allowed | BOOLEAN | DEFAULT false |
| notification_prefs | JSONB | Per-channel notification settings |

### ookmarks and watch_later tables
Both have viewer_id, content_id, created_at, with UNIQUE constraints on (viewer_id, content_id).
