# Component Interaction Diagrams

This document provides text-based component interaction diagrams for major subsystems in the Spark platform.

## Diagram 1: Overall System Context

```

                         Internet
                              |
                              |
                 ▼                    ▼
        +----------------+  +----------------+
        |  Cloudflare CDN |  | Cloudflare WAF |
        |  (330 PoPs)     |  | (DDoS + Rules)  |
        +--------+-------+  +--------+-------+
                 |                    |
                 ▼                    ▼
        +-------------------------------------------+
        |           Envoy API Gateway                |
        |  Auth | Rate Limit | Routing | Tracing     |
        +---+-------+-------+-------+---------------+
            |       |       |       |
    +-------+       |       |       +--------------+
    ▼               ▼       ▼                      ▼
+--------+  +----------+  +--------+  +------------------+
|Identity|  |  Stream  |  | Wallet |  | Moderation       |
|Service |  | Service  |  |Service |  | Service          |
+---+----+  +----+-----+  +---+----+  +-------+----------+
    |            |            |                |
    +------------+------------+----------------+
                 |            |
                 ▼            ▼
        +------------------------------+
        |    Apache Kafka (Events)      |
        |  StreamStarted | GiftSent     |
        |  UserAction | ModerationEvent |
        +---+-------+-------+----------+
            |       |       |
            ▼       ▼       ▼
        +-----+ +----+ +--------+
        |Anal.| |Chat| |Discovery|
        |Svc  | |Svc | |Service  |
        +-----+ +----+ +--------+
```

## Diagram 2: Stream Ingestion Pipeline

```

+------------+    RTMP/WebRTC    +--------------+
| Broadcaster +------------------>+ Edge Ingest  |
+------------+                   +------+-------+
                                        |
                                   Transcoded
                                        |
                                        ▼
                              +-----------------------+
                              |  GPU Transcode Farm    |
                              |  1080p | 720p | 480p  |
                              +---------+-------------+
                                        |
                          +-------------+-------------+
                          |                           |
                          ▼                           ▼
                  +--------------+           +--------------+
                  | Nexus Engine |           |  Kafka Topic |
                  |  (S3 Store)  |           |SegmentsReady |
                  +------+-------+           +------+-------+
                         |                          |
                         ▼                          ▼
                  +--------------+           +--------------+
                  | CDN (Veloc.) |           | Stream Svc   |
                  | Edge Caches  +<----------+ (projection) |
                  +------+-------+           +--------------+
                         |
                         ▼
                  +--------------+
                  |   Viewer     |
                  |  (HLS/CMAF)  |
                  +--------------+
```

## Diagram 3: Multi-Tier Caching

```

+----------+   +----------+   +----------+   +----------+
| Browser  |-->| CDN Edge |-->| Velocity |-->| Redis    |
| Cache    |   | Cache    |   | Engine   |   | Cluster  |
+----------+   +----------+   +----------+   +----------+
     |              |              |              |
     | Miss         | Miss         | Miss         |
     ▼              ▼              ▼              ▼
+---------------------------------------------------------+
|                    Origin / Database                      |
+---------------------------------------------------------+
```

## Diagram 4: Authentication Flow

```

+--------+    +----------+    +----------+    +----------+
| Client |--->| Envoy    |--->| Identity |--->| Postgres |
|        |    | Gateway  |    | Service  |    | (Users)  |
+--------+    +----------+    +----------+    +----------+
     |              |               |               |
     |              |    +----------+               |
     |              |    | JWT Key  |               |
     |              |    | Store    |               |
     |              |    +----------+               |
     |              |               |               |
     |<-------------+---------------+---------------+
     |       JWT Token (Bearer)
     ▼
+--------+
| API    |--> Subsequent requests include JWT
| Calls  |    Envoy validates without DB call
+--------+
```
