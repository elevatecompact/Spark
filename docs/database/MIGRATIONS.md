# Migration Process

Database migrations at SPARK follow a rigorous, automated process designed to ensure zero-downtime deployments and rollback capability.

## Migration Tooling

All schema changes are managed through Flyway for PostgreSQL. Migration files are SQL scripts named following the convention V{version}__{description}.sql. Versions are sequential integers. Each migration is immutable once applied to any environment — corrections are made through new migration scripts.

## Migration Workflow

Developers write migration scripts in their feature branches. The migration must be backward-compatible: new columns must be nullable or have defaults, new tables must not have NOT NULL constraints referencing existing data, and column renames must use a two-phase process (add new column, migrate data, drop old column).

CI runs migrations against a fresh database to verify idempotency and correctness. The staging environment applies migrations automatically during deployment. Production migrations run through a peer-reviewed change request process with automated execution during the maintenance window.

## Zero-Downtime Patterns

For large tables, migrations use batching to avoid long-running locks. CREATE INDEX CONCURRENTLY adds indexes without blocking writes. Adding foreign keys uses NOT VALID followed by VALIDATE CONSTRAINT. Column type changes use a shadow column approach with triggers to keep both columns synchronized during the transition.

## Rollback

Each migration includes a corresponding rollback script. Rollback is automated for the most recent migration but requires manual approval for earlier migrations. The Flyway repair command is available for correcting failed migrations.
