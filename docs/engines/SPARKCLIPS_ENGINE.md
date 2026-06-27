# SparkClips Engine — AI Highlights

## Purpose

SparkClips is Titan's AI highlight generation engine. It automatically identifies and extracts the most engaging moments from videos and live streams, creating short-form clips optimized for social sharing and platform promotion.

## Architecture

SparkClips processes video through multimodal analysis: audio features (crowd reaction, commentary intensity), visual features (motion, scene changes, facial expressions), and engagement signals (chat spike, viewership surge during live streams).

## Tech Stack

- **Language**: Python (ML models), Go (orchestration)
- **Vision**: Fine-tuned VideoMAE for temporal action detection; face detection via MTCNN
- **Audio**: Speech activity detection, excitement classification via Wav2Vec 2.0
- **NLP**: Caption sentiment analysis and highlight-worthy phrase detection
- **Orchestration**: Temporal.io for durable workflow execution
- **Storage**: Nexus for source video access and clip output

## Key Features

- **Multi-modal analysis**: Combines visual, audio, and engagement signals for highlight detection
- **Auto-clip generation**: Configurable clip durations (15s, 30s, 60s) with smart start/end boundary detection
- **Engagement scoring**: Each clip is scored on predicted shareability and viewer retention
- **Real-time generation**: Near-real-time clip generation during live streams (30-second delay)
- **Customizable criteria**: Platform operators can tune sensitivity for different content categories
- **Template engine**: Pre-defined clip styles (goal celebration, funny moment, dramatic reveal)
- **Social export**: Direct publishing to social platforms via webhook integration
- **Feedback loop**: User engagement with generated clips trains the model for future selections

## Performance Targets

| Metric | Target |
|--------|--------|
| Clip generation latency (VOD) | < 3 minutes for 60-minute video |
| Clip generation latency (live) | < 30 seconds from moment to clip |
| Precision (user-shareable clips) | > 80% |
| Recall (capture all highlights) | > 70% |
| Processing throughput | 500 concurrent video jobs |