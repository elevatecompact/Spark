# viewer-service — Event Contracts

## Published Events
The viewer-service publishes domain events to Kafka for asynchronous consumption.

| Topic | Description |
|-------|-------------|
| $(System.Collections.Hashtable.slug).resource.created | Emitted when a new resource is created |
| $(System.Collections.Hashtable.slug).resource.updated | Emitted when a resource is modified |
| $(System.Collections.Hashtable.slug).resource.deleted | Emitted when a resource is removed |

## Consumed Events
| Topic | Source | Handler |
|-------|--------|---------|
| identity.user.deleted | identity-service | Cascade delete user data |
| identity.user.updated | identity-service | Update cached references |
| 
otification.delivery.confirmed | notification-service | Delivery receipts |

## Schema Format
All events use Avro serialization with Confluent Schema Registry.
`json
{
  "eventId": "uuid",
  "source": "viewer-service",
  "type": "resource.created",
  "timestamp": "2026-06-27T12:00:00Z",
  "payload": {},
  "traceId": "uuid"
}
`
Retention: 7 days on Kafka. Compacted topics for keyed events.
