# Runbook

## Startup
1. Run migrations: ./guardian migrate --config config/guardian.toml.
2. Generate keys: ./guardian keys generate --output /etc/guardian/keys/.
3. Start Guardian: ./guardian serve --config config/guardian.toml.
4. Verify: curl http://localhost:8080/v1/token/public-key.
5. Test: POST /v1/auth/login.

## Monitoring
- Dashboard: token validation rate, login success rate, cache hit ratio.
- Alerts: login success rate drop > 20%, latency > 50ms, signing key expiry < 7 days.
