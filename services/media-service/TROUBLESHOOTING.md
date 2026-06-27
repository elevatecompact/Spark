# media-service — Troubleshooting
## Upload fails mid-way: Chunk timeout, checksum mismatch, S3 upload error. Check chunk upload status (Redis), verify client sends correct checksum, check S3 connectivity. Client should resume from last confirmed chunk.
## Transcoding stuck: Worker pod OOM, FFmpeg crash, codec not supported. Check transcode worker logs, verify GPU memory, check source file codec compatibility. Retry with different profiles.
## Video not playing: DRM license failed, CDN not serving, HLS manifest corrupted. Check DRM key delivery, CDN cache status (cache miss → origin fetch), verify HLS manifest syntax.
## Thumbnail generation fails: FFmpeg filter error, timestamp out of range, video duration not available. Check thumbnail interval vs video duration, verify source file seekable.
