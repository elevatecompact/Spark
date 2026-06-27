# Testing Strategy

## Scope

This document defines the testing strategy across all Titan engines and shared libraries. Every team follows the same approach, adapted to their language and domain.

## Quality Gates

### PR Level
All unit tests pass, code coverage does not decrease by more than 1%, linting passes with zero warnings, no new critical/high SAST findings, and contract tests pass for affected APIs.

### Staging Gate
Integration test suite completes (30+ minutes), load test baseline shows no regression greater than 5%, security DAST scan on staging URLs passes, and data validation pipelines pass.

### Production Gate
Canary health metrics meet SLO thresholds, synthetic transaction monitors pass, and chaos experiments from the current game day schedule passed in the current release window.

## Test Ownership

Engine teams own unit, integration, and contract tests for their service. The QA/SDET team owns E2E test suites, load tests, and test data management. The Platform/SRE team owns chaos experiments and production verification. The Security team owns security and fuzz testing.

## Test Environments

Local environment runs unit and integration tests with synthetic data via Docker Compose. CI uses ephemeral environments with Testcontainers in a Kind cluster. Staging uses a full EKS cluster with anonymized production replica data. Production runs canary, chaos, and synthetic tests against real user traffic.

## Continuous Improvement

Test failures are analyzed weekly in the Quality Review meeting. Test suite performance is tracked for total runtime, flake rate, and false positive rate. Obsolete or low-value tests are flagged for removal via a monthly audit.