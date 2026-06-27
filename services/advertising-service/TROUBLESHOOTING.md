# advertising-service — Troubleshooting
## Ads not serving: No eligible campaigns, targeting too narrow, budget exhausted, inventory not configured. Check active campaigns matching placement, verify targeting params, confirm budget remaining.
## Low fill rate: Not enough campaigns bidding, floor price too high, inventory mismatched. Lower floor price, increase campaign targeting breadth, review inventory categorization.
## Impression counting incorrect: ClickHouse ingestion lag, deduplication failure, fraud filter too aggressive. Check impression_worker lag, verify dedup logic, review fraud threshold.
## CTR abnormally high/low: Fraud (bot clicks), ad creative issue, targeting mismatch. Investigate click patterns, review creative quality, analyze audience targeting accuracy.
