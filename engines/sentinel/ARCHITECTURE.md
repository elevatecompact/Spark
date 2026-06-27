# Architecture

Sentinel uses a tiered classification pipeline. Tier 1 runs regex and blocklist matching in-memory for sub-millisecond filtering. Tier 2 loads lightweight DistilBERT and ResNet-18 ONNX models. Tier 3 uses larger models (RoBERTa, CLIP) for ambiguous content, gated by a budget-aware router. Video moderation samples keyframes at configurable intervals. All decisions logged with model confidence scores for audit and active learning sampling.
