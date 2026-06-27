# Runbook

## Startup
1. Verify model files: /models/whisper, /models/clip.
2. Start server: ./sparkclips serve.
3. Start modality workers: ./sparkclips worker --modality audio/visual.
4. Test: POST /v1/highlights/detect.
5. Check clip generation.

## Monitoring
- Dashboard: jobs in flight, queue depth, modality time, shareability scores.
- Alerts: failure rate > 5%, backlog > 100, GPU memory > 90%.
