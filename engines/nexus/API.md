# API

## Upload
- POST /v1/upload - Upload file, returns assetId.
- POST /v1/upload/direct - Direct-to-S3 via presigned URL.
- PUT /v1/upload/multipart - Initiate multipart upload.

## Asset
- GET /v1/asset/:assetId - Get metadata.
- GET /v1/asset/:assetId/download - Presigned download URL.
- DELETE /v1/asset/:assetId - Soft delete.
- POST /v1/asset/:assetId/restore - Restore from trash.

## Transformation
- POST /v1/transform/image - Resize, crop, format conversion.
- POST /v1/transform/video/thumbnail - Video thumbnail.
