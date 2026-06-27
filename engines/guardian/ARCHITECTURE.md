# Architecture

Guardian uses three tiers: API gateways, auth core, and token cache. API gateway terminates TLS, handles OAuth redirect flows, proxies to auth core. Auth core runs OIDC flows, validates credentials, issues Ed25519-signed JWTs. Redis cache stores revoked token IDs and session metadata for fast validation without hitting PostgreSQL. RBAC engine loads policy definitions from PostgreSQL and evaluates requests against user roles and resource permissions.
