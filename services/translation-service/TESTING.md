# translation-service — Testing Guide
## Unit: Language detection accuracy, translation memory lookup (exact/fuzzy), quality score computation, provider selection logic (primary→fallback), text sanitization (HTML strip, length truncation).
## Integration: Translate→cache→lookup, batch processing, WebSocket streaming, review workflow cycle, provider failover (mock DeepL down→verify Google used).
## Quality: BLEU score comparison between providers, human evaluation of sampled translations, latency comparison per provider.
