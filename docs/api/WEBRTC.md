# WebRTC API

WebRTC powers SPARK's peer-to-peer media streaming for live broadcasting, video calls, and real-time audio communication. It enables low-latency direct connections between users without media server bottlenecks.

## Signaling

Signaling for WebRTC connections is handled through the WebSocket API. The signaling flow begins with an offer from the initiating peer containing session description protocol data. The receiving peer responds with an answer. ICE candidates are exchanged bidirectionally as they are discovered. The signaling server relays these messages but does not process media data.

## STUN and TURN

A STUN server is deployed to help peers discover their public IP addresses and ports for direct connections. TURN servers are available as a fallback when direct peer-to-peer connections fail due to NAT traversal issues. TURN relay is used only when necessary to minimize bandwidth costs. Multiple TURN servers are deployed across global regions for low-latency relay.

## Media Tracks

Each WebRTC connection supports multiple simultaneous media tracks: video tracks up to 1080p at 60fps, audio tracks with Opus codec at 48kHz, and data channels for metadata such as chat messages and stream events. Adaptive bitrate encoding adjusts video quality based on network conditions using the SDP offer-answer model.

## Simulcast and SVC

Simulcast sends multiple resolution layers of the same video stream, allowing viewers to receive the appropriate quality for their connection. Scalable Video Coding provides temporal and spatial scalability for adaptive streaming. The SFU selective forwarding unit routes the appropriate layer to each viewer.
