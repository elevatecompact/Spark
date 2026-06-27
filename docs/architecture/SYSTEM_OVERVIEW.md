# System Overview

Spark is a globally distributed live streaming platform designed to connect performers and audiences worldwide in real time. The system architecture is built on microservices principles, event-driven communication, and multi-cloud infrastructure to ensure reliability, scalability, and low-latency interactions.

## High-Level Architecture

The platform is organized into five primary layers:

### Presentation Layer
Client applications across web, mobile, and smart TV platforms connect through an Envoy-based API gateway. The gateway handles TLS termination, rate limiting, authentication, and request routing to internal services. WebRTC signaling and media streams traverse dedicated pathways optimized for real-time communication.

### Application Layer
A set of loosely coupled microservices implements the core business logic. Services are organized by domain boundary including identity, streaming, discovery, monetization, moderation, and analytics. Each service owns its data store and communicates asynchronously via Kafka event streams or synchronously via gRPC for low-latency queries.

### Streaming Layer
The streaming pipeline ingests video from broadcasters through RTMP or WebRTC, transcodes to multiple bitrates, and delivers via a global CDN powered by Cloudflare and the proprietary Velocity edge engine. Real-time messaging uses WebRTC data channels and a global signaling mesh.

### Data Layer
Persistent storage uses S3-compatible object storage backed by the Nexus engine for high-throughput media assets. Relational data lives in horizontally sharded PostgreSQL clusters. Caching leverages a multi-tier strategy with Redis at the application layer, CDN edge caches, and browser-level caching.

### Infrastructure Layer
The platform runs across AWS and GCP in multiple regions with active-active failover. Kubernetes orchestrates containerized services. Infrastructure is managed as code using Terraform with environment-specific configurations.

## Key Design Principles

- **Resilience by Default**: Every component is designed for failure. Circuit breakers, retries, and bulkheads prevent cascading outages.
- **Global by Design**: All services are region-aware and route traffic to the nearest healthy endpoint.
- **Observability**: Distributed tracing via OpenTelemetry, structured logging, and Prometheus metrics provide full system visibility.
