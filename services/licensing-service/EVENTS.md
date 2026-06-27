# licensing-service — Event Contracts
## Published: licensing.license.created, licensing.license.activated, licensing.license.expired, licensing.license.terminated, licensing.usage.recorded, licensing.royalty.calculated, licensing.royalty.paid, licensing.compliance.flag
## Consumed: stream.session.started (verify stream content rights), media.content.uploaded (register copyright), commerce.order.placed (verify licensed content sale), moderation.content.flagged (copyright claim), wallet.payout.completed (confirm royalty payment)
## Schema: LicensingRoyaltyCalculatedEvent {licenseId, contentId, periodStart, periodEnd, usageCount, ratePerUse, totalCents, rightsHolderId, calculatedAt}
