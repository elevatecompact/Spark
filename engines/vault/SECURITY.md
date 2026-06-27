# Security

- PCI-DSS Level 1 compliance. Card data never touches servers - tokenized by Stripe/PayPal.
- AES-256-GCM encryption at rest with KMS key rotation every 90 days.
- Idempotency prevents replay attacks.
- Webhook verification via provider-specific signatures.
- PII minimisation - full card details never persisted.
- Refund/plan endpoints require billing:admin role.
