# notification-service — Testing Guide
## Unit: Template rendering (Handlebars), preference evaluation logic, channel routing, rate limit enforcement, delivery status state machine.
## Integration: Full notification flow (event→route→template→send→deliver), push token registration, email deliverability, preference overrides.
## E2E: Verify actual delivery on real devices (push), email inbox (SendGrid sandbox), SMS (Twilio test numbers).
