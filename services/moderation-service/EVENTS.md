# moderation-service — Event Contracts
## Published: moderation.content.flagged (violation detected), moderation.action.taken (warn/restrict/remove/suspend), moderation.review.completed (human decision), moderation.rule.updated, moderation.report.submitted
## Consumed: chat.message.sent (scan in real-time), media.content.uploaded (scan image/video), creator.channel.updated (scan profile), community.post.created (scan text), stream.session.started (start content monitoring), identity.user.registered (scan for known bad actors)
## Schema: ModerationContentFlaggedEvent {contentId, contentType, scanResults[{ruleId, severity, category, confidence}], autoActionTaken, timestamp}
