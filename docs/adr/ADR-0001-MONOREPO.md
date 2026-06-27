# ADR-0001: Monorepo Structure

## Status

Accepted

## Context

Spark's engineering team is organized into multiple workstreams including platform, AI, streaming, API, and infrastructure. Early in development, each team independently provisioned repositories, leading to fragmented tooling, duplicated CI configuration, and inconsistent dependency management. Sharing type definitions required cross-repository publishing, and coordinated changes across teams demanded synchronized release cycles. As the organization scaled from three to eight teams, coordination overhead grew quadratically. The team evaluated polyrepo, monorepo with sparse checkout, and monorepo with scoped packages. Polyrepo offered team autonomy but at the cost of integration friction. Sparse checkout reduced clone times but added workflow complexity. Scoped packages within a monorepo provided a balance of isolation and discoverability.

## Decision

Adopt a single monorepository using npm workspaces with Yarn Berry and TypeScript project references. Each service, library, and configuration package lives under packages/ with a @spark/ scope prefix. Shared tooling (ESLint, Prettier, Jest, TypeScript) is configured at the root. CI pipelines use dependency graph detection to build and test only affected packages. Code ownership is enforced via CODEOWNERS files, and branch protection rules require review from owning teams.

## Consequences

### Positive
- Single source of truth for all code, types, and configuration
- Atomic cross-service changes with a single commit and CI run
- Shared tooling eliminates version skew across projects
- Simplified dependency management with hoisted node_modules

### Negative
- Clone and install times increase with repository size; mitigated by sparse checkout and zero-installs
- CI must implement smart change detection to avoid full-build overhead
- Requires team discipline around package boundaries and ownership
