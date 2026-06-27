# subscription-service — README

## Overview
The Subscription Service manages all recurring subscription relationships on the Titan platform. It handles subscription plans, billing cycles with automated payment collection via the wallet and payment services, tier management with granular benefits, and subscriber benefits enforcement at access points throughout the platform.

## Purpose
Enable creators to monetize their audience through tiered subscription plans with recurring billing (monthly and yearly). Manage the full subscription lifecycle: sign-up and activation, scheduled billing cycles with retry logic, 3-day grace periods after failed payments, cancellations at period end, mid-cycle plan upgrades with proration, and expired access revocation. Integrates deeply with the gating system for subscriber-only content.

## Ownership
**Team:** Monetization (eng-monetization@titan.dev)
**SLI:** 99.99% uptime, p99 activation < 200ms
**Escalation:** #oncall-subscriptions
