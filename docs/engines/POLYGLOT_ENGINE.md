# Polyglot Engine — AI Translation

## Purpose

Polyglot is Titan's AI-powered translation engine. It provides real-time translation of user-generated content — video subtitles, chat messages, descriptions, and metadata — enabling a multilingual platform experience.

## Architecture

Polyglot uses a cascade of transformer models optimized for different content types. Short-form content (chat messages) uses lightweight distilled models for sub-50ms latency. Long-form content (video subtitles) uses full-size models for maximum translation quality.

## Tech Stack

- **Language**: Python (model inference), Go (API gateway, caching)
- **Models**: Fine-tuned NLLB (No Language Left Behind) for high-quality translation; distilled mBART for low-latency paths
- **Cache**: Redis for translation cache (keyed by content hash + source/target language)
- **Queue**: RabbitMQ for async translation of long-form content
- **Storage**: PostgreSQL for translation memory (reuse of previously translated segments)

## Key Features

- **50+ languages**: Full bidirectional translation support for Titan's content languages
- **Real-time chat translation**: Sub-50ms translation of chat messages with language auto-detection
- **Subtitle generation**: Batch translation of video subtitles with timing alignment preservation
- **Context-aware translation**: Sentence-level context consideration for pronoun resolution and idiomatic expressions
- **Translation memory**: Reuse of previously translated segments for consistency and cost savings
- **Custom glossaries**: Platform-defined terminology overrides (brand names, technical terms)
- **Auto-detection**: Automatic source language detection with confidence scoring
- **Profanity preservation**: Optional mode that preserves profanity patterns for moderation downstream

## Performance Targets

| Metric | Target |
|--------|--------|
| Chat translation latency | < 50ms (p99) |
| Subtitle translation (2-hour video) | < 5 minutes |
| Translation quality (BLEU score) | > 35 for major language pairs |
| Cache hit ratio | > 70% for repeated content |
| Supported language pairs | 2,500+ |