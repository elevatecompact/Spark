# analytics-service — Testing Guide
## Unit: Aggregation window logic, rollup function correctness, funnel step calculation, cohort retention math.
## Integration: Event ingestion→ClickHouse query, dashboard rendering correctness, report generation pipeline, alert threshold evaluation.
## Data quality: Compare raw event count vs aggregated metric count (should match within 0.1%).
## Tools: k6 for event ingestion load test. Locust for dashboard query load.
