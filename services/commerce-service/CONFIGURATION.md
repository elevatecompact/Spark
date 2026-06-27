# commerce-service — Configuration
COMMERCE_PORT=4023, COMMERCE_DB_URL, COMMERCE_REDIS_URL, COMMERCE_KAFKA_BROKERS, MAX_CART_ITEMS=50, CART_TTL_DAYS=7, ORDER_TIMEOUT_MINUTES=30 (unpaid order expiry), MAX_PRODUCTS_PER_CREATOR=100, DIGITAL_DOWNLOAD_EXPIRY_HOURS=72, REVIEW_MODERATION_ENABLED=true, MERCHANT_PAYOUT_CRON="0 0 1 * *" (monthly)
FF: commerce_enabled=true, digital_products=true, product_variants=true, reviews_enabled=true, guest_checkout=false, merchant_payouts_auto=true
Rate limits: 10 products/month per creator, 50 cart operations/h per user, 5 reviews/day per user, 10 orders/min per user
