# competition-service — Runbook
## Alerts: BracketGenerationFailure, MatchDisputeRate > 5% (scoring issues), AutomaticProgressionFailure, PrizeDistributionFailure, CompetitionStartTimeout
## Override match: POST /v1/admin/matches/{id}/override {winnerId, reason}
## Cancel competition: POST /v1/admin/competitions/{id}/cancel — refunds entry fees, cancels all pending matches.
## Redistribute prizes: POST /v1/admin/competitions/{id}/prizes/distribute
## Rebuild bracket: POST /v1/admin/competitions/{id}/bracket/rebuild
