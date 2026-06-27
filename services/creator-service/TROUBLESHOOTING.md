# creator-service — Troubleshooting

## Channel Creation Fails
**Causes:** Name collision (UNIQUE constraint), max channels per creator exceeded, identity service unreachable. Check channels table for duplicates, verify identity service health.

## Verification Docs Not Uploading
**Causes:** S3 bucket CORS misconfiguration, presigned URL expired (> 24h), file > 10MB. Verify S3 policy, regenerate URL, check client-side file size limits.

## Metrics Stale
**Causes:** Materialized view refresh cron job failed. Manually refresh: ./creator refresh-metrics. Check cron job logs in Kubernetes.

## Payout Failed
**Causes:** Wallet service unavailable, insufficient platform balance, invalid payout method. Check wallet-service health, verify payout method status, review payout logs.
