# Sentinel Engine — Content Moderation

## Purpose

Sentinel is Titan's AI-powered content moderation engine. It scans all user-generated content — videos, images, text, and audio — for policy violations, toxic behavior, and prohibited content before it reaches the platform audience.

## Architecture

Sentinel uses a multi-stage pipeline: streaming ingestion for real-time chat moderation, batch processing for video/image analysis, and human-in-the-loop review for borderline cases.

## Tech Stack

- **Language**: Python (ML inference), Go (API gateway)
- **Vision Models**: Fine-tuned CLIP and YOLO for image/video classification, OCR for text extraction
- **NLP Models**: RoBERTa-based classifier for toxicity detection, hate speech detection, spam detection
- **Audio**: Speech-to-text via Whisper, then NLP classification on transcripts
- **Queue**: Kafka for streaming moderation pipeline
- **Storage**: PostgreSQL for decisions and appeals

## Key Features

- **Real-time chat moderation**: Streaming NLP inference on every chat message with < 50ms latency
- **Video moderation**: Keyframe extraction to image classification pipeline for adult content, violence, gore
- **Audio moderation**: Whisper transcription to NLP classification for audio content
- **Text analysis**: Multi-language toxic content detection with contextual understanding
- **Appeal system**: Users can appeal automated decisions with human review workflow
- **Custom policies**: Platform-defined moderation rules per region and content type
- **Confidence scoring**: Each decision includes a confidence score; low-confidence results are queued for human review

## Performance Targets

| Metric | Target |
|--------|--------|
| Chat message latency | < 50ms (p99) |
| Video moderation time | < 30 seconds for 10-minute video |
| Accuracy (precision) | > 95% |
| False positive rate | < 1% |
| Human review queue time | < 5 minutes |