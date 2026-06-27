# Secrets Management

Spark uses a centralized secrets management system to securely store, rotate, and audit access to sensitive credentials, API keys, tokens, and certificates.

## Architecture

Secrets are stored in a dedicated vault (HashiCorp Vault / equivalent) with a hardware security module (HSM) root of trust. Access is governed by policy and requires valid authentication. No secrets are stored in source code, configuration files, or environment variables.

## Secret Types

- **Service account credentials** — API keys, client secrets, database passwords
- **TLS certificates** — Private keys for internal and external TLS
- **Encryption keys** — Key encryption keys (KEKs) and data encryption keys (DEKs)
- **Third-party API tokens** — Tokens for integrated services
- **SSH keys** — Administrative access keys

## Access Control

Secret access follows least privilege and just-in-time (JIT) principles:

- **Read access** — Granted only to services and identities that explicitly require the secret
- **Write/rotate access** — Limited to automation pipelines and designated administrators
- **Dynamic secrets** — Database credentials are generated on-demand with short TTLs (max 24 hours)
- **Leasing** — Secrets have leases; services must renew to continue access

## Rotation

Automated rotation is enforced for all secrets:

| Secret Type | Rotation |
|---|---|
| Database passwords | Every 24 hours (dynamic) |
| API keys | Every 90 days |
| TLS keys | Every 90 days |
| Service tokens | On each deployment |

## Audit

Every secret access is logged with identity, timestamp, secret path, and operation type. Alerts fire on unauthorized access attempts, bulk reads, or secrets accessed outside of normal patterns.
