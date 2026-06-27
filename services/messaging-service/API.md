# messaging-service — API Contract
## Conversations: POST /v1/conversations, GET /v1/conversations, GET /v1/conversations/{id}, DELETE /v1/conversations/{id}
## Messages: POST /v1/conversations/{id}/messages, GET (cursor paginated), PUT (edit within 1h), DELETE (soft delete), POST /v1/conversations/{id}/messages/{msgId}/reactions
## Read State: POST /v1/conversations/{id}/read (mark up to message), GET /v1/conversations/{id}/read-status
## Attachments: POST /v1/conversations/{id}/attachments, GET /v1/attachments/{id} (CDN redirect)
## Groups: POST/DELETE members, PATCH conversation (name/icon)
