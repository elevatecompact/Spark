# Events

## Published Events
- moderation.decision - Payload: { decision, contentId, userId, reasons, modelUsed, latencyMs }.
- moderation.flagged - Payload: { contentId, severity, reviewUrgency } for human review.
- moderation.appealed - Payload: { contentId, userId, reason }.
- moderation.overridden - Payload: { contentId, previousDecision, newDecision, reviewerId }.
- model.drift.detected - Production accuracy dropped below baseline.

## Subscribed Events
- user.reported - Expedited moderation on user report.
- policy.updated - Update policies and rules in real-time.
