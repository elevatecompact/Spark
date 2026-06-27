# payment-service — Event Contracts
## Published: payment.intent.created, payment.intent.succeeded, payment.intent.failed, payment.refund.processed, payment.dispute.opened, payment.dispute.resolved, payment.payout.completed, payment.method.expiring
## Consumed: wallet.transaction.created (initiate payment), subscription.payment.due (recurring), commerce.order.placed (capture)
## Schema: PaymentSucceededEvent {paymentIntentId, externalId("pi_xxx"), processor, amountCents, currency, status, paymentMethodType, succeededAt}
