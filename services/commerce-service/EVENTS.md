# commerce-service — Event Contracts
## Published: commerce.product.created, commerce.product.updated, commerce.product.deleted, commerce.cart.updated, commerce.order.placed, commerce.order.fulfilled, commerce.order.cancelled, commerce.order.refunded, commerce.review.submitted
## Consumed: wallet.transaction.settled (confirm payment → fulfill), notification.push.sent (delivery confirmation), payment.refund.processed (process refund order), identity.user.deleted (cancel pending orders)
## Schema: CommerceOrderPlacedEvent {orderId, buyerId, merchantId, items[{productId, variantId, quantity, priceCents}], totalCents, currency, status, placedAt}
