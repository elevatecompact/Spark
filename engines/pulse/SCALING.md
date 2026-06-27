# Scaling

Pulse scales horizontally by adding SFU nodes behind a UDP-aware load balancer. The control plane uses a Redis-backed room registry that maps stream IDs to SFU node addresses. Ingest gateways (RTMP/SRT acceptors) are stateless and can be scaled independently. Each SFU node registers itself on startup and sends heartbeat pings to the registry. When a viewer requests a stream, the signalling layer queries the registry, assigns the least-loaded SFU, and returns its address. For global distribution, deploy SFU nodes in edge PoPs with Anycast IPs.
