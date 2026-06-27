# API

## Content Management
- PUT /v1/cache/warm - Request cache warming for URLs.
- GET /v1/cache/status/:url - Check cache status across CDNs.
- POST /v1/cache/purge - Purge URLs by pattern, tag, or exact match.
- POST /v1/cache/purge/all - Full cache invalidation.

## Traffic Steering
- POST /v1/steering/rules - Create steering rules.
- GET /v1/steering/rules - List active rules.
- GET /v1/steering/status - Current traffic distribution.

## Provider Management
- GET /v1/providers - List configured CDN providers and health.
- POST /v1/providers/:name/failover - Manual failover.
