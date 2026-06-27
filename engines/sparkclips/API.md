# API

## Highlight Detection
- POST /v1/highlights/detect - Detect highlights. Returns jobId.
- GET /v1/highlights/job/:jobId - Poll status.
- POST /v1/highlights/generate - Generate clip from segment.

## Clip Management
- GET /v1/clip/:clipId - Get metadata and download URL.
- DELETE /v1/clip/:clipId - Delete clip.
- POST /v1/clip/:clipId/publish - Publish to social platforms.

## Templates
- POST /v1/template - Create rendering template.
- GET /v1/template - List templates.

## Feedback
- POST /v1/highlights/feedback - Rate highlight quality.
