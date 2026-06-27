# Troubleshooting

## Job stuck in queue
1. Check worker count: GET /v1/health workers_available.
2. Verify worker can reach queue: check worker logs for queue.connected.
3. Ensure GPU encoder sessions not exhausted: GET /v1/gpu.
4. Check job does not exceed encoder.max_concurrent_sessions.

## Encoding failures
1. Check input media validity: ffprobe inputUri.
2. Verify GPU drivers and encoder firmware loaded.
3. Check for driver timeout: dmesg | grep nvidia.
4. Reduce encoder.target_usage for more reliable encoding.
