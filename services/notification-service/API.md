# notification-service — API Contract
## Notifications: GET /v1/notifications (inbox), PATCH /v1/notifications/{id}/read, POST /v1/notifications/read-all, DELETE /v1/notifications/{id}
## Preferences: GET /v1/preferences (channel per type), PATCH /v1/preferences
## Templates: GET /v1/templates, POST /v1/templates (admin), PATCH /v1/templates/{id} (admin)
## Send: POST /v1/send (internal API for services), POST /v1/send/batch
## Devices: POST /v1/devices (register push token), DELETE /v1/devices/{id}
## Admin: POST /v1/admin/test-push {userId}, POST /v1/admin/test-email {email}, GET /v1/admin/delivery-stats
