# trust-service — Event Contracts
## Published: trust.reputation.updated, trust.trust_level.changed (escalation or demotion), trust.risk.alert (high-risk action detected), trust.fraud.case.opened, trust.fraud.case.resolved, trust.signal.recorded
## Consumed: identity.user.registered (initial trust setup), wallet.transaction.failed (negative signal), moderation.action.taken (negative signal), payment.dispute.opened (negative signal), subscription.activated (positive signal), gift.sent (positive signal), viewer.report.submitted (negative signal for reported user), commerce.order.fulfilled (positive signal)
## Schema: TrustReputationUpdatedEvent {userId, overallScore(0-1000), trustLevel, signalsCount, positiveWeight, negativeWeight, calculatedAt}
