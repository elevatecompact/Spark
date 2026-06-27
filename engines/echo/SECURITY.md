# Security

- JWT authentication - every WebSocket connection validates a signed JWT with userId, allowedChannels, expiry.
- Channel authorization - server-side ACL checks before subscription.
- Message validation - all messages validated against JSON schema.
- Rate limiting - per-connection: 10 msg/s, per-channel: 100 msg/s.
- Origin checking - WebSocket Origin header validated against allowlist.
- WSS enforced at load balancer; internal gRPC uses mutual TLS.
