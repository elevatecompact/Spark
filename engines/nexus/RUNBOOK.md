# Runbook

## Startup
1. Verify S3: ./nexus test-s3.
2. Run migrations: ./nexus migrate.
3. Start API: ./nexus serve.
4. Start workers: ./nexus worker --queue redis://localhost:6379.
5. Test: curl -F "file=@test.jpg" http://localhost:8080/v1/upload.

## Monitoring
- Dashboard: upload throughput, storage by tier, queue depth.
- Alerts: S3 errors > 1%, queue backlog > 10K, storage > 85%.
