# subscription-service — Deployment Guide
Requires PG15+ with cron extension, Redis 7+, Vault for API keys.
K8s: 3 replicas 1GB/500m, HPA min 2 max 8, CronJob for daily billing at 2AM.
Deploy: ./subscription migrate up, ./subscription seed-plans, kubectl apply -f k8s/subscription-service/.
Billing CronJob processes renewals, retries, and expires grace periods.
