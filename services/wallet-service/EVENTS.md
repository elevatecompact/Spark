# wallet-service — Event Contracts
## Published: wallet.transaction.created, wallet.transaction.settled, wallet.transaction.failed, wallet.deposit.received, wallet.withdrawal.processed, wallet.payout.completed, wallet.balance.low
## Consumed: subscription.payment.due (process recurring), gift.purchase.completed (credit recipient), commerce.order.placed (hold funds), identity.user.deleted (close wallet)
## Schema: TransactionCreatedEvent {transactionId, fromWalletId, toWalletId, amountCents, currency, type, status, createdAt}
