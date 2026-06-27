# trust-service — Runbook
## Alerts: ReputationScoreStale > 24h (recalculation job failed), FraudDetectionAlerts > 100/h (possible attack), RiskAssessmentLatency > 500ms, TrustLevelChangesSpike > 1000/h (possible score manipulation)
## Recalculate user reputation: POST /v1/admin/reputation/{userId}/recalculate
## Override trust level: POST /v1/admin/trust/level/override {userId, level: "verified"} — for known trusted users.
## Review fraud cases: GET /v1/admin/fraud/cases?status=open — review and resolve.
## Update risk rule: PATCH /v1/admin/risk/rules/{id} {is_active: false} — disable problematic rule.
## Check score distribution: GET /v1/admin/scores/distribution — monitor overall platform health.
