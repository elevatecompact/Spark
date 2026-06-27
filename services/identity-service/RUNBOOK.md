# identity-service — Runbook

## Monitoring
**Grafana Dashboard:** "Identity Service" — shows auth request rate by method, success/failure ratio, MFA enrollment rate, latency percentiles, and error budgets.

**Alerts:**
- IdentityErrorRate > 1% over 5m → PagerDuty critical
- IdentityP99Latency > 500ms → Slack #oncall-identity
- IdentityDBConnections > 80% → Auto-scale event
- IdentityLoginFailures > 100/min — Possible brute force

## Common Procedures
### Force Session Revocation
`ash
curl -X POST http://identity:4001/v1/admin/users/{id}/revoke-sessions \
  -H "Authorization: Bearer "
`

### Flush Token Cache
`ash
redis-cli --scan --pattern "session:*" | ForEach-Object { redis-cli del  }
`

### Emergency Registration Disable
Set feature flag egistration_open=false via ConfigMap update.

## On-Call Checklist
1. Check dashboards for anomaly window
2. Review recent auth-related deployments
3. Verify upstream DB/Redis health
4. Check Vault for secret rotation events
