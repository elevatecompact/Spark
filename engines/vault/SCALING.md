# Scaling

Vault scales with stateless API nodes behind load balancer. Idempotency keys in Redis ensure safe retries across nodes. Webhooks processed via Kafka consumer group for parallel processing. Subscription billing triggered by scheduler service (single leader via Redis lock) creating invoices for due subscriptions. Dedicated payment processing worker pools prevent billing backlog. PostgreSQL read replicas for listing and analytics.
