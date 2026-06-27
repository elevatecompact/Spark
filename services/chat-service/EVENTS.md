# chat-service — Event Contracts
## Published: chat.message.sent, chat.message.deleted, chat.room.created, chat.room.closed, chat.user.banned, chat.user.muted
## Consumed: stream.session.ended (close room), moderation.filter.updated (update filter rules), identity.user.deleted (purge history)
## Schema: MessageSentEvent {messageId, roomId, userId, content, emotes[], sentAt, moderationResult}
