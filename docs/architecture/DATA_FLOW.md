# End-to-End Data Flow

This document describes the primary data flows through the Spark platform, tracing data from origin to consumption across multiple subsystems.

## Flow 1: Live Stream Broadcast

`
Broadcaster -> RTMP/WebRTC -> Edge Ingest -> Transcoding Farm
                                                              |
                                    +-------------------------+
                                    |
                                    v
            +----------------+     +------------------+
            | Nexus / S3     |     |  Kafka Topics    |
            +----------------+     +------------------+
                    |                       |
                    v                       v
            +----------------+     +------------------+
            | CDN Cache      |     | Stream Service   |
            +----------------+     +------------------+
                    |                       |
                    v                       v
            +----------------+     +------------------+
            | Viewer (HLS)   |     | Analytics Svc    |
            +----------------+     +------------------+
`

Broadcaster connects via RTMP or WebRTC to the nearest edge ingest point. Video frames are forwarded to the transcoding farm for bitrate ladder generation. Transcoded segments are stored in Nexus/S3 object storage. Segment availability events are published to Kafka. CDN edge caches pull segments on first viewer request. Viewers receive HLS/CMAF segments from the nearest CDN edge. WebRTC viewers connect through edge SFUs for sub-second latency.

## Flow 2: Virtual Gift Transaction

`
Viewer -> API Gateway -> Wallet Service
                              |
                    Validate Balance
                    Deduct Balance
                    Publish: GiftSent
                              |
                              v
                        Kafka Topic
                              |
            +-----------------+------------------+
            |                 |                  |
            v                 v                  v
    Chat Service      Stream Service       Analytics
            |                 |                  |
            v                 v                  v
   WebRTC Data Channel  Leaderboard Update   Revenue Calc
            |
            v
   Viewer sees animation
`

Viewer sends gift via API. Wallet Service validates sufficient balance. Balance is deducted atomically. GiftSent event published to Kafka. Chat Service receives event and broadcasts gift animation via WebRTC data channel. Stream Service updates gift leaderboard. Revenue is calculated for broadcaster share.

## Flow 3: Content Discovery

`
Viewer searches -> API Gateway -> Discovery Service -> OpenSearch
                                                                |
                                                           Results cached
                                                           in Redis
                                                                |
                                          +---------------------+
                                          |
                                          v
                              Results returned to viewer
                              (title, viewer count, thumbnail,
                               broadcaster, tags, language)
                                          |
                                          v
                              Viewer selects stream
                              Stream Service provides CDN URL
                              Player connects to nearest edge
`

Discovery Service handles search queries, executes OpenSearch full-text search with personalization boost. Results cache in Redis. Player receives CDN playback URLs, connects to nearest edge cache for segment delivery.
