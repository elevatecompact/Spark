# Troubleshooting

## High false positive rate
1. Check active learning feedback loop processing human overrides.
2. Review recent blocklist updates for overbroad patterns.
3. Compare model accuracy: GET /v1/model/status.
4. Adjust tier2.confidence_threshold upward.

## Slow video moderation
1. Increase video.sample_interval for longer videos.
2. Ensure GPU nodes have sufficient VRAM (> 8GB).
3. Check video download bandwidth bottleneck.
4. Scale tier 2 GPU nodes.
