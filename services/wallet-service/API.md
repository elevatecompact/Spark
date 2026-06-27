# wallet-service — API Contract
## Balances: GET /v1/wallets/me, GET /v1/wallets/me/balances (all currencies)
## Transactions: POST /v1/transactions/deposit, /withdraw, /transfer, /tip, /purchase, GET /v1/transactions (filterable), GET /v1/transactions/{id}
## Payouts: POST /v1/payouts/request, GET /v1/payouts, GET /v1/payouts/{id}
## Admin: GET /v1/admin/ledger, POST /v1/admin/reconcile, GET /v1/admin/audit-log
All amounts in cents. Idempotency-Key header required on all mutations.
