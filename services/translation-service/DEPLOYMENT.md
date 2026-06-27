# translation-service — Deployment Guide
K8s: k8s/translation-service/ — api (3x 1GB), batch-worker (2x 1GB for batch translation), streaming (2x 2GB for WS).
External API calls to DeepL/Google require static egress IPs (allowlisted). Provider failover: if DeepL returns 5xx, auto-failover to Google Translate within 30s.
Deploy: kubectl apply -f k8s/translation-service/. Translation memory pre-seeded with common phrases.
Health: /health (DB+Redis+DeepL API), /ready (provider reachable), /metrics :4115.
