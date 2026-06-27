# wallet-service — Testing Guide
## Unit: Balance calc with optimistic locking, idempotency enforcement, currency math (6 decimal precision), payout eligibility.
## Integration: Deposit→settle→balance verify, transfer with checks, failed transaction rollback, withdrawal lifecycle, daily reconciliation.
## Property-based: balance=sum(incoming)-sum(outgoing) invariant, no double-spend, transaction total matches ledger delta.
