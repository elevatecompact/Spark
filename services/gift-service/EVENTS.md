# gift-service — Event Contracts
## Published: gift.sent, gift.received, gift.subscription.gifted, gift.campaign.match, gift.card.purchased, gift.card.redeemed
## Consumed: wallet.transaction.settled (confirm payment), subscription.activated (confirm gifted sub), creator.stream.started (apply campaign matching)
## Schema: GiftSentEvent {giftId, senderId, recipientId, giftType(virtual_item|subscription|gift_card), amountCents, message, campaignId, sentAt}
