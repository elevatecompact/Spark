# Audit Logging

Spark captures comprehensive audit logs for all security-relevant events. The audit logging system is designed for integrity, tamper resistance, and efficient querying.

## Logged Events

| Category | Events |
|---|---|
| Authentication | Login, logout, MFA enrollment, password change, MFA challenge |
| Authorization | Permission grants, role changes, access denials, policy violations |
| Data Access | Read/write/delete operations on sensitive data, data exports |
| Administrative | Configuration changes, user management, system modifications, privilege escalation |
| Security | Password changes, API key creation/rotation, group membership changes |
| System | Service starts/stops, deployments, certificate rotations |
| Compliance | Data retention actions, privacy requests, consent changes |

## Log Format and Transport

Logs are structured JSON events with the following required fields: `event_id`, `timestamp`, `actor_id`, `actor_type`, `action`, `resource_type`, `resource_id`, `outcome`, `source_ip`, `user_agent`, `correlation_id`. Logs are sent to a centralized log aggregation service via TLS-encrypted transport with authentication.

## Integrity and Tamper Protection

Audit logs are written to an append-only, cryptographically chained log store (immutable ledger). Log entries are hashed and linked to the previous entry, creating a hash chain that prevents undetected modification. Log storage is append-only at the filesystem level.

## Retention

| Log Category | Retention Period |
|---|---|
| Authentication events | 90 days (hot), 2 years (cold archive) |
| Authorization events | 90 days (hot), 2 years (cold archive) |
| Data access events | 1 year (hot), 5 years (cold archive) |
| Administrative events | 2 years (hot), 7 years (cold archive) |

## Monitoring and Alerting

Audit logs are monitored in real time by the SIEM. Alert rules detect mass privilege changes, unusual data access patterns, failed authentication cascades, and configuration drift.
