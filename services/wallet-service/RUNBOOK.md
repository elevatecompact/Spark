# wallet-service — Runbook
## Alerts: TransactionFailureRate > 0.5% (immediate), ReconciliationMismatch ( tolerance), PayoutFailureRate > 2%, WalletDBReplicaLag > 10s
## Manual reconcile: ./wallet reconcile --date 2026-06-26
## Freeze wallet: POST /v1/admin/wallets/{id}/freeze
## Emergency: Disable all withdrawals via feature flag. Halt payment processing on discrepancy.
