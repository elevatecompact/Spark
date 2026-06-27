# stream-service — README

## Overview
The Stream Service is the core of Titan's live streaming infrastructure. It manages stream ingestion via RTMP, WHIP, and SRT protocols, GPU-accelerated transcoding into multiple resolutions, HLS and DASH packaging, edge-origin broadcast distribution, and automatic session recording for VOD archival. Built on a distributed media mesh for global scalability.

## Purpose
Handle the complete lifecycle of a live stream: session creation and configuration, protocol-aware ingestion with automatic bitrate adaptation, real-time transcoding into 720p, 1080p, and source renditions, segment packaging with 4-second HLS chunks for low latency, viewer broadcast via CDN edge nodes, and full session recording to S3 for later playback.

## Ownership
**Team:** Media Infrastructure (eng-media@titan.dev)
**SLI:** 99.99% uptime, end-to-end latency < 5s
**Escalation:** #oncall-stream
