# Forge Engine — Video Transcoding

## Purpose

Forge is Titan's video transcoding engine. It transforms raw uploaded video files into optimized, multi-resolution renditions for streaming delivery. Forge handles the full media processing pipeline from ingestion through packaging.

## Architecture

Forge uses a distributed job queue architecture. Uploaded videos are split into segments, each segment is transcoded independently on GPU-accelerated workers, and results are assembled and packaged into HLS and MPEG-DASH formats.

## Tech Stack

- **Language**: Go (orchestration), C++ (FFmpeg bindings)
- **Media Processing**: FFmpeg with NVIDIA NVENC for GPU-accelerated encoding
- **Job Queue**: RabbitMQ with priority queues
- **Storage**: Nexus (S3-compatible) for source and output
- **Packaging**: Shaka Packager for HLS/DASH output
- **Deployment**: Kubernetes with GPU node pools (A10G, A100)

## Key Features

- **Adaptive encoding**: Per-title encoding optimizes bitrate based on content complexity
- **Resolution ladder**: Automatic generation of 1080p, 720p, 480p, 360p renditions
- **GPU acceleration**: Hardware encoding via NVENC (10x speedup vs CPU)
- **Segment-level parallelism**: Independent segment encoding with automatic reassembly
- **Thumbnail generation**: Keyframe extraction at configurable intervals
- **Quality metrics**: SSIM and VMAF computation for quality assurance
- **Failure recovery**: Segments are retried independently; partial failures do not require full re-encode

## Performance Targets

| Metric | Target |
|--------|--------|
| Transcoding speed (1080p) | < 30 seconds per minute of video (GPU) |
| Concurrent jobs per cluster | 500+ |
| Quality (VMAF) | > 90 for highest rendition |
| Job failure rate | < 0.1% |
| Time to first byte delivered | < 5 seconds after upload complete |