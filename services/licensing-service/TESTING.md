# licensing-service — Testing Guide
## Unit: License date validation (start < end, no overlap for exclusive), royalty calculation (flat, per-use, revenue share), territory matching (ISO country codes, region groups), usage deduplication.
## Integration: Full license lifecycle (create→approve→activate→use→report→royalty→payout), geo-blocking enforcement (mock geo-IP), content rights verification during stream start, royalty dispute workflow.
## Load: 1M usage records/month processing for royalty calculation, 1000 license checks/s during peak stream times.
