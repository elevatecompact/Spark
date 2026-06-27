# Scaling

Oracle scales horizontally across inference nodes, feature store shards, and vector search replicas. Inference nodes are stateless behind gRPC round-robin load balancer. Redis cluster shards by user ID hash. Milvus replication factor of 3 ensures availability. Kafka partitions by item ID with independent consumer groups. The ranker caches business rules and uses hierarchical ranking to avoid O(n) scoring.
