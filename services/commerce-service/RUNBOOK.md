# commerce-service — Runbook
## Alerts: CheckoutFailureRate > 2%, InventoryDesyncDetected (count mismatch), FulfillmentDelay > 5min, MerchantPayoutFailure, OrderRefundRate > 5%
## Refund order: POST /v1/admin/orders/{id}/refund {full:true|partial, amountCents}
## Force fulfill: POST /v1/admin/orders/{id}/fulfill — triggers digital delivery manually.
## Inventory reconcile: POST /v1/admin/products/{id}/inventory-reconcile {actualCount}
## Cancel order: POST /v1/admin/orders/{id}/cancel
