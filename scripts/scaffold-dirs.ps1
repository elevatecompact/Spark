$root = "C:\Users\Dell\Downloads\SPARK"

$dirs = @(
    # Apps
    'apps/web', 'apps/creator-dashboard', 'apps/moderation-console', 'apps/admin',
    'apps/developer-portal', 'apps/tv', 'apps/desktop', 'apps/mobile', 'apps/landing', 'apps/docs',
    # Services
    'services/identity-service', 'services/creator-service', 'services/viewer-service',
    'services/stream-service', 'services/chat-service', 'services/messaging-service',
    'services/wallet-service', 'services/subscription-service', 'services/gift-service',
    'services/payment-service', 'services/analytics-service', 'services/notification-service',
    'services/recommendation-service', 'services/search-service', 'services/translation-service',
    'services/moderation-service', 'services/community-service', 'services/event-service',
    'services/competition-service', 'services/advertising-service', 'services/commerce-service',
    'services/media-service', 'services/licensing-service', 'services/discovery-service',
    'services/trust-service',
    # Engines
    'engines/pulse', 'engines/atlas', 'engines/oracle', 'engines/forge', 'engines/echo',
    'engines/sentinel', 'engines/guardian', 'engines/velocity', 'engines/vault', 'engines/vision',
    'engines/nexus', 'engines/polyglot', 'engines/sparkclips',
    # AI
    'ai/creator-ai', 'ai/recommendation-ai', 'ai/moderation-ai', 'ai/translation-ai',
    'ai/voice-ai', 'ai/vision-ai', 'ai/clip-ai', 'ai/thumbnail-ai', 'ai/ranking-ai',
    'ai/fraud-ai', 'ai/assistant-ai',
    # Packages
    'packages/ui', 'packages/icons', 'packages/tokens', 'packages/components', 'packages/auth',
    'packages/logger', 'packages/telemetry', 'packages/events', 'packages/config', 'packages/database',
    'packages/types', 'packages/utils', 'packages/validators', 'packages/design-system',
    # API
    'api/gateway', 'api/graphql', 'api/rest', 'api/grpc', 'api/websocket', 'api/webrtc',
    'api/stream', 'api/webhooks',
    # SDKs
    'sdk/typescript', 'sdk/go', 'sdk/rust', 'sdk/python', 'sdk/java', 'sdk/dotnet',
    'sdk/flutter', 'sdk/swift', 'sdk/kotlin',
    # Database
    'database/postgres', 'database/redis', 'database/clickhouse', 'database/opensearch',
    'database/neo4j', 'database/migrations', 'database/schemas', 'database/seed',
    # Infrastructure
    'infrastructure/kubernetes', 'infrastructure/terraform', 'infrastructure/helm',
    'infrastructure/cloudflare', 'infrastructure/aws', 'infrastructure/gcp',
    'infrastructure/azure', 'infrastructure/edge', 'infrastructure/network',
    # Deployment
    'deployment/development', 'deployment/staging', 'deployment/production',
    'deployment/disaster-recovery', 'deployment/chaos-testing',
    # Observability
    'observability/grafana', 'observability/prometheus', 'observability/tempo',
    'observability/loki', 'observability/otel', 'observability/sentry', 'observability/dashboards',
    # Security
    'security/oauth', 'security/oidc', 'security/passkeys', 'security/jwt',
    'security/encryption', 'security/audit', 'security/compliance', 'security/permissions',
    'security/secrets',
    # Documents - main
    'docs/architecture', 'docs/engines', 'docs/ai', 'docs/security', 'docs/database',
    'docs/api', 'docs/platform', 'docs/devops', 'docs/testing', 'docs/adr',
    # Tests
    'tests/unit', 'tests/integration', 'tests/contract', 'tests/e2e',
    'tests/load', 'tests/chaos', 'tests/performance',
    # Other top-level
    'platform', 'scripts', 'tools',
    # Future
    'marketplace', 'academy', 'jobs', 'music', 'podcasts', 'virtual-events',
    'vr', 'ar', 'games', 'creator-fund', 'enterprise', 'education', 'sports', 'news'
)

foreach ($d in $dirs) {
    $null = New-Item -ItemType Directory -Path (Join-Path $root $d) -Force
}

Write-Output "Created $($dirs.Count) directories"
