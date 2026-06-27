# chat-service — Testing Guide
## Unit: Message validation, room permission logic, emote parsing, slow mode rate calculation.
## Integration: WS connect/disconnect, message round-trip, room CRUD with permissions, moderation actions, history pagination.
## Load: 10K concurrent WS connections, 1000 msgs/s in single room, emote-heavy throughput.
## Tools: WebSocket test client in tests/clients/. Emulate multiple users with NewTestChatUser().
