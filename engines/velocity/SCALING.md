# Scaling

Velocity scales by adding orchestrator nodes behind load balancer. Warming and purge workers are independent goroutine pools per node. Steering decisions made at edge via Envoy sidecars fetching routing tables from Redis. For high-scale purging, Kafka queue decouples requests from execution. Provider API rate limits managed via per-provider token buckets. Origin shield scales based on cache miss rate.
