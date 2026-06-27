# Abuse Prevention Strategies

Spark implements a multi-layered abuse prevention system to detect and mitigate platform misuse including account takeovers, spam, scraping, API abuse, and fraudulent activity.

## Rate Limiting

All API endpoints are rate-limited using a sliding window algorithm. Limits are applied per user, per IP, and per tenant. Endpoints handling sensitive operations (authentication, password reset, data export) have stricter limits. Excess requests receive HTTP 429 responses with Retry-After headers.

## Bot Detection

A combination of techniques is used to identify automated traffic:

- **JavaScript challenges** — Browser-based proof-of-work challenges for suspicious requests
- **Behavioral analysis** — Mouse movement, keystroke dynamics, and navigation patterns
- **Signature analysis** — TLS fingerprinting, HTTP header order analysis, and browser automation detection
- **CAPTCHA** — reCAPTCHA v3 is triggered for high-risk interactions

## Content Moderation

User-generated content is scanned for abuse using automated classifiers and hash-matching against known abusive content databases. Reports can be submitted by users and are triaged by the trust and safety team.

## Account Protections

- **New account limits** — New accounts have reduced rate limits for the first 48 hours
- **Suspicious registration detection** — Email addresses from known burner domains, disposable phone numbers, and previously flagged devices are blocked
- **Login velocity checks** — Multiple failed logins from disparate geographies trigger account lockout

## API Abuse

API keys are scoped to specific operations and resources. Anomalous API usage patterns (unusual query patterns, unexpected data volumes, rapid pagination) are flagged and escalate to automated or manual review.
