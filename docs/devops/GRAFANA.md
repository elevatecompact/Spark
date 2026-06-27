# Grafana Dashboards

## Dashboard Management

All Grafana dashboards are defined as code in the `observability/dashboards/` directory, managed via Terraform (grafana provider) or Kubernetes ConfigMaps. Dashboards are version-controlled, reviewed in PRs, and automatically provisioned.

## Dashboard Categories

### Service Dashboards (per engine)
Request rate, error rate, duration (RED) — four golden signals. Upstream/downstream dependency latency. Resource utilization and saturation. SLO burn rate with error budget remaining.

### Platform Dashboards
Kubernetes Cluster Overview: node health, pod density, namespace quotas, network flows. Deployment Dashboard: blue-green and canary health, rollout progress, rollback events. Alerting Dashboard: firing alerts, silences, notification delivery status.

### Business Dashboards
User Engagement: active users, session duration, content consumption velocity. Content Pipeline: upload volume, transcoding queue depth, CDN cache hit ratio. Revenue: payment success rate, refund rate, MRR.

## Dashboard Standards

Every dashboard follows common conventions: left Y-axis for rates/totals, right Y-axis for latencies, unit annotations in panel titles, template variables for `$service`, `$cluster`, `$namespace`, `$timeRange`, links to related dashboards and runbooks in panel descriptions, and a minimum 24h time range default for meaningful SLO context.

## Access Control

Read-only access is granted to all engineers. Dashboard editing is restricted to the Platform team via Terraform PRs only. Admin access is reserved for the SRE team.