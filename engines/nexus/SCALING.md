# Scaling

Nexus scales with API nodes behind load balancer. S3 handles storage scaling transparently. Transformation workers scale via Redis job queue. High upload throughput uses S3 Transfer Acceleration with presigned URLs. Metadata reads offloaded to Redis replicas. PostgreSQL read replicas for listing and search. Cold tier migrations batched and throttled for S3 API limits.
