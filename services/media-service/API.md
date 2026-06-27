# media-service — API Contract
## Upload: POST /v1/upload/init (initiate resumable upload), POST /v1/upload/{id}/chunk (upload chunk), POST /v1/upload/{id}/complete (finalize), GET /v1/upload/{id}/status, DELETE /v1/upload/{id} (cancel)
## Processing: POST /v1/media/transcode (video → HLS/DASH), POST /v1/media/thumbnail (generate thumbnails), POST /v1/media/optimize (image resize/compress), GET /v1/media/{id}/status (processing status)
## Delivery: GET /v1/media/{id}/playback (HLS manifest URL), GET /v1/media/{id}/thumbnail/{time} (thumbnail at time), GET /v1/media/{id}/download (direct download), GET /v1/media/{id}/info (metadata)
## DRM: POST /v1/drm/license (Widevine/FairPlay license), POST /v1/drm/policy (create DRM policy — admin), GET /v1/drm/policies
## Admin: GET /v1/admin/storage/usage, POST /v1/admin/cache/purge {path}, GET /v1/admin/processing/queue
