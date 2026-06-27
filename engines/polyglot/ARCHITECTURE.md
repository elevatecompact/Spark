# Architecture

Polyglot uses a pipeline: language detection to preprocessing to translation to quality estimation to postprocessing. Language detection runs M2M100 LID for 200+ languages. Translation loads HuggingFace models (NLLB-200, Opus-MT) exported to ONNX for inference. Term management applies customer glossaries. Quality estimation (CometKiwi) scores translations; low-confidence results flagged for human review. Streaming operates on sentence segments with context window carryover.
