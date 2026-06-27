# translation-service — Troubleshooting
## Translation quality poor: Wrong provider selected, language detection incorrect, translation memory stale. Check PRIMARY_PROVIDER setting, verify language detection accuracy, review translation memory freshness.
## Provider API errors: Rate limit exceeded, API key expired, network issue. Check provider console for rate limits, verify key rotation, check Vault for API key validity.
## Real-time translation slow: WebSocket buffer filling, provider latency high, streaming buffer too large. Reduce STREAMING_BUFFER_SIZE, switch provider, scale streaming pods.
## Translation memory not matching: Exact match expected but fuzzy search returning. Check source_hash computation, verify normalization (lowercase, trim, unicode NFC).
