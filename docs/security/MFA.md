# Multi-Factor Authentication

Spark enforces multi-factor authentication (MFA) for all privileged accounts and provides optional MFA for standard users. The MFA system supports multiple factors and adaptive policies.

## Supported Factors

- **Time-Based One-Time Password (TOTP)** — RFC 6238 compliant, 30-second window, 6-digit codes. Users can enroll via authenticator apps (Google Authenticator, Authy, 1Password, etc.).
- **SMS / Email OTP** — Fallback factor for users without access to their primary authenticator. Rate-limited to 3 attempts per 5 minutes.
- **Hardware Security Keys** — FIDO2 / WebAuthn security keys (YubiKey, Google Titan Key) for phishing-resistant MFA.
- **Push Notifications** — Mobile push via the Spark mobile app. Users approve or deny login attempts from a trusted device.

## Adaptive MFA Policies

MFA requirements adapt based on risk context:

| Scenario | MFA Required |
|---|---|
| Login from new device | Yes |
| Login from known device / location | No (if low risk) |
| Privileged operation (admin panel) | Yes |
| Export of sensitive data | Yes |
| Password change | Confirm existing + MFA |
| High-value financial transaction | Yes |

## Enrollment

Users are prompted to enroll in MFA on first login. A grace period of 7 days is allowed, after which MFA becomes mandatory for accounts flagged as high-risk or privileged.

## Recovery

MFA recovery uses a pre-generated set of 10 single-use recovery codes. Each code is hashed with bcrypt before storage. When fewer than 3 codes remain, the user is prompted to regenerate them. Account recovery through support requires identity verification via government-issued ID.
