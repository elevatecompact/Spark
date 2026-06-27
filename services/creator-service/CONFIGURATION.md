# creator-service — Configuration

## Environment Variables
| Variable | Default | Description |
|----------|---------|-------------|
| CREATOR_PORT | 4002 | HTTP listener |
| CREATOR_DB_URL | — | PostgreSQL connection (required) |
| CREATOR_REDIS_URL | — | Redis connection (required) |
| MAX_TIERS_PER_CHANNEL | 10 | Subscription tier limit |
| VERIFICATION_REQUIRED | true | Must verify to monetize |
| PAYOUT_MINIMUM_CENTS | 5000 | Minimum  for payout |
| CDN_BASE_URL | https://cdn.titan.dev | Base URL for media |
| DOCUMENT_STORAGE_BUCKET | titan-verification-docs | S3 bucket for ID docs |

## Feature Flags
| Flag | Default | Purpose |
|------|---------|---------|
| self_verification | true | Allow creators to self-verify via docs |
| auto_approve_verification | false | Skip manual review for trusted creators |
| tier_pricing_frozen | false | Lock tier pricing changes globally |

Rate limits: Channel creation 1/hour per creator. Tier updates 10/hour. Verification 3/day.
