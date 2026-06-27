# trust-service — README

## Overview
The Trust Service manages reputation scoring, trust signals, fraud detection, and risk assessment for the Titan platform. It computes user reputation scores based on platform behavior, detects fraudulent activity patterns, assesses risk levels for transactions and content actions, and maintains trust signals that other services use to make safety decisions.

## Purpose
Build and maintain a trustworthy platform ecosystem by quantifying user trustworthiness through behavioral analysis. Computes reputation scores combining positive signals (account age, verified status, content quality ratings, successful transactions, community contributions) and negative signals (reports received, moderation actions, payment disputes, suspicious behavior patterns). Provides real-time risk assessment for high-risk actions (large withdrawals, account changes, mass messaging) and fraud detection for coordinated inauthentic behavior, payment fraud, and account takeovers.

## Ownership
**Team:** Trust & Safety (eng-trust@titan.dev)
**SLI:** 99.95% uptime, p99 score calc < 300ms
**Escalation:** #oncall-trust
