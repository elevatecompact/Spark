# Configuration

| Key | Default | Description |
|-----|---------|-------------|
| clickhouse.host | localhost | ClickHouse host |
| clickhouse.port | 9000 | Native protocol port |
| clickhouse.database | titan | Default database |
| clickhouse.compression | lz4 | Compression type |
| kafka.brokers | localhost:9092 | Kafka broker list |
| kafka.batch_size | 10000 | Max events per batch |
| retention.default_days | 30 | Data retention period |
| retention.hot_days | 7 | Hot data retention |
| query.max_result_bytes | 10485760 | Max result size (10MB) |
| alert.evaluation_interval | 60 | Eval interval in seconds |
