# Security Checklist

This checklist covers operational security controls and best practices for the Spark platform. Use this document during deployment reviews, incident response drills, and security audits.

## Infrastructure Security

- [ ] All network traffic uses TLS 1.3 with strong cipher suites
- [ ] mTLS is enforced for all inter-service communication
- [ ] No service runs with root or privileged access
- [ ] All containers and VMs use minimal base images
- [ ] Regular vulnerability scanning (daily) with CVSS > 7.0 remediated within 7 days
- [ ] Network segmentation and micro-segmentation are enforced
- [ ] Firewall rules follow least-privilege; unused ports are closed
- [ ] Backups are encrypted and tested monthly

## Application Security

- [ ] All user input is validated and sanitized
- [ ] No secrets in source code, config files, or environment variables
- [ ] Authentication uses OAuth 2.0 / OIDC; no custom auth protocols
- [ ] MFA is enforced for all privileged accounts
- [ ] Session tokens are short-lived and rotated
- [ ] API rate limiting is configured for all endpoints
- [ ] Content Security Policy (CSP) headers are set on all responses
- [ ] Dependency scanning runs with every build; known-vulnerability policy enforced

## Data Security

- [ ] All data at rest is encrypted with AES-256
- [ ] All data in transit is encrypted with TLS 1.3
- [ ] PII is minimized, tokenized, or anonymized where possible
- [ ] Data retention policies are configured and enforced
- [ ] Access to production data is gated by JIT approval

## Operational Security

- [ ] Access reviews are conducted quarterly
- [ ] Incident response runbooks are reviewed and tested quarterly
- [ ] Penetration testing is performed annually
- [ ] Third-party vendors are assessed for security posture
- [ ] Security awareness training is completed annually by all employees
- [ ] Disaster recovery drills are conducted semi-annually

## Continuous Improvement

- [ ] Threat models are updated for every significant feature change
- [ ] Security debt is tracked and resolved with engineering milestones
- [ ] Post-incident reviews result in actionable improvements
- [ ] Compliance evidence is collected and reviewed continuously
