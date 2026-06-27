# API

## Ingestion
- POST /v1/ingest/event - Ingest single event.
- POST /v1/ingest/batch - Batch ingest with Avro payloads.

## Query
- POST /v1/query - Execute read-only SQL against ClickHouse.
- GET /v1/query/saved - List saved queries.

## Dashboard
- POST /v1/dashboard - Create dashboard.
- GET /v1/dashboard/:id - Render dashboard.

## Alert
- POST /v1/alert - Create alert rule.
- GET /v1/alert - List alerts with status.
- POST /v1/alert/:id/ack - Acknowledge alert.
