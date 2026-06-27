# Incident Response Plan

Spark maintains a structured incident response (IR) plan aligned with the NIST SP 800-61 framework. The plan defines roles, procedures, and communication channels for responding to security incidents.

## Incident Classification

| Severity | Description | Response Time |
|---|---|---|
| P1 (Critical) | Active data breach, full system compromise, ransomware | Immediate, 24/7 |
| P2 (High) | Unauthorized access to restricted data, service outage | Within 1 hour |
| P3 (Medium) | Suspicious activity, isolated compromise | Within 4 hours |
| P4 (Low) | Phishing attempt, policy violation | Next business day |

## Roles and Responsibilities

- **Incident Commander (IC)** — Coordinates response, makes decisions, and communicates with stakeholders
- **SME (Security)** — Leads technical investigation, containment, and eradication
- **Legal** — Advises on legal obligations, data breach notification requirements
- **Communications** — Manages external communications, customer notifications
- **Engineering** — Implements fixes and recovery procedures

## Response Lifecycle

### Preparation
- Incident response runbooks for common scenarios (data breach, ransomware, account takeover)
- Regularly tested backup and restore procedures
- Communication templates pre-approved by legal

### Detection and Analysis
- Events from SIEM, EDR, fraud-ai, and user reports are triaged
- Indicators of compromise (IOCs) are collected and analyzed
- Scope and impact are determined

### Containment, Eradication, Recovery
- Short-term: isolate affected systems, block IOCs
- Long-term: patch vulnerabilities, rotate credentials
- Recovery: restore from clean backups, verify integrity

### Post-Incident
- Root cause analysis (RCA) document
- Lessons learned review
- Remediation items tracked in security backlog
