# Events

## Published Events
- recommendation.served - Payload: { userId, requestId, items, context }.
- recommendation.clicked - Payload: { userId, itemId, requestId, position, timestamp }.
- model.deployed - Payload: { modelId, modelVersion, metrics, deployedAt }.
- model.stale - Alert when online metrics drift exceeds threshold.

## Subscribed Events
- user.activity - Consumed from Kafka: { userId, itemId, eventType, timestamp }.
- catalog.update - Content metadata refresh to rebuild embeddings.
- abtest.assignment - User assigned to experiment variant.
