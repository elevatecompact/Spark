# Competitions

Competitions provide a framework for time-bound challenges that drive viewer engagement, creator participation, and community excitement.

## Competition Types

**Viewer Challenges.** Time-limited challenges where viewers earn points or badges by completing specific actions such as watching a certain number of hours, sending gifts, or sharing content. Leaderboards show top participants with rewards for top finishers.

**Creator Competitions.** Competitive events where creators compete on specific metrics like follower growth, engagement rate, or content output over a defined period. Winners receive platform promotion, monetization bonuses, and physical awards.

**Team Battles.** Creators or viewers form teams to compete in collaborative challenges. Team scores are aggregated from individual member contributions. Team-based competitions foster community building and cross-promotion.

**Seasonal Events.** Recurring competitions tied to holidays, seasons, or platform milestones. These events feature themed rewards, exclusive items, and special leaderboards.

## Competition Engine

The competition engine manages competition lifecycle including creation, participant registration, score tracking, leaderboard updates, and reward distribution. Scoring rules are configurable per competition type. Real-time leaderboard updates use the WebSocket system to show current standings. Score verification prevents manipulation through rate limiting and anomaly detection.

## Rewards

Competition rewards include platform currency, exclusive badges and emotes, profile customization items, subscription credits, and promotional placement on the platform. Rewards are distributed automatically at competition conclusion through the wallet and inventory systems.

## Anti-Abuse

Competition participation is monitored for fraudulent activity including bot accounts, coordinated manipulation, and terms of service violations. Suspicious participants are removed and their scores are voided.
