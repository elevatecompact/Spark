# translation-service — API Contract
## Translate: POST /v1/translate (text, sourceLang, targetLang), POST /v1/translate/batch (up to 100 texts), POST /v1/translate/streaming (WebSocket for real-time), GET /v1/translate/languages (supported languages with coverage %)
## Detect: POST /v1/detect (text → language + confidence), POST /v1/detect/batch
## Translation Memory: GET /v1/tm/lookup (exact/fuzzy match), POST /v1/tm/store (save translation), DELETE /v1/tm/entries/{id}
## Review: GET /v1/review/queue (pending review), POST /v1/review/{id}/approve, POST /v1/review/{id}/reject, POST /v1/review/{id}/correct
## Admin: GET /v1/admin/usage (API call counts by provider), POST /v1/admin/provider/switch, GET /v1/admin/coverage (language coverage metrics)
