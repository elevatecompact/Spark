# Atlas Engine

**Purpose:** Service discovery and routing engine for the Titan microservices ecosystem.
**Tech Stack:** Go, etcd, gRPC, mDNS, Consul API compatibility layer.

Atlas provides service registration, health-based discovery, and dynamic routing across all Titan engines. It supports multiple backends (etcd, Consul) via a unified interface, enables blue-green deployments through weight-based routing, and maintains a real-time topology graph of the entire service mesh.
