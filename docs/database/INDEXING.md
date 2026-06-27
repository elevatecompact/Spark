# Indexing Strategy

Indexing is critical to SPARK's query performance. This document outlines index design principles, types, and maintenance procedures.

## Index Design Principles

Indexes are designed based on measured query patterns, not speculation. Every index must correspond to a specific query or set of queries identified through application profiling. Indexes are reviewed quarterly for usage and redundancy through the pg_stat_user_indexes view to identify unused indexes for removal.

## PostgreSQL Index Types

**B-tree indexes** are the default for equality and range queries on columns with high cardinality. Composite B-tree indexes follow the leftmost prefix rule, with the most selective column first.

**GIN indexes** are used for JSONB queries, full-text search columns, and array columns. They efficiently handle queries checking containment (?) and key existence (?|, ?&).

**GiST indexes** support exclusion constraints and proximity searches, used primarily for geospatial queries.

**BRIN indexes** are used on very large, append-only tables where the data has natural ordering by timestamp. They provide space-efficient indexing for time-series data on the created_at column.

## Partial and Covering Indexes

Partial indexes with WHERE clauses index only relevant rows, reducing index size and write overhead. Covering indexes with INCLUDE columns allow index-only scans by including payload columns in the index without affecting the sort order.

## Index Maintenance

The autovacuum daemon maintains index health. REINDEX is run during maintenance windows for indexes that show significant bloat. Concurrent index creation is used in production to avoid locking writes.
