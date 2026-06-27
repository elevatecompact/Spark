# Events

## Published Events
- service.registered - Payload: { serviceId, serviceName, endpoints, tags, ttl }.
- service.deregistered - Payload: { serviceId, reason }.
- service.health.changed - Payload: { serviceId, previousStatus, currentStatus }.
- service.topology.changed - Full topology diff snapshot emitted every 30s or on change.
- agent.connected / agent.disconnected - Agent lifecycle events.

## Subscribed Events
- discovery.force.refresh - Trigger immediate cache invalidation across all agents.
- discovery.routing.update - Update routing weights for a service.
