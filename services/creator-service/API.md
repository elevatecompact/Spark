# creator-service — API Contract

## Channels
- POST /v1/channels — Create channel with name, category, description
- GET /v1/channels/{id} — Full channel profile with metrics
- PATCH /v1/channels/{id} — Update branding, category, description
- DELETE /v1/channels/{id} — Archive channel (30-day recovery window)
- POST /v1/channels/{id}/verify — Submit identity verification documents

## Subscription Tiers
- POST /v1/channels/{id}/tiers — Create tier (name, price_cents, benefits[])
- GET /v1/channels/{id}/tiers — List active tiers ordered by price
- PATCH /v1/tiers/{id} — Update pricing or benefits
- DELETE /v1/tiers/{id} — Archive tier (existing subs unaffected)

## Payouts
- GET /v1/creators/{id}/payouts — Payout history with status
- POST /v1/creators/{id}/payout-preferences — Set method and minimum threshold
- GET /v1/creators/{id}/revenue — Revenue breakdown by source

## Verification
- GET /v1/verification/status — Current verification level and next steps
- POST /v1/verification/documents — Upload ID, tax forms, or business docs

All responses standard envelope. Cursor pagination on list endpoints.
