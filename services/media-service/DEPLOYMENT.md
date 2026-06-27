# media-service — Deployment Guide
Components: API server (3x 2GB), upload-worker (4x 2GB for post-upload processing), transcode-worker (GPU, 4x 16GB for video transcoding), thumbnail-worker (2x 1GB).
K8s: k8s/media-service/ — api, upload, transcode (GPU nodeSelector), thumbnail deployments.
Transcoding: FFmpeg with NVIDIA NVENC acceleration. HLS packaging with 4s segments. Thumbnail generation with FFmpeg thumbnail filter.
Deploy: kubectl apply -f k8s/media-service/. S3 bucket policies and CORS configured separately via Terraform.
Health: /health (DB+Redis+Kafka+S3), /ready (transcoder pool available), /metrics :4122.
