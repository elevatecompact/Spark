# Vault Engine

**Purpose:** Secure payment processing and billing engine for the Titan platform.
**Tech Stack:** Go, Stripe API, PayPal API, PostgreSQL, Redis, gRPC, PCI-DSS compliant.

Vault handles subscription billing, one-time payments, invoicing, tax calculation, refunds, and payout reconciliation. Integrates with Stripe and PayPal while keeping sensitive data within PCI-DSS boundaries. Supports usage-based billing, coupon management, dunning, and billing analytics.
