# chat-service — Ownership
**Real-Time Infrastructure** — eng-realtime@titan.dev. TL: Hassan Ali, WS: Yuki Tanaka, SRE: Omar Farouk.
Owns: WebSocket gateway, message persistence, chat moderation, emote system, pub/sub fan-out.
Deps: stream-service (auto-create rooms), moderation-service (filters), identity-service (JWT auth), notification-service (@mentions), analytics-service (volume metrics).
