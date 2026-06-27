# search-service — API Contract
## Search: GET /v1/search?q=keyword&type=all|creators|streams|recordings|clips&filters={}&sort=relevance|date|popularity&page=1&size=20
## Autocomplete: GET /v1/autocomplete?q=partial&size=10
## Indexing: POST /v1/index/{contentType} (index document), PUT /v1/index/{contentType}/{id} (update), DELETE /v1/index/{contentType}/{id} (remove), POST /v1/index/reindex (full reindex - admin)
## Admin: GET /v1/admin/stats (index sizes, query rates), POST /v1/admin/synonyms (manage synonym sets), PUT /v1/admin/analyzers (custom analyzers), GET /v1/admin/health (cluster health)
