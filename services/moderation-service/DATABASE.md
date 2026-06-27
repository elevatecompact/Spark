# moderation-service — Database Schema
## PostgreSQL — Rules, actions, queue
### moderation_rules: id UUID PK, name, category(harassment,spam,nsfw,violence,hate_speech), severity(warn,restrict,remove,suspend), conditions JSONB (regex, ML model thresholds), is_active, priority INT, created_at
### moderation_actions: id UUID PK, user_id FK, content_id, rule_id FK, action_type, severity, status(pending,applied,appealed,reversed), applied_by(automated|moderator_id), applied_at
### eview_queue: id UUID PK, content_type, content_id, flagged_by(automated|report), reasons TEXT[], assigned_moderator_id UUID nullable, status(pending,reviewing,resolved), resolution, resolved_at
### content_reports: id UUID PK, reporter_id FK, content_type, content_id, reason ENUM, description TEXT, status(open,investigating,resolved), created_at
## Redis — Rate limit counters for reports, temporary action cache, real-time filter rules cache
## S3 — Evidence storage (flagged images, chat logs for appeal review)
