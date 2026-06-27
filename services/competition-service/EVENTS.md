# competition-service — Event Contracts
## Published: competition.created, competition.started, competition.ended, competition.participant.registered, competition.match.scheduled, competition.match.completed, competition.match.disputed, competition.prize.distributed, competition.leaderboard.updated
## Consumed: wallet.transaction.settled (prize distribution), notification.push.sent (match reminder), stream.session.started (competition stream going live)
## Schema: CompetitionMatchCompletedEvent {competitionId, matchId, winnerId, loserId, scores, round, bracketPosition, completedAt}
