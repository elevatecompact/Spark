# Troubleshooting

## No highlights detected
1. Check modality config - fusion.excitement_window appropriate for content.
2. Review per-modality scores in Redis.
3. Lower audio excitement_threshold for quiet content.
4. Ensure video duration > fusion.excitement_window.

## Clip generation failing
1. Verify FFmpeg with required codecs.
2. Check output format supported.
3. Ensure disk space for temp rendering (3x clip size).
4. Verify source video still accessible.
