# gift-service — Runbook
## Alerts: GiftSendFailureRate > 2%, GiftCardRedeemErrors > 5%, CampaignMatchBudgetExceeded, LeaderboardStale > 5m
## Refund: POST /v1/admin/gifts/{id}/refund
## Manual match: POST /v1/admin/campaigns/{id}/apply-match {giftId}
## Rebuild leaderboard: ./gift leaderboard rebuild
## Extend gift card: PATCH /v1/admin/gift-cards/{id} {expiresAt}
