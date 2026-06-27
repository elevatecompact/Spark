# viewer-service — README

## Overview
The Viewer Service handles all viewer-side interactions on the Titan platform. It manages watch history, viewing preferences, content bookmarks, recommendations feedback signals, and viewer engagement metrics. This service optimizes the content consumption experience across both live and recorded media by tracking behavior and personalizing the interface.

## Purpose
Provide a personalized viewing experience by tracking viewer behavior, managing watch history with progress persistence, saving content preferences (categories, language, maturity filters), and collecting implicit feedback signals (likes, dislikes, watch time, completion rate) for the recommendation engine. It serves as the viewer's home base for all content discovery and consumption activities.

## Ownership
**Team:** Viewer Experience (eng-viewers@titan.dev)
**SLI:** 99.95% uptime, p99 preference load < 100ms
**Escalation:** #oncall-viewer
