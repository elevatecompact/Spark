# Testing — Titan Platform

## Overview

Testing at Titan is a multi-layered discipline spanning unit-level verification through production-scale chaos experiments. Every change is validated at multiple points in the delivery pipeline, with automated gates at each stage.

## Testing Philosophy

- **Shift Left**: Find defects as early as possible. Unit tests run on every keystroke; integration tests run on every PR.
- **Test the Behavior, Not the Implementation**: Tests focus on observable outcomes. Internal refactoring should not break tests.
- **Coverage Is a Signal, Not a Target**: Coverage helps identify untested paths, but the goal is meaningful assertions, not line counts.
- **Deterministic by Default**: Flaky tests are treated as production incidents. Any test that fails nondeterministically is quarantined and fixed within 48 hours.

## Test Pyramid

```
    /\
   /  \        E2E (5%) — full user journeys
  /    \
 /      \      Integration (20%) — API contracts, DB, cross-engine
/________\
 Unit (75%) — individual functions, services, business logic
```

## Toolchain

| Layer | Tool |
|-------|------|
| Unit | Go testing, pytest, Jest |
| Integration | Testcontainers, Docker Compose |
| Contract | Pact (CDC testing) |
| E2E | Playwright, k6 |
| Load | k6, Vegeta |
| Chaos | Chaos Mesh, Litmus |
| Security | OWASP ZAP, Semgrep, Trivy |
| Test Data | Faker, factory_boy, gofake |

All test runs produce a JUnit XML report, published to a centralized test intelligence dashboard.