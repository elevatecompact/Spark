# Configuration

Configuration is loaded from config/pulse.toml or environment variables prefixed PULSE_.

| Key | Default | Description |
|-----|---------|-------------|
| ingest.rtmp.port | 1935 | RTMP listen port |
| ingest.srt.port | 9998 | SRT listen port |
| webrtc.ice.port_range | 50000-50100 | UDP port range for ICE candidates |
| webrtc.stun_servers | stun:stun.l.google.com:19302 | STUN server list |
| sfu.max_viewers | 1000 | Max viewers per SFU node |
| sfu.congestion_control | gcc | GCC or PCC congestion control algorithm |
| recording.enabled | false | Enable automatic stream recording to S3 |
| recording.segment_duration | 6 | HLS segment duration in seconds |
