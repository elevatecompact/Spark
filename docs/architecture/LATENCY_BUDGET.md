# Latency Budget and SLOs

Spark defines explicit latency budgets for all critical user journeys. These budgets guide architectural decisions, drive performance optimization, and form the basis of service level objectives (SLOs) and service level indicators (SLIs).

## User Journeys and Budgets

### Live Stream Viewing (End-to-End)
Broadcaster Camera to Viewer Screen: 3,000ms total budget

| Segment | Allocation | Description |
|---------|-----------|-------------|
| Camera to Edge Ingest | 200ms | FEC + network transmission |
| Ingest to Origin | 100ms | Regional fiber transport |
| Transcoding | 500ms | GPU encode + segment packaging |
| Origin to CDN Edge | 200ms | Cache fill propagation |
| CDN Edge to Player | 1,000ms | Segment download |
| Player Buffer | 500ms | Smooth playback buffer |
| Decode + Render | 500ms | Client-side processing |

### Chat Message Delivery
Sender to Receiver: 200ms total budget

| Segment | Allocation |
|---------|-----------|
| Sender to SFU | 50ms |
| SFU Processing | 20ms |
| SFU to Receiver | 130ms |

### API Request (Read)
Client to Response: 200ms p99 budget

| Segment | Allocation |
|---------|-----------|
| Network (client to edge) | 30ms |
| Cloudflare processing | 5ms |
| Envoy routing + auth | 15ms |
| Application logic | 50ms |
| Database query | 80ms |
| Response encoding | 20ms |

### API Request (Write)
Client to Async Confirmation: 500ms p99 budget. Includes network (30ms), gateway + auth (20ms), validation (50ms), event persist (200ms Kafka ack), and response handling (200ms).

## SLOs by Service

| Service | p50 | p95 | p99 | Target |
|---------|-----|-----|-----|--------|
| API Gateway | 20ms | 80ms | 150ms | 99.9% under 200ms |
| Stream Service | 50ms | 200ms | 500ms | 99.9% under 500ms |
| Search Service | 30ms | 100ms | 200ms | 99.5% under 300ms |
| CDN Delivery | 100ms | 800ms | 2,000ms | 99% under 2s |
| Transcoding | 300ms | 1,000ms | 3,000ms | 99% under 3s |

## Error Budget Policy

Each service has a monthly error budget of 0.1% (for 99.9% SLO). When burn rate exceeds 10x the budget consumption rate, alerts fire and feature deployments are gated until the budget recovers. Budget alerts page the on-call team through PagerDuty with severity based on projected budget exhaustion time.
