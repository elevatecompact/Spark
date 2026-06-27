# Helm Charts

## Chart Structure

Titan maintains a single umbrella chart and individual subcharts per engine under `charts/`. This allows deploying the full platform or individual engines.

## Repository

Charts are published to an OCI-compatible registry in Amazon ECR. Versioning follows semver and aligns with platform releases.

## Chart Architecture

```
charts/
  titan/                  Umbrella chart
    Chart.yaml
    values.yaml           Global defaults
    charts/
      pulse/              Pulse engine subchart
      oracle/             Oracle engine subchart
      atlas/              Atlas engine subchart
      ...
  titan-infra/            Infrastructure dependencies (Redis, Postgres, MinIO)
```

## Values Management

Global values set defaults for image tags, replica counts, resource requests/limits. Environment overrides live in `charts/titan/values-{env}.yaml`. Secrets are never stored in values files; referenced via `externalSecrets` annotation.

## Deployment Flow

CI builds image, pushes to ECR, updates `values-{env}.yaml` with new image tag. ArgoCD detects drift and syncs the Helm release. Helm upgrade uses `--atomic --cleanup-on-fail` to prevent partial deployments.

## Testing & Best Practices

`helm test` hooks validate service readiness post-install. Chart linting runs via `ct lint` in CI. Dependencies are managed via `helm dependency update`. `helm template` output is diffed against the live cluster in staging before production promotion. Annotations include `app.kubernetes.io/name`, `app.kubernetes.io/instance`, and `app.kubernetes.io/version`.