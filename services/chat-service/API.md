# chat-service — API Contract
## WebSocket: GET /ws/chat/{roomId}?token=<jwt> — Bidirectional event stream
WS Events: message:send, message:received, user:join, user:leave, mod:action, room:update

## REST
### Rooms: POST /v1/rooms, GET /v1/rooms/{id}, DELETE /v1/rooms/{id}
### Messages: POST /v1/rooms/{id}/messages, GET /v1/rooms/{id}/messages (cursor paginated), PUT /v1/rooms/{id}/messages/{msgId}, DELETE /v1/rooms/{id}/messages/{msgId}
### Moderation: POST /v1/rooms/{id}/mute/{userId}, POST /v1/rooms/{id}/ban/{userId}, POST /v1/rooms/{id}/slow-mode
### Emotes: GET /v1/emotes (global), GET /v1/rooms/{id}/emotes (room-specific)
