# AI-Powered Content Moderation — Sentinel

Sentinel is Spark's real-time content moderation engine designed to detect and mitigate harmful content at platform scale while minimizing false positives and respecting free expression.

## Detection Capabilities

Sentinel operates across multiple content modalities:
- **Text**: Hate speech, harassment, violence incitement, spam, phishing, and policy-violating language using transformer-based NLP models (fine-tuned BERT/RoBERTa)
- **Image**: NSFW detection, graphic violence, weapon recognition, and policy-violating visual content via ConvNeXt and ViT architectures
- **Video**: Frame-by-frame analysis with temporal coherence checking to catch context-dependent violations
- **Audio**: Speech-to-text transcription followed by text moderation for spoken content in live streams and recordings

## Architecture

Sentinel follows a two-stage pipeline:

1. **Fast Filter**: Lightweight ONNX-optimized models run on every upload with sub-100ms latency. This catches obvious violations and passes uncertain cases to the deep inspector.
2. **Deep Inspector**: Full transformer ensemble with cross-modal fusion runs asynchronously. Uses attention-based reasoning to understand context, sarcasm, and cultural nuance.

## Human-in-the-Loop

All automated moderation decisions are subject to appeal. Edge cases, borderline content, and appeals are routed to human moderators via a specialized review queue. Moderators see the model's reasoning, confidence scores, and flagged regions. Feedback from human review continuously fine-tunes the model through active learning loops.

## Platform Integration

Sentinel is integrated at every content ingestion point: uploads, live streams, comments, DMs, and profile content. Violations trigger immediate action (block, flag, or blur) and produce an evidence package for the moderation team.
