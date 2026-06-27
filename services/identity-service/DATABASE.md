# identity-service — Database Schema

## PostgreSQL — Primary Store
### users table
| Column | Type | Constraints |
|--------|------|-------------|
| id | UUID | PK, DEFAULT gen_random_uuid() |
| email | VARCHAR(320) | UNIQUE, NOT NULL |
| password_hash | VARCHAR(255) | Nullable for OAuth-only accounts |
| display_name | VARCHAR(100) | NOT NULL |
| avatar_url | TEXT | Nullable |
| mfa_secret | VARCHAR(64) | Encrypted at rest, Nullable |
| mfa_enabled | BOOLEAN | DEFAULT false |
| status | ENUM('active','suspended','deleted') | DEFAULT 'active' |
| role | ENUM('user','moderator','admin') | DEFAULT 'user' |
| created_at | TIMESTAMPTZ | NOT NULL |
| updated_at | TIMESTAMPTZ | NOT NULL |
| deleted_at | TIMESTAMPTZ | Nullable, soft delete |

### sessions table
Stores refresh token hashes. Active sessions tracked for forced logout capability.
| Column | Type | Notes |
|--------|------|-------|
| id | UUID | PK |
| user_id | UUID | FK -> users.id, indexed |
| refresh_token_hash | VARCHAR(255) | SHA-256 hash |
| device_info | JSONB | User agent, IP, platform |
| expires_at | TIMESTAMPTZ | TTL: 7 days |
| revoked_at | TIMESTAMPTZ | Nullable |

### pi_keys table
Scoped API keys with last-used tracking. Only SHA-256 prefix stored for security.

## Redis — Session Cache & Rate Limiter
Session TTL: 15 minutes (access tokens). Rate limit counters: sliding window per IP/endpoint. MFA challenge tokens: TTL 5 minutes.

Partition users by created_at quarterly. Indexes on email (UNIQUE), status, role.
