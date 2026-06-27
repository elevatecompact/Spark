# event-service — Testing Guide
## Unit: Event status state machine, ticket inventory management (prevent oversell), RSVP capacity enforcement, event series recurrence calculation (daily/weekly/monthly).
## Integration: Full event lifecycle (create→publish→ticket sales→start→complete), recurring series generation, ticket purchase flow with wallet, event cancellation + refund.
## Load: 1000 concurrent ticket purchases, 5000 RSVPs on single event, series generation for 1000 series x 12 months.
