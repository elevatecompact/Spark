# Security

- Provider credentials stored in HashiCorp Vault; retrieved via Vault agent sidecar.
- Purge API requires JWT with purge:write scope; rate-limited.
- Steering rules signed and verified; only admin roles can modify.
- Origin shield in dedicated VPC with strict ingress rules.
- All purge, warm, steering changes logged with user identity.
