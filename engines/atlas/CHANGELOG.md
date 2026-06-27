# Changelog

## 1.2.0 (2026-06-10)
- Added Consul API compatibility layer for legacy service migration.
- Introduced locality-preference routing strategy.
- Reduced agent memory by 40% through shared cache buffer.
- Added blue-green deployment weight override endpoint.

## 1.1.0 (2026-04-05)
- Switched to bidirectional gRPC streams for heartbeat delivery.
- Added SSE watch endpoint for real-time topology changes.
- etcd compaction job integrated into registrar lifecycle.

## 1.0.0 (2026-01-20)
- First stable release with etcd backend.
- Service registration, discovery, and health checking.
