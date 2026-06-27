# Notifications

The notification system delivers timely, relevant updates to users across multiple channels. It supports push notifications, in-app notifications, email digests, and webhook integrations.

## Notification Types

Social notifications include new follower, content liked, comment received, and mention alerts. Content notifications include content published, live stream started, and scheduled event reminders. System notifications include account updates, security alerts, and platform announcements. Transactional notifications include subscription renewals, payment confirmations, and payout notifications. Moderation notifications include content removal, warning notices, and appeal updates.

## Delivery Channels

Push notifications are delivered through Firebase Cloud Messaging for Android and Apple Push Notification Service for iOS, with configurable sound, badge, and alert settings. In-app notifications appear in the notification center accessible from the main navigation, with read/unread tracking and bulk actions. Email notifications provide digest summaries for low-priority notifications and individual emails for high-priority alerts. WebSocket notifications enable real-time delivery for users currently on the platform.

## Preference Management

Users have granular control over notification preferences. Channel-level preferences allow enabling or disabling push, email, and in-app delivery independently. Category-level preferences control which notification categories are received. Frequency controls allow users to choose between instant delivery, daily digest, weekly digest, or disabled. Quiet hours prevent notifications during specified time periods.

## Delivery Infrastructure

The notification service processes notifications through a queue-based architecture. Notifications are created, enriched with template data, and dispatched through the appropriate channel providers. Delivery tracking logs success, failure, and delivery latency. Failed deliveries are retried with exponential backoff.
