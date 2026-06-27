# Forge Engine

**Purpose:** High-performance video transcoding and packaging engine.
**Tech Stack:** Rust, FFmpeg bindings, NVIDIA NVENC, VAAPI, CMAF, HLS, DASH.

Forge transforms raw ingested video into adaptive bitrate ladders with optimized encoding. It supports hardware-accelerated encoding on NVIDIA and AMD GPUs, produces CMAF fragments for HLS/DASH packaging, and applies per-title encoding optimisation. Designed to replace traditional FFmpeg pipelines with a safer, faster Rust-native core.
