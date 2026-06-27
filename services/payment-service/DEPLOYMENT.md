# payment-service — Deployment Guide
Critical: PCI DSS compliant infra, Vault for API keys, SOC2 Type II, multi-region active-passive DR.
K8s: 3 replicas 1GB/1CPU, PDB min 2, HPA min 3 max 10. No PCI data stored — all card data handled by Stripe/PayPal client-side.
Deploy: canary 5% via Istio, monitor error rate 30min before full rollout.
Health: /health (DB+Redis+Stripe API), /ready (webhooks registered), /metrics :4110.
