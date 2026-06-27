# trust-service — Troubleshooting
## Reputation score not updating: Signal not being consumed, recalculation cron missed, score calculation error. Check Kafka consumer lag for trust signals, verify cron job ran, check recalculation logs for errors.
## Fraud false positives: Risk rule too aggressive, threshold too low, ML model not calibrated. Review triggered rules, check feature distributions, recalibrate model thresholds.
## Trust level wrong: Score threshold misconfigured, signals expired prematurely, weight assignment issue. Verify TRUST_LEVEL_THRESHOLDS config, check signal expiration logic, review weight assignments per category.
## Risk assessment slow: Too many rules evaluated, Redis cache miss, downstream dependency slow. Review rule complexity, warm risk assessment cache, check Redis performance.
