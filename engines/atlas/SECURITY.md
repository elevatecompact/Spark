# Security

- Mutual TLS required for all gRPC connections between agents and registrars.
- Registration tokens - services present a signed token to register; tokens encode service identity and allowed tags.
- Watch authentication - SSE watchers must present a valid API key in the X-Atlas-Token header.
- etcd encryption - etcd peer and client traffic encrypted via TLS; etcd data encrypted at rest.
- Audit logging - all registration and deregistration events logged with originating service identity.
- Rate limiting - per-service registration rate limited to 10 req/s; per-IP discovery rate limited to 100 req/s.
