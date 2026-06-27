# Events

## Published Events
- session.connected - Payload: { sessionId, userId, device, ip, connectedAt }.
- session.disconnected - Payload: { sessionId, reason, duration, messagesSent }.
- channel.subscribed - Payload: { sessionId, channel, subscribedAt }.
- channel.unsubscribed - Payload: { sessionId, channel }.
- message.published - Payload: { messageId, channel, senderId, size, timestamp }.
- message.delivered - Payload: { messageId, sessionId, deliveryLatencyMs }.

## Subscribed Events
- system.maintenance - Broadcast maintenance notification.
- user.presence.request - Query user presence status.
