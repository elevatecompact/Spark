# commerce-service — Troubleshooting
## Checkout fails: Payment service down, inventory unavailable, cart expired. Check payment-service health, verify inventory count, confirm cart hasn't expired (TTL 30min idle).
## Inventory oversold: Race condition on concurrent purchases, optimistic lock missed. Check product.inventory_count vs sold total, reconcile with actual count, set lower concurrent purchase limit.
## Digital download not working: Presigned URL expired, S3 bucket permissions wrong, fulfillment not processed. Check URL expiry (72h default), verify S3 bucket policy, check fulfillment status.
## Merchant payout missing: Payout cron not run, wallet service unavailable, insufficient platform balance. Check cron logs, verify wallet health, confirm platform balance for payouts.
