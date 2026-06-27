# Database Coding Standards

This document defines mandatory coding standards for all database development at SPARK.

## SQL Style

SQL keywords use uppercase (SELECT, INSERT, CREATE). Identifiers use snake_case. Queries are formatted with each major clause on a new line with consistent indentation. Subqueries are indented and aliased clearly. JOIN conditions specify the join type explicitly (INNER, LEFT, RIGHT) and conditions are in the ON clause, not the WHERE clause.

## Naming Conventions

Tables are plural nouns (users, content_items, subscription_plans). Columns are singular nouns (username, display_name, created_at). Boolean columns use affirmative names with is_ prefix (is_active, is_verified, is_deleted). Timestamp columns end with _at (created_at, updated_at, deleted_at). Date columns end with _date. Foreign key columns match the referenced table name in singular form with _id suffix.

## Migration Standards

Every migration must be reversible. Migrations must be backward compatible with the previous application version. No migration may rename a column without a two-phase process. Data migrations run as separate scripts after schema changes. Each migration script includes a header comment with the author, date, and purpose.

## Performance Standards

All queries in application code must use parameterized statements. EXPLAIN ANALYZE is required before shipping any new query to production. Indexes must be justified by query patterns. N+1 query patterns are detected through application tracing and rejected in code review.
