# licensing-service — Configuration
LICENSING_PORT=4025, LICENSING_DB_URL, LICENSING_REDIS_URL, LICENSING_KAFKA_BROKERS, ROYALTY_CALCULATION_CRON="0 2 1 * *" (monthly), ROYALTY_PAYOUT_CRON="0 0 15 * *" (15th monthly), USAGE_RECORDING_RATE_LIMIT=1000/min, LICENSE_APPROVAL_TIMEOUT_HOURS=72 (auto-reject if not reviewed), CONTENT_RIGHTS_REGISTRATION_REQUIRED=true
FF: licensing_enabled=true, auto_royalty_calculation=true, auto_royalty_payout=true, geo_restriction_enforcement=true, usage_reporting=true, compliance_alerts=true
Rate limits: 100 license creations/day per rights holder, 1000 usage records/min (total), 10 royalty disputes/month per user
