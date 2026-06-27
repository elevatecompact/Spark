# API

## Moderation Endpoints
- POST /v1/moderate/text - Moderate text content. Returns decision: allow/flag/block.
- POST /v1/moderate/image - Moderate image (base64 or URL).
- POST /v1/moderate/video - Moderate video (URL). Async, returns jobId.
- GET /v1/job/:jobId - Poll async moderation result.

## Review Endpoints
- POST /v1/review/appeal - User appeals a decision.
- POST /v1/review/human - Human overrides AI decision.

## Model Endpoints
- POST /v1/model/update - Push updated model weights.
- GET /v1/model/status - Current model versions and accuracy.
