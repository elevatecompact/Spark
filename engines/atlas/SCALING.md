# Scaling

Atlas scales horizontally by adding registrar nodes. The etcd cluster handles consensus for up to 7 registrar nodes. Agents are stateless and auto-discover registrars via the first healthy registrar response. For multi-datacenter deployments, each DC runs its own Atlas ring with a global etcd watcher that cross-replicates critical services. Agent cache TTL reduces registrar load; under failure, agents fall back to direct etcd reads. Use locality-preference routing to keep traffic within the same DC.
