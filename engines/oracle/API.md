# API

## Recommendation Endpoints
- POST /v1/recommend - Request recommendations. Returns ranked list of item IDs with scores.
- POST /v1/recommend/personalized - With real-time session context override.
- POST /v1/recommend/similar/:itemId - Get similar items based on content embedding.

## Feedback Endpoints
- POST /v1/feedback - Record user interaction.
- POST /v1/feedback/batch - Batch feedback submission.

## Model Endpoints
- POST /v1/model/deploy - Deploy a new model version.
- GET /v1/model/active - Return currently active model metadata.
