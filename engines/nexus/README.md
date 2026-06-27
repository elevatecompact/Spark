# Nexus Engine

**Purpose:** Distributed media storage and asset management engine.
**Tech Stack:** Go, S3-compatible storage, PostgreSQL, Redis, gRPC, FFmpeg, ImageMagick.

Nexus stores, organises, and serves media assets - videos, images, thumbnails, audio. Content-addressed storage with deduplication, lifecycle policies, and origin protection. Supports hot/warm/cold storage tiers with automated migration based on access patterns.
