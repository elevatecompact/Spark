# Troubleshooting

## Viewer cannot connect
1. Check ICE candidate gathering - verify STUN server is reachable.
2. Verify UDP ports 50000-50100 are open on SFU node.
3. Check the signalling WebSocket is connected; look for session.created event.

## Ingest stream drops
1. Check encoder bitrate is within published maximum.
2. Verify SRT passphrase matches config/pulse.toml.
3. Examine ingest logs for ingest.timeout messages.
4. Ensure TCP port 1935 (RTMP) or UDP port 9998 (SRT) is reachable from encoder.

## High latency
1. Check SFU CPU usage - scale out if > 80%.
2. Review bitrate ladder - lower the top rung if viewers buffer.
3. Enable sfu.congestion_control = pcc for better variable-network performance.
