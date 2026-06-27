# competition-service — API Contract
## Competitions: POST /v1/competitions, GET /v1/competitions/{id}, PATCH /v1/competitions/{id}, POST /v1/competitions/{id}/start, POST /v1/competitions/{id}/end
## Registration: POST /v1/competitions/{id}/register, POST /v1/competitions/{id}/withdraw, GET /v1/competitions/{id}/participants
## Brackets: GET /v1/competitions/{id}/bracket (tree structure), POST /v1/matches/{id}/score (submit), POST /v1/matches/{id}/confirm, POST /v1/matches/{id}/dispute
## Judging: POST /v1/competitions/{id}/judges (assign), POST /v1/submissions/{id}/score (judge scores), GET /v1/submissions/{id}/scores
## Leaderboard: GET /v1/competitions/{id}/leaderboard (live), GET /v1/competitions/{id}/results (final)
## Prizes: GET /v1/competitions/{id}/prizes, POST /v1/competitions/{id}/prizes/distribute (admin)
## Admin: POST /v1/admin/competitions/{id}/cancel, POST /v1/admin/matches/{id}/override
