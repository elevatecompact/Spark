# identity-service — Event Contracts

## Published Events
| Topic | Event | Description |
|-------|-------|-------------|
| iam.user.registered | UserRegisteredEvent | New account created via any auth method |
| iam.user.updated | UserUpdatedEvent | Profile fields changed (displayName, avatar) |
| iam.user.deleted | UserDeletedEvent | Account purged after grace period |
| iam.user.logged_in | UserLoggedInEvent | Successful authentication with method info |
| iam.user.mfa_enabled | UserMfaEnabledEvent | TOTP multi-factor activated |
| iam.user.suspended | UserSuspendedEvent | Admin-initiated account suspension |
| iam.token.revoked | TokenRevokedEvent | Token blacklisted (logout or rotation) |

## Consumed Events
| Topic | Source | Handler |
|-------|--------|---------|
| moderation.user.flagged | moderation-service | Apply login restrictions |
| payment.dispute.opened | payment-service | Flag account for review |

## Schema (UserRegisteredEvent)
`json
{
  "eventId": "uuid", "userId": "uuid", "email": "user@example.com",
  "authMethod": "email|google|github|discord", "registeredAt": "ISO8601"
}
`
All events use Avro with schema registry. Retention: 7 days.
