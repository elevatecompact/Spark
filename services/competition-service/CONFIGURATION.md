# competition-service — Configuration
COMPETITION_PORT=4021, COMPETITION_DB_URL, COMPETITION_REDIS_URL, COMPETITION_KAFKA_BROKERS, MAX_COMPETITIONS_PER_CREATOR=5 (active), MAX_PARTICIPANTS_DEFAULT=128, MIN_PARTICIPANTS=4, AUTO_PROGRESS_DELAY_MINUTES=30 (after match ends), JUDGE_SCORE_RANGE=1-100, PRIZE_DISTRIBUTION_CRON="0 0 * * *" (daily)
FF: competitions_enabled=true, bracket_competitions=true, score_competitions=true, leaderboard_competitions=true, judging_panels=true, prize_pools_enabled=true
Rate limits: 5 competitions/month per creator, 10 registrations/day per user, 5 judge assignments per competition
