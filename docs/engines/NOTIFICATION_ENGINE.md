# Notification Engine — Multi-Channel Notifications

## Purpose

The Notification engine delivers time-sensitive communications to Titan users across multiple channels — push notifications (FCM, APNs), email (SES), SMS, and in-app notifications (Echo messaging). It handles templating, delivery routing, and preference management.

## Architecture

Notifications are submitted to a unified API, which routes them through channel-specific delivery pipelines. Delivery status is tracked end-to-end with retries, fallbacks, and aggregate delivery analytics.

## Tech Stack

- **Language**: Go
- **Push**: Firebase Cloud Messaging (Android), Apple Push Notification Service (iOS)
- **Email**: Amazon SES with HTML template rendering
- **SMS**: Twilio for transactional SMS
- **In-app**: Echo engine (WebSocket delivery)
- **Queue**: RabbitMQ for async delivery with per-channel queues
- **Template Engine**: Go templates for email content with localization

## Key Features

- **Unified API**: Single notification submission endpoint that handles channel routing
- **Channel preference**: Per-user and per-notification-type channel preferences with hierarchical overrides
- **Template system**: Versioned notification templates with locale, personalization, and conditional blocks
- **Batching**: Coalesce multiple notifications into digest emails and push summaries
- **Delivery tracking**: End-to-end delivery status (submitted, delivered, opened, clicked) per notification
- **Rate limiting**: Per-user and per-channel rate limiting to prevent notification fatigue
- **Fallback routing**: If push fails, in-app notification, then email (configurable chain)
- **Unsubscribe**: One-click unsubscribe with preference persistence and compliance (CAN-SPAM, GDPR)
- **Scheduled delivery**: Time-zone-aware delivery scheduling (send during user's daytime hours)

## Performance Targets

| Metric | Target |
|--------|--------|
| Push delivery latency | < 5 seconds (p99) |
| Email delivery latency | < 30 seconds (p95) |
| Delivery success rate | > 99% |
| Throughput | 100,000 notifications/second |
| Template rendering | < 10ms per notification |