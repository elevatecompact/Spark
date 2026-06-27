# Runbook

## Startup
1. Verify provider keys loaded from Vault.
2. Run migrations: ./vault migrate --config config/vault.toml.
3. Start API: ./vault serve.
4. Verify: POST /v1/payment/intent.
5. Register webhook endpoints.

## Monitoring
- Dashboard: success rate, revenue, refund rate, dunning recovery.
- Alerts: success rate < 95%, reconciliation failure, webhook lag > 5min.
