# media-service — Runbook
## Alerts: UploadFailureRate > 3%, TranscodingQueueDepth > 200, TranscodeFailureRate > 5%, CDNPurgeDelay > 60s, StorageUsage > 80% of S3 bucket, DRMLicenseFailureRate > 2%
## Purge CDN: POST /v1/admin/cache/purge {path: "/media/*"} — use sparingly, full purge takes 15min.
## Retry transcoding: POST /v1/admin/media/{id}/retry — re-queues failed transcode job.
## Cancel upload: DELETE /v1/admin/uploads/{id} — removes incomplete upload and cleans S3 chunks.
## Check storage: GET /v1/admin/storage/usage — per-bucket breakdown. Archive old source files if > 80%.
