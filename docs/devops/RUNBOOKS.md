# Runbooks

## Runbook Structure

Every production incident must have a documented runbook. Runbooks live in `runbooks/` at the repository root and are written in Markdown with standardized frontmatter including title, affected service, severity (P0-P3), and tags.

## Required Sections

1. **Symptoms** — What the on-call engineer will see in alerts and dashboards
2. **Severity Assessment** — Criteria to classify the incident (P0-P3)
3. **Immediate Actions** — Step-by-step instructions to mitigate (numbered steps, no ambiguity)
4. **Verification** — How to confirm the mitigation worked (dashboards, health checks, log queries)
5. **Resolution** — Permanent fix steps once the immediate incident is handled
6. **Communication** — Who to notify, what Slack channel to update, and at what cadence
7. **Post-Incident** — Links to template files for postmortem

## Runbook Catalog

Key runbooks include service-down procedures for all engines, database failover for RDS multi-AZ, certificate expiry renewal, deployment rollback for blue-green and canary, data corruption recovery via point-in-time restore, and security incident handling for credential leaks and unauthorized access.

## Maintenance

Every runbook is tested during Chaos Engineering GameDays on a quarterly basis. Runbooks are updated within 1 week of any infrastructure change that affects the documented procedure. A `runbooks/README.md` index file lists all runbooks with tags and last-reviewed date.

## Tooling

Runbooks are linked directly from Alertmanager notifications via the `runbook_url` annotation. Grafana dashboard panels include a View Runbook link in panel descriptions. A `runbookctl` CLI tool allows searching by tag, service, or keyword from the terminal.