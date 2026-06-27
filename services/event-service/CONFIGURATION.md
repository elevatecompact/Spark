# event-service — Configuration
EVENT_PORT=4020, EVENT_DB_URL, EVENT_REDIS_URL, EVENT_KAFKA_BROKERS, MAX_EVENTS_PER_CREATOR=50 (active), MAX_TICKETS_PER_EVENT=5 (tiers), EVENT_DISCOVERY_PAGE_SIZE=50, RSVP_GRACE_PERIOD_HOURS=2 (cancellation window), TICKET_SALES_END_HOURS=1 (before start)
FF: events_enabled=true, ticketed_events=true, recurring_series=true, hybrid_events=false, event_discovery=true
Rate limits: 10 events/month per creator, 100 RSVPs/day per user, 5 ticket tiers per event
