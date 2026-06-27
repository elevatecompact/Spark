# subscription-service — Event Contracts
## Published: subscription.activated, subscription.renewed, subscription.cancelled, subscription.expired, subscription.upgraded, subscription.downgraded, subscription.payment.failed, subscription.grace.ended
## Consumed: wallet.transaction.settled (confirm payment), creator.tier.pricing_changed (update billing), identity.user.deleted (cancel all)
## Schema: SubscriptionActivatedEvent {subscriptionId, userId, planId, creatorId, billingPeriod, amountCents, currency, activatedAt, nextBillingAt}
