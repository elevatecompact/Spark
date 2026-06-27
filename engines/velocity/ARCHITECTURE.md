# Architecture

Velocity uses provider-abstracted architecture. Core defines CDNProvider interface implemented by adapters for Fastly, CloudFront, Akamai, Cloudflare. Orchestrator handles cache warming by prefetching popular content to edge nodes. Traffic steering module uses real-time latency and availability data to route viewers to optimal CDN. Purging layer supports instant, tag-based, and wildcard purges. Redis-backed config store propagates routing rules to edge Envoy proxies.
