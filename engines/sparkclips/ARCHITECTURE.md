# Architecture

Sparkclips uses signal fusion architecture. Parallel analysers process: Whisper for speech/text excitement, CLIP for visual event classification, audio spectrogram for crowd reaction, and chat sentiment for engagement spikes. Fusion layer combines signals with configurable weights to produce excitement scores. Peaks become clip candidates. Clipping pipeline extracts segments, applies transitions, renders via FFmpeg. Ranking model sorts clips by predicted shareability.
