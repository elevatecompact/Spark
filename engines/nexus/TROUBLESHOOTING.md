# Troubleshooting

## Upload failures
1. Check file size under s3.max_upload_size.
2. Verify S3 reachable: ./nexus test-s3.
3. Check mime type validation.
4. Review S3 bucket policy allows writes.

## Transformation stuck
1. Check worker process running.
2. Verify Redis queue accessible.
3. Check FFmpeg/ImageMagick installed.
4. Ensure disk space for temp files (2x input).
