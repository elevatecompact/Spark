# advertising-service — API Contract
## Campaigns: POST /v1/campaigns, GET /v1/campaigns/{id}, PATCH, DELETE, GET /v1/campaigns (advertiser's campaigns)
## Ad Units: POST /v1/ad-units (creative), GET /v1/ad-units/{id}, PATCH, DELETE, POST /v1/ad-units/{id}/approve (admin)
## Ad Serving: GET /v1/ads/request (SSP: return best ad for placement), POST /v1/ads/impression (record impression), POST /v1/ads/click (record click)
## Targeting: GET /v1/targeting/options (available targeting dimensions), POST /v1/targeting/segments (create custom audience)
## Analytics: GET /v1/analytics/campaigns/{id}/performance, GET /v1/analytics/creators/{id}/revenue, GET /v1/analytics/realtime (live ad metrics)
## Inventory: GET /v1/inventory/available (available ad slots), POST /v1/inventory/pricing (floor price recommendations)
## Admin: POST /v1/admin/campaigns/{id}/pause, POST /v1/admin/campaigns/{id}/resume, GET /v1/admin/revenue
