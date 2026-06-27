# Vault Engine — Payments

## Purpose

Vault is Titan's payment processing engine. It handles subscription billing, one-time purchases, payout management, and wallet operations. Vault abstracts multiple payment providers behind a unified API with idempotent operations and comprehensive reconciliation.

## Architecture

Vault uses a dual-write pattern: events are written to the payment provider (Stripe, crypto) and to Vault's own ledger database. Reconciliation jobs run hourly to detect and resolve discrepancies between the two systems.

## Tech Stack

- **Language**: Go
- **Database**: PostgreSQL (ledger, subscriptions, invoices, wallets)
- **Cache**: Redis for idempotency keys and rate limiting
- **Payment Providers**: Stripe (primary), PayPal, cryptocurrency (Solana, USDC)
- **Messaging**: RabbitMQ for asynchronous billing events
- **Tax**: Automated tax calculation via Stripe Tax and region-specific logic

## Key Features

- **Subscription management**: Plans, trials, upgrades, downgrades, cancellations with proration
- **One-time purchases**: Digital goods, credits, in-app purchases
- **Multi-provider**: Stripe as primary with fallback to PayPal and crypto
- **Idempotent operations**: Every payment operation has an idempotency key to prevent double charges
- **Wallet system**: Platform credit balance for prepaid usage
- **Invoicing**: Automated invoice generation, delivery, and dunning
- **Payouts**: Creator/subscriber payout processing with configurable thresholds
- **Tax compliance**: Automated sales tax, VAT, GST calculation and remittance reporting
- **Refunds**: Full and partial refunds with automatic fee reversal

## Performance Targets

| Metric | Target |
|--------|--------|
| Payment processing latency | < 500ms (p99) |
| Idempotency deduplication | 100% |
| Reconciliation accuracy | 100% |
| Dunning success rate | > 90% |
| Uptime | 99.99% |