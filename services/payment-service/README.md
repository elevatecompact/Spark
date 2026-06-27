# payment-service — README

## Overview
The Payment Service is the payment processing engine for the Titan platform. It integrates with external payment gateways (Stripe as primary, PayPal as secondary), manages payment method tokenization, handles incoming webhooks with signature verification, and provides a unified payment abstraction for all services that need to accept or disburse money.

## Purpose
Abstract away the complexity of multiple payment processors behind a single internal API. Handles payment intent creation with idempotency, confirmation and capture, full and partial refunds, dispute management (chargebacks), payout initiation to connected accounts, and PCI-compliant payment method storage via client-side tokenization. All sensitive card data is handled by Stripe/PayPal directly — Titan never touches raw PAN data.

## Ownership
**Team:** Financial Systems (eng-finance@titan.dev)
**SLI:** 99.999% uptime, p99 intent < 200ms
**Escalation:** #oncall-payment
