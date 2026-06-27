# Scaling

Sparkclips scales by modality and stream. Each modality analyser runs as independent worker pool. Audio uses CPU workers; visual uses GPU workers; chat colocated with Echo. Redis job queue dispatches segments. Live streams get dedicated workers for duration. Clip rendering workers GPU-accelerated based on queue depth. Fusion layer stateless, combining signals when all modalities for time window complete.
