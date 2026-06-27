# creator-service — Event Contracts

## Published Events
| Topic | Event | Description |
|-------|-------|-------------|
| creator.channel.created | ChannelCreatedEvent | New channel goes live on platform |
| creator.channel.verified | ChannelVerifiedEvent | Identity verification approved |
| creator.channel.suspended | ChannelSuspendedEvent | ToS violation suspension |
| creator.tier.created | TierCreatedEvent | New subscription tier configured |
| creator.tier.pricing_changed | TierPricingChangedEvent | Price adjustment on existing tier |
| creator.payout.preferences_set | PayoutPreferencesSetEvent | Payout method configured |

## Consumed Events
| Topic | Source | Handler |
|-------|--------|---------|
| wallet.payout.completed | wallet-service | Update payout history |
| subscription.tier.activated | subscription-service | Notify creator of new subscriber |
| moderation.content.flagged | moderation-service | Flag creator channel for review |

## Schema (ChannelCreatedEvent)
`json
{"channelId":"uuid","creatorId":"uuid","channelName":"string","category":"gaming|music|education|entertainment","createdAt":"ISO8601"}
`
"@ | Set-Content (Join-Path C:\Users\Dell\Downloads\SPARK\services\creator-service "EVENTS.md") -Encoding UTF8

@"
# creator-service — Database Schema

## PostgreSQL — Primary Store
### channels table
| Column | Type | Notes |
|--------|------|-------|
| id | UUID | PK |
| creator_id | UUID | FK -> users.id, indexed |
| name | VARCHAR(100) | UNIQUE per platform |
| slug | VARCHAR(100) | URL-friendly, UNIQUE |
| description | TEXT | Markdown supported |
| category | VARCHAR(50) | From allowed categories |
| verification_status | ENUM | unverified, pending, verified, rejected |
| avatar_url | TEXT | CDN URL |
| banner_url | TEXT | CDN URL |
| status | ENUM | active, suspended, archived |
| created_at | TIMESTAMPTZ | |
| updated_at | TIMESTAMPTZ | |

### subscription_tiers table
| Column | Type | Notes |
|--------|------|-------|
| id | UUID | PK |
| channel_id | UUID | FK -> channels.id |
| name | VARCHAR(100) | e.g. "Gold Supporter" |
| price_cents | INTEGER | In USD cents, >= 100 |
| benefits | JSONB | Array of benefit descriptions |
| sort_order | INTEGER | Display ordering |

### creator_metrics table (materialized view)
Refreshed every 15 minutes. Contains subscriber count, 30-day revenue, content count, avg watch time.

Indexes: channels.creator_id, channels.slug UNIQUE, tiers.channel_id.
