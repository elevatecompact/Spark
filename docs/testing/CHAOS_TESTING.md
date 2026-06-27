# Chaos Testing

## Philosophy

Titan embraces chaos engineering to build resilience into the platform. We proactively inject failures into production-like environments to verify that the system degrades gracefully and recovers automatically.

## Tooling

Chaos Mesh provides Kubernetes-native chaos engineering for pod failures, network latency, network partition, CPU stress, and disk I/O faults. Litmus handles more complex scenarios involving infrastructure such as RDS failover and AZ outage simulation. Custom fault injection is applied through application-layer faults like delayed responses, corrupted payloads, and authentication failures via Chaos Mesh HTTPChaos and custom sidecars.

## Experiment Catalog

### Network Faults
Network delay experiments add 50ms, 200ms, and 1000ms between specific service pairs. Network partition experiments isolate a service from its database or cache. Packet loss experiments test retry logic at 1%, 5%, and 10% loss rates.

### Infrastructure Faults
Pod kill experiments test graceful shutdown, leader election, and connection draining. Node failure experiments simulate EC2 instance failure validating pod rescheduling and data durability. AZ outage experiments test cross-region failover.

### Application Faults
Dependency timeout experiments simulate external API timeouts. Rate limit exceeded tests verify upstream 429 handling. Database connection pool exhaustion tests validate connection management.

## Schedule

GameDays are quarterly full-day events where SRE and engineering teams run a curated experiment list. Continuous chaos runs a subset of low-risk experiments automatically in staging every night. Release chaos runs the full chaos suite against every major release candidate.