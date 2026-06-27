# Rate Limiting

SPARK implements rate limiting across all API endpoints to ensure fair resource usage, protect against abuse, and maintain platform stability.

## Rate Limit Policies

Unauthenticated requests are limited to 60 requests per minute per IP address. Authenticated requests are limited to 1,000 requests per minute per user. Individual endpoint limits vary based on resource cost: read endpoints allow 5,000 requests per minute, write endpoints allow 500 requests per minute, and expensive query endpoints allow 100 requests per minute.

## Rate Limit Headers

Every API response includes three rate limit headers: X-RateLimit-Limit shows the maximum number of requests allowed in the current window, X-RateLimit-Remaining shows the number of requests remaining, and X-RateLimit-Reset shows the Unix timestamp when the rate limit window resets. When a rate limit is exceeded, the API returns a 429 Too Many Requests status code with a Retry-After header indicating the number of seconds to wait before retrying.

## Implementation

Rate limiting uses a sliding window algorithm implemented in Redis. Each window is tracked as a sorted set with timestamps as scores. Expired entries are removed on each request to keep the data structure compact. The sliding window approach provides smoother rate enforcement compared to fixed windows which can allow burst traffic at window boundaries.

## Exemptions

Internal service-to-service communication through the service mesh is exempt from rate limits. Approved third-party partners may request higher limits through the developer program. Rate limit configurations are managed through a central configuration service and can be updated without deployments.
