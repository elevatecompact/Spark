# Spark Engineering Laws

These laws govern all engineering work across the Spark platform. They are not guidelines. They are binding rules. Every commit, every PR, every deployment, every line of code must comply.

Violation of these laws is a system defect, not a disagreement.

---

## Law 1 — Reality First

Reality is the source of truth. Production behavior beats staging. Staging beats tests. Tests beat specifications. Specifications beat assumptions.

When observability data conflicts with a developer's belief, the data wins.

Every incident must produce a permanent artifact. Every hypothesis must be tested against production signals before being accepted.

---

## Law 2 — No Unverified Execution

Nothing runs without proof.

Every change must be:
- **Reviewed** — at least one human or automated reviewer
- **Tested** — unit, integration, or contract tests covering the change
- **Observable** — logs, metrics, or traces proving the change works
- **Traceable** — linked to a ticket, PR, and deployment event

A change without evidence is not complete.

---

## Law 3 — No Authority Bypass

No individual, team, service, or external dependency may bypass governance.

All exceptions require:
1. Written justification
2. Time-bound expiration
3. Review by two engineers outside the requesting team
4. Documentation in the relevant ADR

Permanent exceptions do not exist.

---

## Law 4 — Continuous Improvement

Every failure is information. Every success is evidence. Every cycle must improve the platform.

Post-incident reviews are mandatory and blameless. Findings must produce:
- Automated checks that prevent recurrence
- Monitoring that detects similar patterns
- Documentation that captures the lesson

If it happens once, it's an incident. If it happens twice, it's a process failure.

---

## Law 5 — Causal Accountability

Every decision must have a cause. Every cause must be discoverable.

This means:
- Every PR description must state *why* the change exists
- Every configuration change must be version-controlled
- Every deployment must link to the set of changes included
- Every feature flag must have an owner and expiration

If a decision cannot be explained six months later, it was not made properly.

---

## Law 6 — Economic Rationality

Resources are finite. Time is finite. Attention is finite. Compute is finite.

Every service must:
- Define its expected resource footprint
- Justify its cost-to-value ratio quarterly
- Remove unused code, dependencies, and infrastructure
- Optimize for the marginal cost per request

Performance budgets are not targets. They are ceilings.

---

## Law 7 — Multi-Reality Evaluation

Never commit to one future when many futures can be explored.

Before significant architectural decisions:
- Evaluate at least three alternatives
- Document trade-offs in an ADR
- Prototype the riskiest unknown first
- Use feature flags to test in production

If you cannot articulate why the chosen approach is better than the alternatives, you are not ready to decide.

---

## Law 8 — Governed Autonomy

Teams own their services. Autonomy is permitted. Ungoverned autonomy is forbidden.

Service teams may:
- Choose internal implementation details
- Own their deployment cadence
- Define their testing strategy

Service teams must:
- Maintain their documentation
- Meet platform-wide SLOs
- Participate in on-call rotations
- Publish their API contracts

Freedom within framework. Not freedom from framework.

---

## Law 9 — Preservation of Knowledge

Knowledge gained must not be lost.

Every service must maintain:
- Architecture decisions and their rationale
- Runbooks for common operational tasks
- Troubleshooting guides for known failure modes
- A changelog documenting significant changes

If a team member leaves and knowledge leaves with them, the system has failed.

---

## Law 10 — Service to Human Intent

The platform serves humans. Humans define purpose. Technology defines execution.

Every feature must answer:
- What human need does this serve?
- How do we measure whether it serves that need?
- How do we know if it causes harm?

Technology choices are subordinate to user outcomes. If a technical decision cannot be explained in terms of user value, it is not justified.

---

## Enforcement

These laws are enforced through:

| Mechanism | Frequency |
|-----------|-----------|
| Automated CI checks | Every PR |
| Architecture review | Every significant change |
| Service ownership review | Quarterly |
| Incident post-mortem | Every incident |
| Architecture audit | Annually |

---

## Amendments

Changes to these laws require:
1. An ADR documenting the proposed change
2. Review by the Architecture Review Board
3. A 2-week comment period for all engineering teams
4. Ratification by the CTO

No single individual may amend these laws.

---

*Last updated: 2026-06-27*
*Version: 1.0.0*
