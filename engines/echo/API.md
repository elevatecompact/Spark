# API

## WebSocket
- ws://host/v1/ws?token=jwt - Connect with auth token.
- Text frames: JSON messages with type, channel, payload, id.
- Server pushes: type message, channel, payload, timestamp.

## REST
- POST /v1/channel/:channel/publish - Publish message to channel.
- POST /v1/channel/:channel/broadcast - Broadcast to all subscribers.
- GET /v1/channel/:channel/subscribers - List connected subscribers.
- POST /v1/user/:userId/send - Send direct message.
- DELETE /v1/session/:sessionId - Disconnect a session.
- GET /v1/health - Connection count, message throughput.
