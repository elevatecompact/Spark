# identity-service — Ownership

## Engineering Team
**Platform Security** — eng-security@titan.dev

| Role | Name | Contact |
|------|------|---------|
| Tech Lead | Alex Chen | alex.chen@titan.dev |
| Engineering Manager | Sarah Wu | sarah.wu@titan.dev |
| Security Architect | David Kim | david.kim@titan.dev |

### Ownership Boundary
| Component | Owner |
|-----------|-------|
| Auth flows (login, register, MFA) | Platform Security |
| OAuth provider integrations | Platform Security |
| Session & token management | Platform Security |
| User profile CRUD | Platform Security |
| API key lifecycle | Platform Security |

### On-Call
- Primary: Weekly rotation among 6 team members
- Secondary: Escalates to Tech Lead
- Schedule: PagerDuty identity-oncall

### Code Review
All PRs require approval from at least one Platform Security engineer. Use @eng-security for review requests.

### Key Dependencies
- **notification-service** — Auth event emails (welcome, password reset, suspicious login)
- **moderation-service** — Account suspension enforcement
- **analytics-service** — Login funnel and conversion metrics
