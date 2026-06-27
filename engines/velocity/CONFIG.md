# Configuration

| Key | Default | Description |
|-----|---------|-------------|
| providers | [] | List of CDN provider configs |
| warming.concurrency | 10 | Max concurrent warming requests |
| warming.queue_size | 100000 | Max queued warm requests |
| purge.max_batch_size | 1000 | Max URLs per purge request |
| steering.probe_interval | 60 | Edge probe interval (seconds) |
| steering.failover_threshold | 0.05 | Failover when error rate > 5% |
| steering.default_provider | fastly | Default CDN provider |
| origin.shield.enabled | true | Enable origin shield |
