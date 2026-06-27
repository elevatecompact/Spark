# Events

## Published Events
- asset.uploaded - Payload: { assetId, bucket, key, size, mimeType, hash }.
- asset.deleted, asset.restored.
- asset.transformed - Payload: { jobId, transformationType, outputKey }.
- asset.tier.changed - Payload: { previousTier, newTier, reason }.
- upload.multipart.failed, transform.started/completed/failed.

## Subscribed Events
- storage.bucket.created, storage.retention.updated.
- media.transcode.requested - From Forge engine.
