# event-service — API Contract
## Events: POST /v1/events, GET /v1/events/{id}, PATCH /v1/events/{id}, DELETE /v1/events/{id}, GET /v1/events (discover, filterable by date/category/creator)
## Ticketing: POST /v1/events/{id}/tickets (create tier), GET /v1/events/{id}/tickets, POST /v1/events/{id}/rsvp (free events), POST /v1/tickets/{id}/purchase, GET /v1/tickets/{id}
## Schedule: GET /v1/events/{id}/schedule, POST /v1/events/{id}/schedule/sessions, PATCH /v1/sessions/{id}
## Series: POST /v1/series (recurring events), GET /v1/series/{id}, PATCH /v1/series/{id}, DELETE /v1/series/{id}
## Admin: POST /v1/admin/events/{id}/cancel, POST /v1/admin/events/{id}/refund-all, GET /v1/admin/events/stats
