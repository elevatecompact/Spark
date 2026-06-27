# Alerting

## Alerting Philosophy

Titan follows a tiered alerting model. Every alert is actionable, documented, and has a defined owner. Alerts that do not trigger a specific response are removed or demoted to warnings.

## Alert Tiers

| Tier | Response Time | Channel | Examples |
|------|---------------|---------|----------|
| P0 (Critical) | 5 minutes | PagerDuty phone call | Service down, data loss, security breach |
| P1 (High) | 15 minutes | PagerDuty push | High error rate, latency spike, certificate expiry < 24h |
| P2 (Medium) | 1 hour | Slack notification | Warning thresholds, degraded performance |
| P3 (Low) | 8 hours | Email digest | Disk usage > 80%, old certificate nearing renewal |

## Alert Definition

Alerts are defined as PrometheusRule custom resources, co-located with service code in `deploy/prometheus-rules.yaml`. Every alert rule includes a `summary`, `description` with contextual details, `runbook_url` linking to the relevant runbook, and `severity` (P0-P3).

## Notification Routing

P0/P1 alerts route through PagerDuty with escalation from primary to secondary to SRE team. P2 alerts post to Slack in `#alerts` with `@here` for the on-call team. P3 alerts go to a weekly ops review email.

## Alert Fatigue Prevention

Flapping detection suppresses notifications after 3 state changes within 10 minutes. Similar alerts are grouped by `service` and `alertname`. Burn rate alerts based on SLO consumption reduce noise compared to raw error counts. Alert rules are unit-tested using `promtool test rules`.