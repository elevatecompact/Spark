# moderation-service — API Contract
## Scan: POST /v1/scan/text (content → policy violations), POST /v1/scan/image (URL → NSFW/violence categories), POST /v1/scan/video (stream ID → real-time analysis), POST /v1/scan/batch (up to 50 items)
## Rules: GET /v1/rules (policy rules), POST /v1/rules (create rule — admin), PATCH /v1/rules/{id} (update — admin), DELETE /v1/rules/{id} (admin)
## Review Queue: GET /v1/queue (pending items), POST /v1/queue/{id}/approve, POST /v1/queue/{id}/reject, POST /v1/queue/{id}/escalate, GET /v1/queue/stats (queue depth by type)
## Actions: POST /v1/actions/warn {userId, reason}, POST /v1/actions/restrict {userId, duration}, POST /v1/actions/remove {contentId}, POST /v1/actions/suspend {userId, reason, duration}, POST /v1/actions/reverse {actionId}
## Reports: POST /v1/reports (user report), GET /v1/reports/{id}, GET /v1/reports (admin — filterable)
## Admin: GET /v1/admin/stats (scans, actions, queue), GET /v1/admin/accuracy (auto vs human decision match rate)
