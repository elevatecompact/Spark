# gift-service — Configuration
GIFT_PORT=4011, GIFT_DB_URL, GIFT_REDIS_URL, GIFT_CARD_CODE_LENGTH=12, GIFT_CARD_EXPIRY_DAYS=365, MAX_GIFT_AMOUNT_CENTS=1000000 (), MIN_GIFT_AMOUNT_CENTS=100 (), CAMPAIGN_MAX_DURATION_DAYS=30, LEADERBOARD_UPDATE_INTERVAL=60s
FF: gift_sending_enabled=true, gift_cards_enabled=true, campaign_matching=true, anonymous_gifting=true, gift_leaderboard=true
RL: 50 gifts/h per sender, 5 gift cards/day per user, 3 campaigns/month per creator
