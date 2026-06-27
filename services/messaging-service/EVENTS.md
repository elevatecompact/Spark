# messaging-service — Event Contracts
## Published: messaging.conversation.created, messaging.message.sent, messaging.message.read, messaging.conversation.updated, messaging.attachment.uploaded
## Consumed: identity.user.deleted (cleanup), moderation.content.flagged (remove message), notification.push.sent (confirm delivery)
## Schema: MessageSentEvent {messageId, conversationId, senderId, contentType, content, replyTo, sentAt}
