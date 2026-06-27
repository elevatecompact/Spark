# Architecture

Vault has three layers: payment gateway abstraction, billing core, and reconciliation engine. Gateway abstraction implements PaymentProvider interface for Stripe, PayPal, and mock providers. Billing core manages subscriptions, invoices, and payment intents using PostgreSQL with Redis for idempotency keys. Reconciliation engine runs hourly to match bank statements against processed payments. Webhook handler processes async events with idempotent replay protection.
