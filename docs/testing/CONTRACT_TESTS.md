# Contract Testing

## Consumer-Driven Contracts

Titan uses Pact for consumer-driven contract testing (CDC). Each service defines the contracts it expects from its upstream dependencies. These contracts are published to a Pact Broker and verified by the provider during CI.

## Workflow

The consumer (downstream service) writes a Pact test defining expected request/response pairs. The Pact file is published to the Pact Broker during consumer CI. The provider (upstream service) fetches all active contracts from the broker. Provider CI runs Pact verification tests against its actual API. The result is published back to the broker — failures block the provider's deployment.

## Contract Scope

Contracts cover HTTP APIs including path, method, headers, query params, request body, response status, and response body. They also cover event-driven message queue payload schemas (RabbitMQ, Kafka) and gRPC protobuf message schemas with RPC signatures.

## Benefits

Early detection of provider changes that break consumers before deployment. Independent deployments without team coordination if contracts remain valid. The Pact Broker serves as a living catalog of all inter-service contracts.

## Tooling

Pact is used with pact-go, pact-python, and pact-js language bindings. PactFlow hosted broker manages contract storage. Pact CLI integrates with GitHub Actions. Webhooks notify consumer teams when a provider changes a contract.

## Best Practices

Contracts are versioned and a provider may support multiple contract versions simultaneously. WIP pacts allow experimental endpoints without blocking CI. Deleted or deprecated contracts are cleaned up via broker webhooks. Contract test suites should complete in under 30 seconds per service pair.