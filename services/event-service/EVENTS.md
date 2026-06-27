# event-service — Event Contracts
## Published: event.created, event.updated, event.cancelled, event.started (stream begun), event.ended, event.ticket.purchased, event.rsvp.confirmed, event.series.occurrence.created
## Consumed: creator.channel.created (announce upcoming events), wallet.transaction.settled (confirm ticket purchase), stream.session.started (mark event as live), notification.push.sent (event reminder confirmation)
## Schema: EventStartedEvent {eventId, creatorId, streamId, sessionId, ticketCount, concurrentViewers, startedAt}
