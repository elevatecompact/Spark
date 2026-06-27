# trust-service — Changelog
## [1.2.0] — 2026-06-15
- ML-based fraud detection (XGBoost model for payment fraud), device fingerprinting for account takeover prevention, IP reputation integration, real-time risk assessment on all high-value actions. Reputation decay factor (daily 0.95) to reduce stale signal weight.
## [1.1.0] — 2026-04-08
- Risk rule engine (configurable conditions + actions), fraud case management workflow, trust level escalation/demotion with notifications, reputation recalculation batch job (daily).
## [1.0.0] — 2026-01-20
- Reputation scoring (0-1000), trust signals ingestion from platform events, trust level classification, basic risk assessment API.
