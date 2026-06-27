# commerce-service — API Contract
## Products: POST /v1/products, GET /v1/products/{id}, PATCH /v1/products/{id}, DELETE /v1/products/{id}, GET /v1/products (storefront), POST /v1/products/{id}/variants
## Cart: GET /v1/cart (current user), POST /v1/cart/items (add), PATCH /v1/cart/items/{id} (update qty), DELETE /v1/cart/items/{id}, DELETE /v1/cart (clear)
## Checkout: POST /v1/checkout (create order), GET /v1/orders/{id}, GET /v1/orders (history), POST /v1/orders/{id}/cancel
## Fulfillment: POST /v1/orders/{id}/fulfill (digital delivery), GET /v1/orders/{id}/downloads (digital items), POST /v1/fulfillment/retry {orderItemId}
## Merchant: GET /v1/merchant/dashboard (sales, revenue), GET /v1/merchant/payouts, GET /v1/merchant/products (manage), POST /v1/merchant/storefront (configure)
## Reviews: POST /v1/products/{id}/reviews, GET /v1/products/{id}/reviews
## Admin: POST /v1/admin/products/{id}/feature, POST /v1/admin/orders/{id}/refund, GET /v1/admin/revenue
