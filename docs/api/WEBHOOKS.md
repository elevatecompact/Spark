# Webhooks

Webhooks provide event-driven HTTP callbacks that notify external systems about events occurring in the SPARK platform. They enable third-party integrations to react to platform events in real-time without polling.

## Event Types

Webhooks are available for a wide range of platform events. User events include user.created, user.updated, and user.deleted. Content events include content.created, content.updated, content.published, and content.deleted. Transaction events include transaction.completed, transaction.refunded, and transaction.failed. Subscription events include subscription.created, subscription.renewed, and subscription.cancelled. Moderation events include content.flagged, content.reviewed, and account.suspended.

## Delivery

Webhook delivery uses HTTP POST requests to the subscriber's configured endpoint URL. Each request includes a JSON payload with the event type, event ID, timestamp, and the relevant entity data. The request includes signature headers for verification. Delivery retries follow an exponential backoff schedule with a maximum of 10 retries over 24 hours.

## Security

Every webhook payload is signed using HMAC-SHA256 with the subscriber's secret key. The signature is sent in the X-Spark-Signature header. Subscribers must verify the signature before processing the payload. Endpoint URLs must use HTTPS. IP allowlisting is available for enhanced security.

## Management

Webhook subscriptions are managed through the developer dashboard or the REST API. Each subscription specifies the event types to receive, the endpoint URL, and optional filters. Webhook delivery logs are available for debugging with a 30-day retention period.
