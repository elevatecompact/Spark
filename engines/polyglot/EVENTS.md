# Events

## Published Events
- translation.requested/completed/failed - Request lifecycle.
- translation.quality.low - Score below threshold.
- model.loaded - Payload: { modelId, languagePair, loadTimeMs }.
- stream.translation.started/ended - Streaming session lifecycle.
- glossary.updated - Payload: { glossaryId, entryCount }.

## Subscribed Events
- content.translate.request, caption.translate.request.
- terminology.updated.
