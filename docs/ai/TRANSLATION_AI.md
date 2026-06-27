# AI Translation — Polyglot

Polyglot is Spark's neural machine translation system, providing real-time multilingual content localization across 100+ languages with cultural and contextual awareness.

## Core Technology

Polyglot is built on a transformer-based sequence-to-sequence architecture with the following components:

- **Encoder**: Deep transformer encoder with language-agnostic embeddings. Supports mixed-language inputs and code-switching.
- **Decoder**: Autoregressive decoder with language-specific adapters for efficient fine-tuning without full retraining.
- **Attention**: Multi-head cross-lingual attention mechanism with alignment tracking for provenance.

## Key Features

- **Real-Time Streaming**: Sub-second translation latency for live chat, comments, and stream captions using inference optimization (int8 quantization, ONNX runtime, KV-cache batching)
- **Cultural Adaptation**: Beyond literal translation — adapts idioms, humor, references, and culturally sensitive content for the target audience. Uses a learned cultural adaptation layer trained on parallel corpora annotated for cultural context.
- **Domain Specialization**: Fine-tuned models for gaming, tutorials, news, entertainment, and educational content. Domain detection routes translations to the appropriate specialized model.
- **Preservation**: Retains emoji, formatting, hashtags, mentions, and embedded media links during translation.

## Workflow

1. Content is published on Spark
2. Language detection identifies source language
3. Translation request is published to Kafka with content ID and target locales
4. Polyglot generates translations using the appropriate domain model
5. Translations are stored alongside original content in the content store
6. Viewers receive localized content transparently based on their locale preferences
