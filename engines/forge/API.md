# API

## Transcoding Endpoints
- POST /v1/job - Submit transcoding job with inputUri, outputUri, profile, ladder.
- GET /v1/job/:jobId - Query job status and progress.
- POST /v1/job/:jobId/cancel - Cancel a running job.
- POST /v1/job/:jobId/pause - Pause encoding, keep state.

## Profile Endpoints
- GET /v1/profiles - List available encoding profiles.
- POST /v1/profiles - Create custom encoding profile.

## Health
- GET /v1/health - GPU utilization, queue depth, node status.
- GET /v1/gpu - Per-GPU metrics: memory, encoder sessions, temperature.
