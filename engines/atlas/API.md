# API

## Service Registration
- PUT /v1/register - Register a service instance with metadata.
- DELETE /v1/register/:serviceId - Deregister a service instance.
- POST /v1/register/:serviceId/ttl - Heartbeat TTL renewal.

## Service Discovery
- GET /v1/discover/:serviceName - Return healthy endpoints for a service.
- GET /v1/discover/:serviceName/:tag - Filter by metadata tag.
- GET /v1/watch/:serviceName - Long-poll SSE stream for topology changes.

## Health
- GET /v1/health - Agent health status.
- GET /v1/health/:serviceId - Detailed health of a specific instance.
