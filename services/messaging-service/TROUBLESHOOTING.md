# messaging-service — Troubleshooting
## Messages not delivering: Recipient WS disconnected, message stuck in outbound queue, membership revoked. Check delivery status in DB.
## Read receipts not updating: last_read_message_id stale, Redis unread counter out of sync. Recalculate: ./messaging recalculate-unreads.
## Attachment fails: File too large, S3 misconfigured, presigned URL expired. Verify ATTACHMENT_MAX_SIZE_MB, S3 policy, regenerate URL.
