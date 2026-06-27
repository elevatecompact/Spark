# OAuth 2.0 / OpenID Connect Implementation

Spark uses OAuth 2.0 and OpenID Connect (OIDC) as the standard protocol for delegated authorization and authentication. The implementation follows RFC 6749, RFC 7519, and the OIDC Core 1.0 specification.

## Supported Grant Types

- **Authorization Code + PKCE** — Primary flow for public clients (SPAs, mobile apps). PKCE (Proof Key for Code Exchange) ensures authorization codes cannot be intercepted.
- **Client Credentials** — Machine-to-machine communication. Clients authenticate with a client ID and secret.
- **Refresh Token** — Long-lived access without requiring user re-authentication. Refresh tokens are sender-constrained (DPoP binding).

## Token Formats

- **Access Token** — JWT format, signed with RS256, containing claims: `iss`, `sub`, `aud`, `exp`, `iat`, `scope`, `sid`, and custom claims for tenant context.
- **ID Token** — OIDC identity token containing user profile claims. Never used for API authorization.
- **Refresh Token** — Opaque token stored server-side as a SHA-256 hash. Rotated on each use.

## Authorization Server

Spark operates a dedicated OAuth 2.0 Authorization Server built on the identity platform. It handles token issuance, introspection, revocation, and public key publication at `/.well-known/openid-configuration` and `/.well-known/jwks.json`.

## Security Controls

- Token binding via DPoP (OAuth 2.0 Demonstrating Proof of Possession)
- Redirect URI strict validation (exact match, no wildcards)
- Authorization code binding to PKCE challenge
- Refresh token rotation and reuse detection
- Client authentication with private_key_jwt for confidential clients
