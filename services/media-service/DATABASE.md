# media-service — Database Schema
## PostgreSQL — Media metadata
### media_assets: id UUID PK, uploader_id FK, content_type(video,image,audio), source_filename, file_size_bytes, mime_type, status(uploading,processing,ready,failed,deleted), storage_path (S3 key), cdn_url, duration_seconds (video/audio), width, height, checksum SHA-256, created_at
### media_renditions: id UUID PK, media_id FK, profile(thumbnail,720p,1080p,source), format(hls,dash,mp4,webp,jpg), file_size_bytes, storage_path, cdn_url, status, created_at
### 	ranscoding_jobs: id UUID PK, media_id FK, profile TEXT[], status(pending,processing,completed,failed), worker_id, started_at, completed_at, error_message
### drm_policies: id UUID PK, name, content_id FK nullable, key_system(widevine,fairplay), license_duration_seconds, security_level, is_active
## S3 — Primary storage buckets: titan-media-uploads (source), titan-media-renditions (processed), titan-media-thumbnails
## Redis — Upload session state (chunk tracking, TTL 24h), processing queue status, CDN cache status
