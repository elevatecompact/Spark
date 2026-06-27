# Scaling

Guardian scales behind stateless HTTP load balancer. All nodes share PostgreSQL for persistence and Redis for cache and token blacklist. Token validation is the hottest path - Redis lookup with local LRU cache for extreme throughput. OAuth flows use Redis-stored state for stickiness. Multi-region deployments use Redis CRDTs for cross-region token blacklist replication. PostgreSQL read replicas handle audit queries.
