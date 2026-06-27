# payment-service — Testing Guide
## Unit: Payment intent state machine, webhook signature verification (Stripe+PayPal), idempotency, refund validation, currency formatting.
## Integration: Create→confirm→succeed, failed payment handling, refund lifecycle, payment method save/reuse, webhook processing.
## Mocks: Stripe API via WireMock (tests/mocks/stripe/), PayPal similarly. Network failure scenarios (timeout, 500, reset).
