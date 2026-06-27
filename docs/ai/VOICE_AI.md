# Voice AI — Real-Time Speech

Aura is Spark's voice AI system, providing real-time text-to-speech and speech-to-text capabilities for live streaming, voiceovers, accessibility features, and interactive experiences.

## Text-to-Speech (TTS)

Aura's TTS engine uses a neural vocoder architecture (HiFi-GAN V2 + FastSpeech 2) with the following capabilities:

- **Voice Library**: 50+ pre-built voices across 20+ languages, with multiple age groups, accents, and speaking styles
- **Custom Voice Cloning**: Creators can train a custom voice model from as little as 30 minutes of recorded speech. Voice cloning uses a fine-tuned VITS model with speaker embedding extraction.
- **Emotion Control**: SSML-style tags for happiness, sadness, excitement, calm, and anger. Emotion-aware prosody modeling adjusts pitch, tempo, and energy.
- **Real-Time Streaming**: Streaming inference with <300ms first-audio latency. Supports chunked synthesis for live narration and interactive voice responses.

## Speech-to-Text (STT)

Aura's STT pipeline uses a fine-tuned Whisper large-v3 model optimized for streaming:

- **Real-Time Transcription**: Streaming ASR with word-level timestamps. Supports punctuation restoration and capitalization via a secondary transformer model.
- **Language Detection**: Automatic language identification with mid-sentence language switching support.
- **Speaker Diarization**: Identifies and labels different speakers in a conversation, critical for multi-participant streams and interviews.
- **Domain Adaptation**: Specialized language models for gaming terminology, creator slang, and platform-specific jargon.

## Latency Budget

Aura maintains strict streaming SLAs: TTS p95 < 500ms from text input to first audio byte, STT p95 < 800ms from audio input to final transcription output.
