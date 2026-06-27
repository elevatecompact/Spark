# chat-service — README

## Overview
The Chat Service powers all real-time messaging experiences on the Titan platform, including live stream chat, channel conversations, and transient event chats. Built on a WebSocket-first architecture with message persistence in PostgreSQL, real-time delivery via Redis pub/sub, moderation filtering, and rich media embedding. Designed for thousands of concurrent users in a single chat room.

## Purpose
Deliver low-latency, scalable chat for live streams with thousands of concurrent users. Supports emotes and badges for self-expression, real-time moderation actions (mute, ban, slow mode), scrollback history via cursor-paginated API, and tenant isolation between different chat rooms. Messages are scanned in real-time by the moderation service for policy violations.

## Ownership
**Team:** Real-Time Infrastructure (eng-realtime@titan.dev)
**SLI:** 99.99% uptime, p99 delivery < 100ms
**Escalation:** #oncall-chat
