# Security

- Model access control - model files signed and verified before loading.
- User feature isolation - Redis keys namespaced by tenant.
- gRPC TLS - mutual TLS for all internal RPC.
- Input validation - sanitized against injection and adversarial manipulation.
- Feedback rate limiting - per-user limit of 100 events/minute.
- PII scrubbing - feature pipeline strips PII before storage.
