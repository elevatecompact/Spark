# stream-service — Runbook
## Alerts: StreamE2ELatency > 10s, TranscoderQueueDepth > 200, IngestDisconnectRate > 5%, RecordingFailureRate > 1%
## Force end: POST /v1/admin/streams/{id}/force-stop
## Reprocess recording: POST /v1/admin/recordings/{id}/reprocess
## Scale transcoders: kubectl scale deployment/stream-transcoder --replicas=20
## Incident: Confirm alert → identify component → mitigate → postmortem
