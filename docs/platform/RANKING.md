# Ranking

Ranking systems determine the ordering of content across various surfaces on the SPARK platform, including feeds, search results, recommendations, and trending lists.

## Ranking Models

SPARK uses machine learning models to rank content based on predicted user engagement. The primary ranking model predicts watch time probability, which correlates strongly with user satisfaction and retention. Secondary models predict like probability, share probability, follow probability, and report probability for quality control. These predictions are combined into a single ranking score using a weighted ensemble approach.

## Feature Engineering

Ranking features are organized into several categories. Content features include content age, duration, category, language, and content quality score. User features include watch history, interaction patterns, session context, and user preferences. Context features include time of day, device type, and previous session engagement. Social features include creator relationship, friend activity, and network engagement. Real-time features include trending velocity, current viewer count, and recent engagement spikes.

## Model Training

Ranking models are trained on historical interaction data using gradient boosted decision trees. Training data includes impressions with corresponding user actions as positive and negative labels. Models are trained daily using the previous 30 days of data. Online evaluation tracks offline metric correlation with online performance through A/B testing.

## Serving

Models are served through a low-latency inference pipeline. Feature computation runs as a streaming pipeline with sub-second latency. Model inference targets 10ms p99 latency per request. Ranking results are cached with short TTLs for frequently requested surfaces. Fallback ranking using heuristic rules ensures service availability during model failures.
