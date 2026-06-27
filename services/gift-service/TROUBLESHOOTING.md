# gift-service — Troubleshooting
## Gift not delivered: Kafka failure, recipient wallet not found, amount validation failed post-settlement. Check gift status in DB, verify recipient wallet exists, check Kafka gift.sent topic.
## Gift card invalid: Code entered wrong, expired, already redeemed. Look up by code, check expiry/redeemed_at.
## Campaign match not applied: Budget exhausted, campaign not active, gift sent too early. Verify campaign, check total matched vs max, manually apply.
