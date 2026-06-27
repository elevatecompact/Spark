# chat-service — Runbook
## Alerts: ChatDeliveryLatency > 500ms, WSConnectionErrors > 5%, ChatModerationLatency > 200ms, MessageBacklog > 10000
## Drain WS node: POST /v1/admin/drain then kubectl delete pod
## Clear rate limit: redis-cli DEL "ratelimit:chat:{roomId}:{userId}"
## Capacity: Each WS node handles ~5000 concurrent connections. Scale at 80%.
