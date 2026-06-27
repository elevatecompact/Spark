# Troubleshooting

## Slow queries
1. Check if hitting materialised view vs raw table.
2. Verify partition pruning by date or tenantId.
3. Use EXPLAIN SELECT to check index usage.
4. Check query.max_result_bytes - large results spill to disk.
5. Review ClickHouse merge backlog.

## Missing data
1. Check Kafka consumer lag.
2. Verify Avro schema matches incoming events.
3. Check ClickHouse insert logs for rejected rows.
4. Verify retention.default_days - old data auto-dropped.
