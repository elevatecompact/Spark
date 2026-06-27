# API Documentation

This directory contains comprehensive documentation for SPARK's API ecosystem. The platform exposes multiple API styles to support diverse client needs, from browser-based web applications to mobile clients, third-party integrations, and real-time services.

## API Styles

| API | Protocol | Use Case |
|-----|----------|----------|
| GraphQL | HTTP/2 | Primary API for web and mobile clients |
| REST | HTTP/1.1 | Third-party integrations, simple CRUD |
| gRPC | HTTP/2 | Internal service-to-service communication |
| WebSocket | WS/WSS | Real-time notifications, chat, live events |
| WebRTC | UDP/TCP | Peer-to-peer media streaming |
| Webhooks | HTTP | Event-driven outgoing integrations |

## Key Documents

- **API_GUIDELINES.md** — Design principles, conventions, and best practices
- **AUTHENTICATION.md** — Authentication flows and token management
- **RATE_LIMITING.md** — Rate limiting policies and headers
- **PAGINATION.md** — Cursor and offset pagination patterns
- **ERROR_CODES.md** — Standard error response format and codes
- **VERSIONING.md** — API versioning strategy and deprecation policy

Each API style has a dedicated document with detailed specification, usage patterns, and implementation guidance.
