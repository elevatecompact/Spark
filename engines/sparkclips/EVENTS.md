# Events

## Published Events
- highlights.detection.started/progress/completed/failed.
- clip.generated - Payload: { clipId, startTime, endTime, score, outputUri }.
- clip.published - Payload: { clipId, platform, publishedUrl }.
- clip.rated - Payload: { clipId, rating }.
- highlight.shareability.estimated.

## Subscribed Events
- stream.started - Begin live detection.
- stream.ended - Finalise highlights.
- social.account.linked - Social auth for publishing.
