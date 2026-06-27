# DDoS Protection

Spark employs a layered defense against distributed denial-of-service (DDoS) attacks at the network, application, and infrastructure levels.

## Network Layer Defense

- **Anycast routing** — Traffic is distributed across multiple globally distributed points of presence (POPs). An attack against one POP is absorbed without affecting others.
- **Volumetric filtering** — At the network edge, traffic is filtered based on source IP reputation, protocol anomalies, and connection rates. Scrubbing centers remove malicious traffic before it reaches origin infrastructure.
- **Bandwidth capacity** — The edge network maintains 10x normal peak traffic capacity to absorb large-scale volumetric attacks.

## Application Layer Defense

- **Web Application Firewall (WAF)** — The WAF inspects HTTP/HTTPS traffic for application-layer attacks including HTTP floods, Slowloris, and protocol violations. WAF rules are tuned based on observed attack patterns.
- **Rate limiting** — Per-IP and per-user rate limits prevent any single source from overwhelming API endpoints. Limits are dynamically adjusted based on traffic patterns.
- **Challenge-based mitigation** — Suspicious clients receive JavaScript challenges, CAPTCHA, or proof-of-work challenges before reaching application servers.

## Infrastructure Hardening

- **Auto-scaling** — All services are configured for horizontal auto-scaling. During a DDoS event, additional instances are provisioned to maintain availability.
- **Connection limits** — Per-server connection limits prevent resource exhaustion. Idle connections are terminated after a configurable timeout.
- **Caching** — Static and cacheable dynamic content is served from edge caches, reducing load on origin servers.

## Monitoring and Response

The security operations center (SOC) monitors DDoS alerts in real time. Automated mitigation playbooks trigger on attack detection. Post-incident analysis identifies attack vectors and updates defenses.
