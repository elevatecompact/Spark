# stream-service — Troubleshooting
## Stream not starting: Wrong RTMP endpoint/key, ingest node full, tier limit. Verify: SELECT ingest_key FROM stream_sessions.
## Buffering/high latency: Transcoding profile too aggressive, origin underprovisioned, CDN miss. Reduce profiles, scale origins, pre-warm CDN.
## Recording missing: auto_record disabled, transcoding failed, S3 error. Check recording_enabled on session, transcoding_jobs status, S3 permissions.
