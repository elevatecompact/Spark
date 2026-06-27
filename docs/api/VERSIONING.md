# API Versioning

SPARK maintains backward compatibility for all APIs through a versioning strategy that allows the platform to evolve without breaking existing clients.

## Versioning Strategy

The API uses a URL-based versioning scheme where the major version is included in the URL path: /v1/users, /v2/users. Breaking changes require a new major version. Minor and patch changes are additive and do not require a version bump. This approach makes version explicit in every request and easy to route at the API gateway level.

## Version Lifecycle

Each API version goes through three phases. Active versions receive full support including new features, bug fixes, and performance improvements. Deprecated versions are announced to developers with a minimum 90-day notice before the deprecation date. Deprecated versions receive only critical security and bug fixes. Sunset versions are no longer available. Deprecated versions include a Sunset header in responses indicating the planned removal date.

## Breaking Changes

Changes that require a new major version include removing fields from responses, changing field types, making previously optional fields required, changing endpoint URLs, removing endpoints, and changing error codes. Changes that do not require a new version include adding new endpoints, adding new optional fields to responses, adding new optional request parameters, and extending error code sets.

## Client Communication

Version deprecations are communicated through multiple channels including email notifications to registered developers, the developer dashboard banner, and the Sunset HTTP header in API responses. Migration guides are provided for each version transition, with code examples showing how to update from the old to the new version.

## Internal Versioning

Internal service-to-service gRPC APIs follow semantic versioning on the proto package declaration. Breaking proto changes require coordination across all consuming services.
