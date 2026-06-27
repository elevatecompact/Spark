# media-service — Testing Guide
## Unit: Upload chunk assembly (checksum verification), transcoding profile resolution (codec selection), thumbnail time extraction, DRM license validation, CDN URL generation.
## Integration: Resumable upload (chunk with network interruption), transcode→rendition→CDN pipeline, thumbnail generation at multiple time points, image optimization (resize + webp conversion).
## Load: 500 concurrent chunked uploads (5MB chunks), 50 simultaneous transcode jobs, 10000 CDN URL requests/s.
## Tools: k6 for upload load test. FFmpeg mocking for unit tests.
