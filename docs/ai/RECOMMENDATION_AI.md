# AI Recommendation Engine — Oracle

Oracle is Spark's multi-modal recommendation engine, delivering personalized content feeds that maximize engagement, discovery, and creator sustainability.

## Architecture

Oracle implements a hybrid recommendation approach combining multiple paradigms:

- **Collaborative Filtering**: Two-tower neural network using user-video interaction embeddings. Trained on implicit feedback signals (watch time, completion rate, likes, shares, skips).
- **Content-Based Filtering**: Video embeddings derived from visual features, audio features, transcript topics, and metadata. Uses CLIP-style multi-modal encoders.
- **Knowledge Graph**: Entity-aware recommendations leveraging a graph of users, creators, topics, and trends. Captures serendipitous discovery paths.
- **Reinforcement Learning**: Contextual bandit framework for exploration-exploitation balance. Dynamically adjusts the trade-off based on user session state and novelty tolerance.

## Feature Pipeline

The recommendation pipeline processes features in real time:
- **User Features**: Watch history, search queries, session context, device type, time-of-day, location
- **Video Features**: Content embeddings, upload recency, creator popularity, topic tags, engagement velocity
- **Context Features**: Trending signals, seasonal events, platform-wide viral coefficients

## Serving

Oracle serves recommendations via gRPC with strict latency budgets (p99 < 50ms). Precomputed candidate pools are refreshed every 5 minutes via Spark jobs. Real-time reranking uses a lightweight transformer that scores ~500 candidates per request.

## Fairness and Diversity

Oracle incorporates diversity constraints to avoid filter bubbles. Recommendation mixes include exploration slots, creator diversity targets, and content-type variety constraints. Auditing dashboards track representation metrics across demographics.
