# identity-service — README

## Overview
The Identity Service is the central authentication and authorization authority for the Titan platform. It manages user identities, credentials, session tokens, OAuth flows (Google, GitHub, Discord), multi-factor authentication via TOTP, and role-based access control (RBAC). Every request entering the platform is validated against this service to establish trust and identity context before being proxied to downstream services.

## Purpose
Provide a unified identity layer that supports email/password authentication, social login, MFA enrollment, and API key management for programmatic access. Maintains the user profile root record from which all other services derive identity data. Handles password hashing with bcrypt, JWT issuance with RS256 signing, refresh token rotation, and session revocation. Integrates with the moderation service for account suspension enforcement and with the notification service for auth-related alerts.

## Ownership
**Team:** Platform Security (eng-security@titan.dev)
**SLI:** 99.99% uptime, p99 auth response < 200ms
**Escalation:** #oncall-identity in Slack
