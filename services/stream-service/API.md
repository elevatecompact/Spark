# stream-service — API Contract
## Stream Management
POST /v1/streams (create), GET /v1/streams/{id}, PATCH /v1/streams/{id}, DELETE /v1/streams/{id} (archive), POST /v1/streams/{id}/start, POST /v1/streams/{id}/stop

## Ingest
GET /v1/ingest/endpoints (nearest by geo), POST /v1/ingest/credentials (generate RTMP key)

## Playback
GET /v1/playback/{id}/manifest.m3u8 (HLS), GET /v1/playback/{id}/manifest.mpd (DASH), GET /v1/playback/{id}/status

## Transcoding
GET /v1/transcoding/profiles (720p, 1080p, source available)

## Recording
POST /v1/streams/{id}/recording/start, POST /v1/streams/{id}/recording/stop
