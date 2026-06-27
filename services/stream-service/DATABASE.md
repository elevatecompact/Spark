# stream-service — Database Schema
## stream_sessions: id UUID PK, creator_id UUID FK, title, status ENUM(created,starting,live,ending,ended,archived), ingest_protocol(rtmp,whip,srt), ingest_endpoint, is_subscriber_only BOOLEAN, started_at, ended_at
## 	ranscoding_jobs: id UUID PK, stream_id FK, profiles TEXT[], status ENUM(pending,processing,completed,failed), output_manifest TEXT (S3 URL)
## Redis: stream metadata (live only), viewer count per stream, ingest node assignments
## S3: Transcoding segments, manifests, recording VOD files, thumbnail captures
