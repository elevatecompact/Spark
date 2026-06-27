# Security

- Row-level security via ALLOWED_ACCOUNT policies.
- Query sandbox - DDL and ALTER TABLE blocked.
- Input validated against Avro schema; malformed events rejected.
- Alert webhooks support HMAC signatures.
- GDPR right-to-deletion via per-row deletion policies.
- ClickHouse data encrypted at rest; TLS for inter-node and Kafka.
