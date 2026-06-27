# gift-service — Database Schema
## gift_items (catalog): id UUID PK, name, price_cents, image_url, category(emote,badge,effect,sub), is_active, sort_order
## gifts (transactions): id UUID PK, sender_id FK, recipient_id FK, gift_item_id FK nullable, amount_cents, message, campaign_id FK nullable, status(pending,completed,refunded)
## gift_cards: id UUID PK, code VARCHAR UNIQUE, purchaser_id FK, balance_cents, expires_at, redeemed_at
## gift_campaigns: id UUID PK, creator_id FK, match_ratio FLOAT, max_match_cents, start_at, end_at
