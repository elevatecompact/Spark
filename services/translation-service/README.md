# translation-service — README

## Overview
The Translation Service provides automated content localization across the Titan platform. It translates stream titles, descriptions, chat messages, UI labels, and community content into multiple languages. Uses a combination of machine translation (DeepL, Google Translate) and human-reviewed translation memory for high-quality localization.

## Purpose
Break language barriers on the platform by automatically translating content into viewers' preferred languages. Supports 40+ languages with automatic language detection for incoming content. Translations are cached with translation memory to reduce API costs and ensure consistency. Human translators can review and correct automated translations through the review queue. Supports real-time chat translation for multilingual streams.

## Ownership
**Team:** Platform Engineering (eng-platform@titan.dev)
**SLI:** 99.9% uptime, p99 translate < 2s
**Escalation:** #oncall-translation
