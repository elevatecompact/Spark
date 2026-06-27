# Terminology

A shared glossary for the Spark platform.

| Term | Definition |
|------|-----------|
| **Creator** | A user who publishes content on Spark. |
| **Viewer** | A user who consumes content on Spark. |
| **Stream** | A live video broadcast, ingested via WebRTC or RTMP and delivered via HLS. |
| **Clip** | A short-form recorded video, typically under 60 seconds. |
| **VOD** | Video-on-demand — a pre-recorded, transcoded, and stored video asset. |
| **Channel** | A creator's dedicated page containing their profile, streams, clips, and VODs. |
| **Transcoding** | The process of converting source video into multiple resolutions and codecs for adaptive bitrate streaming. |
| **Adaptive Bitrate (ABR)** | Streaming technique that switches video quality based on the viewer's network conditions. |
| **SFU** | Selective Forwarding Unit — relays WebRTC streams to many viewers without transcoding. |
| **Event** | A domain-level occurrence published to Kafka (e.g., `StreamStarted`, `ContentUploaded`). |
| **Journey** | A real-time interaction session combining video, chat, reactions, and moderation. |
| **Treasure** | A virtual good or tip sent by a viewer to a creator during a stream. |
| **Curation** | A playlist or collection of content curated by a creator or algorithm. |
| **Moderation** | Tools and policies for managing chat, content, and user behavior. |
| **Edge** | Globally distributed points of presence for low-latency content delivery. |
| **Spark ID** | A verified identity system linking a creator's profile across platforms. |
