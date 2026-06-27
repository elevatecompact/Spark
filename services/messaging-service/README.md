# messaging-service — README

## Overview
The Messaging Service enables direct person-to-person and group messaging across the Titan platform. Unlike the chat service which is room/stream-oriented, this service focuses on persistent threaded conversations with read receipts, typing indicators, media sharing, message reactions, and optional end-to-end encryption for private conversations.

## Purpose
Provide a reliable, low-latency messaging experience for direct conversations between users. Supports text, images, voice messages, file attachments up to 100MB, emoji reactions, read state synchronization across multiple devices, and typing indicator broadcast. Designed for both one-on-one DMs and group conversations of up to 500 participants with admin roles.

## Ownership
**Team:** Real-Time Infrastructure (eng-realtime@titan.dev)
**SLI:** 99.99% delivery rate, p99 delivery < 200ms
**Escalation:** #oncall-messaging
