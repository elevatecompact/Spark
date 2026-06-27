# stream-service — Configuration
STREAM_PORT=4004, STREAM_DB_URL, STREAM_REDIS_URL, INGEST_REGIONS=us-east,eu-west,ap-southeast, TRANSCODING_PROFILES=720p,1080p,source, MAX_STREAM_DURATION_MINUTES=1440, RECORDING_BUCKET=titan-stream-recordings, HLS_SEGMENT_DURATION=4, RTMP_PORT=1935
FF: subscriber_only_streams=true, auto_record=false, transcoding_enabled=true, low_latency_mode=true, thumbnail_capture=true
Rate limits: 5 streams/h per creator, 2 concurrent (free), 10 (pro), 20 ingest keys/day
