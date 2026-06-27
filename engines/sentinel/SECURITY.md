# Security

- Adversarial input detection for prompt injection.
- Model isolation - each model runs in sandboxed process.
- Decision audit trail in immutable ClickHouse store.
- Human reviewer privacy - only flagged content visible, not user identity.
- Policy access control - only admins update blocklists or deploy models.
- Evasion monitoring for homoglyph and Unicode attacks.
