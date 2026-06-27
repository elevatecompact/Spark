# translation-service — Runbook
## Alerts: TranslationLatency > 5s, ProviderErrorRate > 5% (per provider), TranslationMemoryHitRate < 30%, BatchJobFailure > 2%, RealtimeLatency > 3s
## Switch provider: POST /v1/admin/provider/switch {provider: "google"} — emergency if DeepL down.
## Clear cache: POST /v1/admin/clear-cache — flushes Redis translation cache.
## Warm memory: Pre-compute translations for top 1000 content items into top 10 languages.
## Check usage: GET /v1/admin/usage — monitor character count vs API plan limits.
