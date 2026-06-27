# Penetration Testing Program

Spark operates a continuous penetration testing program that combines automated scanning, manual testing by internal teams, and third-party assessments.

## Scope

Penetration tests cover all externally facing applications, internal services, APIs, mobile applications, and infrastructure components. Tests are conducted against production and staging environments, with appropriate safeguards to prevent service disruption.

## Testing Cadence

| Test Type | Frequency |
|---|---|
| Automated vulnerability scanning (DAST, SAST) | Continuous (every build) |
| Internal manual penetration testing | Quarterly |
| Third-party penetration testing | Annually (or after major architecture changes) |
| Bug bounty program | Continuous |

## Methodology

Tests follow the OWASP Testing Guide and NIST SP 800-115. The testing process includes:

1. **Reconnaissance** — Subdomain enumeration, service discovery, technology fingerprinting
2. **Threat Modeling** — Identifying attack vectors specific to the target
3. **Exploitation** — Attempting to exploit identified vulnerabilities in a controlled manner
4. **Privilege Escalation** — Attempting to move laterally or elevate access
5. **Persistence** — Testing whether persistent access can be established
6. **Reporting** — Detailed findings with reproduction steps, impact, and remediation

## Bug Bounty Program

Spark invites security researchers to find and report vulnerabilities through a managed bug bounty platform. Rewards are tiered by severity. All valid reports receive acknowledgment and are tracked to resolution.

## Remediation SLAs

| Severity | Remediation Timeframe |
|---|---|
| Critical | 24 hours |
| High | 7 days |
| Medium | 30 days |
| Low | 90 days |

All penetration test findings are logged in the security tracker and assigned to the responsible engineering team for remediation.
