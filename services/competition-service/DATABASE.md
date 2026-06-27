# competition-service — Database Schema
## competitions: id UUID PK, creator_id FK, title, description, format(elimination,score,leaderboard), bracket_type(single,double), status(draft,open,started,completed,cancelled), max_participants INT, registration_start/end, competition_start/end, prize_pool_cents, rules TEXT[]
## competition_participants: id UUID PK, competition_id FK, user_id FK, seed_number INT, status(registered,active,eliminated,withdrawn,winner), bracket_position, registered_at
## matches: id UUID PK, competition_id FK, round INT, bracket_position, participant1_id FK, participant2_id FK, winner_id FK nullable, scores JSONB, status(scheduled,completed,disputed), scheduled_at, completed_at
## judging_submissions: id UUID PK, competition_id FK, participant_id FK, scores JSONB {judgeId: score}, average_score, status(pending,scored,confirmed)
## Redis: Live leaderboard (sorted set by score), bracket cache (tree structure), match status cache
