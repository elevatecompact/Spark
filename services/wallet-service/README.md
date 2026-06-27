# wallet-service — README

## Overview
The Wallet Service manages all monetary balances, transactions, and financial operations on the Titan platform. It handles multiple currency types (USD, Titan Tokens, creator coins), maintains append-only transaction ledgers, and ensures double-entry accounting integrity for all financial movements. Every financial operation flows through this service.

## Purpose
Provide a secure, auditable financial system that supports deposits via Stripe/PayPal, withdrawals to external accounts, peer-to-peer transfers, creator tipping, content purchases, and automated creator payouts. Implements strict reconciliation controls, fraud detection triggers, optimistic locking for concurrent balance updates, and comprehensive audit logging for regulatory compliance.

## Ownership
**Team:** Financial Systems (eng-finance@titan.dev)
**SLI:** 99.999% uptime, p99 balance read < 50ms
**Escalation:** #oncall-wallet
