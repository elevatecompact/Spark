# notification-service — Event Contracts
## Published: notification.delivered, notification.read, notification.bounced (delivery failure), notification.preferences.updated
## Consumed: iam.user.registered (welcome email), subscription.activated (new sub alert), gift.received (gift notification), stream.session.started (stream live alert), wallet.payout.completed (payout confirmation), chat.message.sent (@mention), messaging.message.sent (DM alert)
## Schema: NotificationDeliveredEvent {notificationId, userId, channel(push|email|sms|inapp), type, deliveryStatus(delivered|failed|bounced), deliveredAt}
