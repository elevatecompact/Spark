# Sentinel Engine

**Purpose:** AI-powered content moderation engine for text, image, and video.
**Tech Stack:** Python, PyTorch, Transformers, CLIP, Faster R-CNN, ONNX, gRPC, Redis.

Sentinel scans user-generated content for policy violations using a multi-stage pipeline: fast blocklist matching, lightweight classifier models, and deep multimodal analysis for edge cases. Models are fine-tuned on platform-specific data with continuous active learning.
