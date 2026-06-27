# moderation-service — README

## Overview
The Moderation Service is the content safety and policy enforcement engine for the Titan platform. It provides automated content scanning using ML-based image and text classifiers, configurable rule-based filtering, human review queue management, and policy violation enforcement actions. Every piece of user-generated content passes through this service.

## Purpose
Keep the Titan platform safe by automatically detecting and acting on policy violations across all content types: chat messages, stream video/audio, profile content, uploaded media, community posts, and reported content. Supports configurable policy rules with different severity levels (warn, restrict, remove, suspend). Human moderators review flagged content through a dashboard with efficient review workflows and decision consistency tools.

## Ownership
**Team:** Trust & Safety (eng-trust@titan.dev)
**SLI:** 99.95% uptime, p99 scan < 500ms
**Escalation:** #oncall-moderation
