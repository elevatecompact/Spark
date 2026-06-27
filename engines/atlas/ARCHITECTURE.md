# Architecture

Atlas runs as a sidecar agent on every Titan node plus a small cluster of registrar nodes. The agent performs health checks and reports status to the registrar cluster via gRPC bidirectional streams. Registrars maintain an in-memory service graph backed by etcd for durability. Client lookups hit the local agent first (cache), then fall back to the registrar cluster. The routing layer supports weighted random, consistent hash, and locality-preference strategies. Every registration includes metadata tags for environment, version, and canary status.
