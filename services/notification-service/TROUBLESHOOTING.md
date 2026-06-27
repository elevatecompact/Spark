# notification-service — Troubleshooting
## Push notifications not arriving: FCM/APNs certificate expired, device token stale, rate limited. Check FCM console, verify token last used, check rate limit counters.
## Emails going to spam: Sender reputation issue, missing DKIM/SPF, content flagged. Verify domain authentication in SendGrid, check spam reports.
## Digest not sent: CronJob failed, template rendering error, user preference excluded. Check CronJob logs, verify template syntax, test user preference evaluation.
## Bounced emails: Invalid address, mailbox full, server rejected. Process SendGrid webhook, update user email status, flag for re-verification.
