# licensing-service — API Contract
## Licenses: POST /v1/licenses, GET /v1/licenses/{id}, PATCH /v1/licenses/{id}, DELETE /v1/licenses/{id}, GET /v1/licenses (filter by rights holder or content)
## Rights: POST /v1/rights/content (register content rights), GET /v1/rights/{contentId} (check rights), POST /v1/rights/verify (check if content can be used in context)
## Royalties: GET /v1/royalties/calculate (projected), POST /v1/royalties/statement (generate period statement), GET /v1/royalties/statements, GET /v1/royalties/pending (unpaid amounts)
## Usage: POST /v1/usage/record (log content usage), GET /v1/usage/report (rights holder view), GET /v1/usage/content/{id} (usage history for content)
## Admin: POST /v1/admin/licenses/{id}/approve, POST /v1/admin/licenses/{id}/reject, POST /v1/admin/royalties/process (trigger payout), GET /v1/admin/compliance/report
