# media-service — README

## Overview
The Media Service manages the storage, processing, and delivery of all media assets on the Titan platform. It handles file uploads, image transcoding and optimization, video transcoding, thumbnail generation, CDN content delivery, and DRM-protected playback. Every media file on the platform flows through this service.

## Purpose
Provide a unified media pipeline for all content types: images (avatars, banners, thumbnails), video (stream recordings, uploaded videos, clips), and audio (voice messages, music). Handles upload with resumable chunked uploads for large files, automatic transcoding into multiple resolutions and formats (HLS, DASH, MP4), thumbnail generation at configurable time intervals, content delivery via CDN with edge caching, and DRM encryption for premium content using Widevine and FairPlay.

## Ownership
**Team:** Media Infrastructure (eng-media@titan.dev)
**SLI:** 99.99% uptime, p99 upload < 1s
**Escalation:** #oncall-media
