# moderation-service — Troubleshooting
## False positives: ML model threshold too low, rule too broad, training data not representative. Check sensitivity vs specificity metrics, adjust threshold, retrain model with more diverse data.
## False negatives (missed violations): Threshold too high, model miss, rule not covering new pattern. Review recent missed violations, lower threshold, add new rule pattern.
## Queue backing up: Insufficient moderator staffing, review tool latency, auto-resolution not configured. Scale human review team, verify dashboard performance, enable auto-approve for low-severity with confidence > 0.95.
## ML model down: GPU OOM, model version mismatch, inference server crash. Restart ML pod, verify model artifacts in S3, rollback model version.
