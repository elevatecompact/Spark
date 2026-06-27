# API

## Translation
- POST /v1/translate - Translate text. Returns translatedText, detectedSourceLang, confidence.
- POST /v1/translate/batch - Batch up to 100 texts.
- POST /v1/translate/stream - WebSocket for streaming translation.
- POST /v1/translate/detect - Language detection.

## Glossary
- PUT /v1/glossary - Create glossary.
- GET /v1/glossary/:id - Retrieve entries.

## Quality
- POST /v1/quality/estimate - Estimate quality.
- POST /v1/quality/feedback - Submit feedback.

## Model
- GET /v1/model/languages - Supported pairs.
- POST /v1/model/swap - Swap model for A/B test.
