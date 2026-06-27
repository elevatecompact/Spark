# Authorization Model

Spark implements a hybrid Role-Based Access Control (RBAC) and Attribute-Based Access Control (ABAC) model, combining the manageability of roles with the granularity of attribute-driven policies.

## RBAC Core

Roles are hierarchical and support inheritance. Every user is assigned one or more roles, each of which grants a set of permissions (actions on resources). Standard roles include:

- **Viewer** — Read-only access to assigned resources
- **Operator** — Ability to modify operational state but not configuration
- **Admin** — Full control over a defined scope
- **Super Admin** — System-wide administrative access (strictly limited)

Custom roles can be defined per tenant, allowing organizations to mirror their internal structure.

## ABAC Extensions

Where RBAC alone is insufficient, ABAC policies evaluate contextual attributes at runtime:

- **User attributes** — department, clearance level, employment status
- **Resource attributes** — data classification, sensitivity, owner
- **Environment attributes** — time of day, network location, device posture
- **Action attributes** — read, write, delete, export

## Policy Evaluation

Policies are defined in a declarative policy language (Rego/OPA-compatible) and evaluated at every request. The evaluation engine combines RBAC permissions with ABAC constraints, returning a binary allow/deny decision along with a reason code for audit purposes.

## Audit Trail

Every authorization decision is logged with the identity, resource, action, policy matched, and outcome. This provides a complete chain of evidence for compliance audits and forensic investigations.
