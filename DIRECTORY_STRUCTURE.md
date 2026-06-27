# Directory Structure

```
spark/
+-- apps/                          # Deployable applications
¦   +-- web/                       # Next.js web application
¦   ¦   +-- src/                   # Application source code
¦   ¦   ¦   +-- app/               # Next.js App Router pages
¦   ¦   ¦   +-- components/        # Shared UI components
¦   ¦   ¦   +-- lib/               # Utility functions and hooks
¦   ¦   +-- public/                # Static assets
¦   ¦   +-- package.json
¦   +-- mobile/                    # Flutter mobile app
¦   ¦   +-- lib/                   # Dart source code
¦   ¦   +-- pubspec.yaml
¦   +-- admin/                     # Admin dashboard (Next.js)
¦
+-- packages/                      # Shared libraries
¦   +-- ui/                        # React component library
¦   +-- sdk/                       # TypeScript/Go SDK for API clients
¦   +-- config/                    # Shared configuration presets
¦   +-- types/                     # TypeScript type definitions
¦   +-- eslint-config/             # Shared ESLint configuration
¦
+-- services/                      # Backend microservices
¦   +-- gateway/                   # Envoy configuration and auth filter
¦   +-- auth/                      # Authentication service (Go)
¦   +-- content/                   # Content management service (Rust)
¦   +-- streaming/                 # Live streaming orchestration (Go)
¦   +-- analytics/                 # Real-time analytics service (Python)
¦   +-- moderation/                # Content moderation service (Python)
¦
+-- infra/                         # Infrastructure definitions
¦   +-- terraform/                 # Cloud resource provisioning
¦   +-- kubernetes/                # Kubernetes manifests and Helm charts
¦   +-- docker/                    # Base Docker images
¦
+-- tools/                         # CLI tools and scripts
¦   +-- cli/                       # Spark CLI (Go)
¦   +-- scripts/                   # Build and automation scripts
¦
+-- docs/                          # Documentation
¦   +-- architecture/              # ADRs and design docs
¦   +-- runbooks/                  # Operational runbooks
¦   +-- guides/                    # Developer guides
¦
+-- .github/                       # GitHub configuration
¦   +-- workflows/                 # CI/CD pipeline definitions
¦   +-- ISSUE_TEMPLATE/            # Issue and PR templates
¦
+-- pnpm-workspace.yaml           # pnpm workspace definition
+-- package.json                   # Root package.json
+-- turbo.json                     # Turborepo pipeline config
+-- tsconfig.json                  # Root TypeScript configuration
+-- .env.example                   # Environment variable template
```
