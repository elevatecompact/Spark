# messaging-service — Runbook
## Alerts: MsgDeliveryRate < 99.9%, MsgP99Latency > 1s, WSConnectionErrors > 3%, AttachmentUploadFailures > 2%
## Purge: DELETE /v1/admin/conversations/{id}/purge
## Recalculate unreads: ./messaging recalculate-unreads --userId {userId}
## Force device sync: POST /v1/admin/users/{id}/sync-devices
