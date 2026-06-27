# Schema Design Principles

This document defines the schema design principles and conventions used across SPARK's database systems.

## General Principles

**Normalization.** Relational schemas follow third normal form (3NF) by default. Denormalization is permitted only when justified by measured query performance and documented with the rationale.

**Consistent Naming.** All table, column, index, and constraint names use snake_case. Table names are plural nouns (users, content, subscriptions). Column names are singular nouns (user_id, content_type, created_at). Foreign key columns match the referenced table name in singular form followed by _id.

**Primary Keys.** All tables use UUIDv4 primary keys stored as UUID type. Auto-increment integers are avoided for primary keys to prevent enumeration attacks and simplify migration between systems.

**Timestamps.** Every table includes created_at (NOT NULL, default NOW()) and updated_at (NOT NULL, default NOW()). The updated_at column is automatically updated by a trigger. Soft deletes use a deleted_at TIMESTAMPTZ column, NULL indicates active records.

**Immutable Audit Log.** Mutations to sensitive tables are recorded in immutable audit log tables. Each audit entry captures the old and new values, the acting user, timestamp, and the source IP address.

## Type Conventions

Monetary values use NUMERIC(12,4) with currency stored as a ISO 4217 code. JSONB is used for flexible metadata with GIN indexes for queries. Enumerated types use domain types or lookup tables, never raw string columns.
