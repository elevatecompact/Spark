# API Authentication

SPARK APIs use token-based authentication with JWT (JSON Web Tokens) for all API access. This document describes the authentication flows and token management procedures.

## Authentication Flows

**OAuth 2.0 Authorization Code Flow** is the primary authentication mechanism for user-facing applications. The flow begins with the client redirecting the user to the authorization endpoint. Upon successful authentication, the user is redirected back with an authorization code. The client exchanges this code for an access token and a refresh token. This flow supports PKCE for public clients such as mobile and single-page applications.

**Client Credentials Flow** is used for server-to-server communication. The client authenticates with its client ID and client secret to obtain an access token. This flow is appropriate for backend services and cron jobs that act on behalf of the application rather than a specific user.

## Token Format

Access tokens are JWTs containing the issuer, subject, audience, expiration time (1 hour), issued-at time, and custom claims including user roles and permissions. Tokens are signed using RS256 with a 2048-bit key pair. Public keys for signature verification are available at the JWKS endpoint.

## Token Management

Access tokens expire after 1 hour. Refresh tokens expire after 30 days and are single-use. When a refresh token is used, both a new access token and a new refresh token are returned, invalidating the previous refresh token. If a refresh token is compromised, it can be revoked through the developer dashboard. Token revocation propagates within 5 minutes through a distributed invalidation list.

## API Key Authentication

Third-party integrations use API keys for authentication. API keys are passed in the X-API-Key header. Keys are long-lived but can be rotated through the developer dashboard. Each key is scoped to specific permissions.
