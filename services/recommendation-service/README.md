# recommendation-service — README

## Overview
The Recommendation Service powers personalized content discovery on the Titan platform. It serves curated feeds for each user using collaborative filtering, content-based filtering, and real-time behavioral signals. The service operates a machine learning inference pipeline that combines offline-trained models with online feature computation for low-latency recommendations.

## Purpose
Deliver relevant content recommendations to viewers based on their watch history, search behavior, subscription patterns, and implicit engagement signals. Supports personalized home feeds, "up next" suggestions during streams, similar creators discovery, trending content ranking, and search result boosting. Models are trained offline using TensorFlow and served via ONNX runtime with sub-100ms inference latency.

## Ownership
**Team:** ML Platform (eng-ml@titan.dev)
**SLI:** 99.9% uptime, p99 inference < 100ms
**Escalation:** #oncall-recs
