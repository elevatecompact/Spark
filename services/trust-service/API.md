# trust-service — API Contract
## Reputation: GET /v1/reputation/{userId} (full score breakdown), GET /v1/reputation/{userId}/history (score changes over time), POST /v1/reputation/{userId}/recalculate (admin)
## Trust Signals: GET /v1/trust/signals/{userId} (all signals), POST /v1/trust/signals (record signal — internal), GET /v1/trust/level/{userId} (trust level: low,medium,high,verified)
## Risk Assessment: POST /v1/risk/assess (score an action context), GET /v1/risk/assessment/{id}, POST /v1/risk/rules (create — admin), PATCH /v1/risk/rules/{id} (admin)
## Fraud Detection: POST /v1/fraud/check-payment {transactionContext}, POST /v1/fraud/check-account {accountAction}, POST /v1/fraud/report {userId, reason}, GET /v1/fraud/cases (admin — open cases), POST /v1/fraud/cases/{id}/resolve (admin)
## Admin: GET /v1/admin/dashboard (trust metrics), GET /v1/admin/scores/distribution, POST /v1/admin/thresholds/update, GET /v1/admin/flagged-users
