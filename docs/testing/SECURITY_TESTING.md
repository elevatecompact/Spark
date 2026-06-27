# Security Testing

## Approach

Security at Titan is integrated into the development lifecycle, not bolted on at the end. Every PR is scanned, every dependency is vetted, every deployment is hardened.

## SAST (Static Analysis)

Semgrep runs custom rules for Titan-specific vulnerability patterns including mass assignment, insecure deserialization, and JWT manipulation. CodeQL provides GitHub-native analysis for JavaScript, TypeScript, Python, and Go. golangci-lint includes security linters such as `gosec` and `errcheck`.

## DAST (Dynamic Analysis)

OWASP ZAP runs automated scans of staging APIs and web UI on every staging deployment. Custom fuzzing targets API endpoints with invalid JSON, large payloads, and boundary values. Rate limit testing verifies that abusive clients are properly blocked.

## Dependency Scanning

All dependencies are scanned continuously. Trivy scans container images and filesystems for CVEs. Dependabot automates PRs for dependency updates. FOSSA ensures license compliance and dependency vulnerability mapping.

## Secret Detection

Gitleaks runs as a pre-commit hook and CI scanner for accidental credential commits. Pre-commit hooks run `detect-secrets` on every `git commit`. Scanning covers repositories, CI logs, build artifacts, and container images.

## Auth & Compliance

OAuth 2.0 and OIDC flow validation, JWT token validation for expiration and signature, RBAC matrix testing, and API key rotation testing are performed. PCI-DSS controls cover payment processing, SOC 2 controls are verified during every release, and GDPR data handling validation covers deletion and anonymization. Titan maintains a bug bounty program via HackerOne with critical vulnerabilities patched within 48 hours.