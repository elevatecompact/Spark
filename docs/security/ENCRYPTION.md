# Encryption at Rest and in Transit

Spark employs industry-standard encryption to protect data in all states. The encryption architecture is designed to meet the requirements of SOC2, GDPR, HIPAA, and PCI DSS.

## Encryption in Transit

All network communication uses TLS 1.3 with strong cipher suites (TLS_AES_256_GCM_SHA384 and TLS_CHACHA20_POLY1305_SHA256). Perfect forward secrecy (PFS) is enforced via ECDHE key exchange. mTLS is required for inter-service communication within the Spark mesh. Certificate management uses automated issuance via ACME with 90-day certificate lifetimes and automatic renewal.

## Encryption at Rest

- **Database Encryption** — All databases use AES-256-GCM transparent data encryption (TDE). Standalone encryption keys are managed by a Hardware Security Module (HSM).
- **Object Storage** — All objects are encrypted with server-side AES-256 using a unique per-object key. The per-object key is itself encrypted with a master key stored in the key management system.
- **Backup Encryption** — All backups are encrypted with AES-256-CBC before leaving the source environment. Encryption keys are stored separately from backup data.

## Key Management

Spark uses a dedicated key management service (KMS) with automatic key rotation schedules:

| Key Type | Rotation Cadence |
|---|---|
| TLS certificate private keys | Every 90 days |
| Database master keys | Every 365 days |
| Object storage master keys | Every 180 days |
| User data encryption keys | Every 365 days |

## Client-Side Encryption

For sensitive fields (PII, payment data), Spark supports application-layer encryption where data is encrypted by the client SDK before transmission. The server never has access to the raw plaintext of these fields.
