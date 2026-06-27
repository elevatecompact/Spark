# Release Process

Spark follows a **time-based release cadence** with a continuous delivery model for patch releases.

## Cadence

| Release Type | Frequency | Version Bump | Examples |
|-------------|----------|-------------|----------|
| Major | Every 6 months | X.0.0 | 1.0.0, 2.0.0 |
| Minor | Bi-weekly | 0.X.0 | 0.3.0, 0.4.0 |
| Patch | As needed | 0.0.X | 0.3.1, 0.3.2 |
| Hotfix | Emergency | 0.0.X | Critical CVEs, outages |

## Release Process

1. **Cut release branch** from `main` named `release/v<version>`.
2. **Run full suite** of CI checks — lint, type-check, unit tests, integration tests, E2E tests.
3. **Generate changelog** from conventional commit messages since last release.
4. **Create release candidate** (RC) and deploy to staging environment.
5. **Smoke test** staging — verify critical paths, run performance benchmarks.
6. **Tag and sign** the release commit: `git tag -s v<version> -m "Spark v<version>"`.
7. **Deploy to production** using the infrastructure pipeline.
8. **Monitor** dashboards for 24 hours post-release — rollback immediately if error rate exceeds 1%.

## Versioning Convention

We follow [SemVer 2.0](https://semver.org/): `MAJOR.MINOR.PATCH`. Breaking changes increment the MAJOR version. Backward-compatible features increment MINOR. Bug fixes increment PATCH.

## Hotfix Process

Hotfixes bypass the normal release cycle. A fix is applied to the latest release branch, tested, and deployed immediately. The fix is then merged back into `main` and staged for the next regular release.
