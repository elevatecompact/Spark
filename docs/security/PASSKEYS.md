# Passkey Authentication

Spark supports passkeys (WebAuthn / FIDO2) as a passwordless, phishing-resistant authentication method. Passkeys use public-key cryptography to eliminate credential theft and provide seamless cross-device authentication.

## How It Works

When a user registers a passkey, their device generates a public-private key pair. The private key never leaves the device and is protected by the platform biometric or PIN. The public key is registered with Spark. During authentication, the device signs a challenge, proving possession of the private key without ever transmitting it.

## Implementation

Spark implements the WebAuthn Level 2 specification. The registration and authentication ceremonies follow the W3C WebAuthn standard. Platform authenticators (Windows Hello, Touch ID, Face ID, Android biometric) and roaming authenticators (YubiKey, Titan Key) are both supported.

## Storage

Public keys are stored in a dedicated, encrypted database table. Each credential record includes the credential ID, public key (in COSE format), sign count, device type, and creation timestamp. Private keys are never transmitted to or stored by Spark.

## Security Benefits

- **Phishing resistant** — Passkeys are bound to the origin, preventing use on fraudulent sites
- **No shared secrets** — Server stores only public keys; credential theft is impossible
- **Revocable** — Users can revoke individual passkeys from their account settings
- **Cross-platform** — Passkeys sync across devices via platform providers (iCloud Keychain, Google Password Manager, Windows Hello)

## Recovery

Users who lose all passkeys can regain access via recovery codes (generated at enrollment), backup email OTP, or an administrative reset. Recovery events are logged and alert the account holder.
