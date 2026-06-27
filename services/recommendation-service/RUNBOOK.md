# recommendation-service — Runbook
## Alerts: InferenceLatency > 200ms, FeedCTRDrop > 20% (model regression), FeatureCacheMiss > 30%, ModelVersionRollback, EmbeddingCorruption
## Refresh features: POST /v1/admin/refresh-features
## Rollback model: POST /v1/admin/models/deploy {version: "previous"}, monitor CTR recovery.
## Warm cache: Pre-compute top 10K user feeds and cache in Redis.
## Evaluate: Run offline eval on current vs candidate model.
