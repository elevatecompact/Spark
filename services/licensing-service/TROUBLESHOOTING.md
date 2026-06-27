# licensing-service — Troubleshooting
## License check failing for valid content: Cache stale, rights not registered, license expired. Flush Redis cache, verify content_rights registration, check license end_date. Override temporarily via admin API.
## Royalty calculation wrong: Usage log incomplete, rate misconfigured, period boundary error. Verify usage_log entries for period, check license rate_type and rate_cents, ensure period start/end alignment.
## Content usage not tracked: Kafka event lost, usage recording rate limited, missing license association. Check dead-letter queue for failed usage events, verify rate limit counters, confirm content_id has active license.
