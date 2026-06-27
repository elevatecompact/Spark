# Scaling

Sentinel scales with tier separation. Tier 1 runs on CPU-only nodes. Tier 2 on GPU nodes with model batching. Tier 3 on dedicated GPU nodes with large VRAM for CLIP and RoBERTa. Budget-aware router distributes traffic across tiers. Redis-backed queue handles async video jobs. Models deploy via blue-green to avoid inference downtime.
