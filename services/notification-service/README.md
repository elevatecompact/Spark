# notification-service — README

## Overview
The Notification Service manages all user-facing communications across the Titan platform. It delivers notifications through multiple channels: push notifications (FCM/APNs), email (SendGrid), SMS (Twilio), and in-app notification center. Handles template rendering, user notification preferences, delivery scheduling, and read/unread tracking.

## Purpose
Provide a unified notification delivery system that ensures users receive timely, relevant alerts about platform activity. Supports event-driven notifications (new subscriber, gift received, stream started), digest emails (daily/weekly summaries), and transactional messages (password reset, payment confirmation). Each user controls channel preferences and notification frequency per notification type through the preferences API.

## Ownership
**Team:** User Engagement (eng-engagement@titan.dev)
**SLI:** 99.95% uptime, p99 delivery < 5s
**Escalation:** #oncall-notifications
