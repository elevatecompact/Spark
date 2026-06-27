# Authentication Strategies

Spark supports multiple authentication mechanisms to accommodate a wide range of use cases while maintaining strong security guarantees. The authentication layer is built on open standards and designed for interoperability.

## OAuth 2.0 and OpenID Connect

OAuth 2.0 provides delegated authorization, while OpenID Connect (OIDC) extends it with identity and authentication capabilities. Spark uses OIDC as the primary authentication protocol for user-facing applications. All tokens are signed using RS256 or ES256 and have configurable expiration windows. Refresh tokens are bound to client secrets and rotated on each use.

## Passkeys (WebAuthn / FIDO2)

Passkeys offer phishing-resistant, passwordless authentication using public-key cryptography. Spark supports passkeys as both a primary and secondary authentication factor. User credentials never leave the device; only public keys are stored server-side. This eliminates credential theft via server breach.

## Multi-Factor Authentication (MFA)

Spark enforces MFA for all privileged operations, administrative access, and high-risk transactions. Supported factors include time-based one-time passwords (TOTP), SMS codes, hardware security keys (FIDO2/WebAuthn), and push notifications.

## Session Management

Sessions are tracked via short-lived access tokens (15 minutes) and longer-lived refresh tokens (7 days). Refresh token rotation and revocation are enforced. Concurrent session limits and geographic anomaly detection prevent token abuse.

## Password Policies

For legacy password-based authentication, Spark enforces NIST SP 800-63B guidelines: minimum 12 characters, no composition rules, credential checks against breached password databases, and rate limiting after 5 failed attempts.
