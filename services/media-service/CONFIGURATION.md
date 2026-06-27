# media-service — Configuration
MEDIA_PORT=4024, MEDIA_DB_URL, MEDIA_REDIS_URL, MEDIA_KAFKA_BROKERS, AWS_S3_REGION=us-east-1, UPLOAD_BUCKET=titan-media-uploads, RENDITION_BUCKET=titan-media-renditions, CDN_URL=https://cdn.titan.dev, MAX_UPLOAD_SIZE_BYTES=5368709120 (5GB), TRANSCODING_PROFILES=720p,1080p,source, THUMBNAIL_INTERVAL_SECONDS=30, DRM_ENABLED=false, UPLOAD_CHUNK_SIZE_BYTES=5242880 (5MB), UPLOAD_URL_TTL_MINUTES=60
FF: resumable_uploads=true, auto_transcode=true, thumbnail_generation=true, drm_enabled=false, cdn_purge_on_delete=true, image_optimization=true
Rate limits: 10 concurrent uploads per user, 100MB/min upload rate per user, 10 transcode requests/min per user
