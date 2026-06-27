# commerce-service — Database Schema
## products: id UUID PK, creator_id FK, name VARCHAR, description TEXT, type(digital,physical,download,bundle), price_cents, currency, category, tags TEXT[], media_urls TEXT[], inventory_count INT (-1=unlimited), is_active, is_featured, created_at
## product_variants: id UUID PK, product_id FK, name VARCHAR (e.g. "Large"), price_cents (override), inventory_count, sort_order
## cart_items: id UUID PK, user_id FK, product_id FK, variant_id FK nullable, quantity INT, added_at (TTL-based cleanup for abandoned carts)
## orders: id UUID PK, buyer_id FK, merchant_id FK, status(pending,paid,fulfilled,cancelled,refunded), total_cents, currency, payment_intent_id FK nullable, shipping_address JSONB, placed_at, fulfilled_at
## order_items: id UUID PK, order_id FK, product_id FK, variant_id FK, quantity, unit_price_cents, fulfillment_status(pending,fulfilled), download_url (digital), fulfilled_at
## eviews: id UUID PK, product_id FK, user_id FK, rating INT(1-5), title VARCHAR, body TEXT, is_verified_purchase, created_at. UNIQUE(product_id, user_id)
## Redis: Cart cache (TTL 7 days for logged-in, 30 min for guest), product cache (TTL 1min), inventory counters (real-time deduction)
