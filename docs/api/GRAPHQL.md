# GraphQL API

GraphQL is SPARK's primary API protocol, serving as the main interface for web and mobile applications. It provides a strongly-typed, declarative data fetching model that allows clients to request exactly the data they need.

## Schema Design

The GraphQL schema is organized around the concept of nodes and edges, following the Relay specification. Every entity implements the Node interface with a globally unique ID. Connections implement the Connection interface with edges, pageInfo, and totalCount fields.

## Query Patterns

The root Query type exposes fields for fetching single entities by ID (user, content, subscription) and lists via connections (users, contents, subscriptions). Filtering is provided through arguments on connection fields using input types such as UserFilter, ContentFilter, and DateRangeFilter.

## Mutations

Mutations follow the Relay input pattern: each mutation accepts a single input argument and returns a payload type. Mutations are named in the imperative mood (createUser, updateContent, deleteSubscription). Optimistic updates are supported through the @live query directive for real-time data.

## Performance

The DataLoader pattern batches and caches database queries within a single request, preventing N+1 query problems. Query complexity analysis limits the maximum query depth to 8 levels and caps the number of requested nodes per connection to 100. Persisted queries are used for high-traffic operations to reduce parsing overhead.

## Federation

The GraphQL gateway federates schemas from multiple microservices using Apollo Federation. Each service declares its types and extends types from other services through the @key directive.
