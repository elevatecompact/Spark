# Events

## Published Events
- stream.started - Payload: { streamId, ingestProtocol, resolution, bitrate, timestamp } - Emitted when an encoder begins pushing.
- stream.stopped - Payload: { streamId, duration, bytesIngested, reason } - Emitted on stream end or encoder disconnect.
- viewer.joined - Payload: { sessionId, streamId, viewerIp, userAgent }.
- viewer.left - Payload: { sessionId, streamId, watchDuration, reason }.
- stream.quality.changed - Payload: { streamId, previousBitrate, newBitrate, reason } - Emitted when adaptive bitrate switch occurs.

## Subscribed Events
- stream.recording.request - Request to start/stop recording a stream.
- stream.transcode.request - Request to adjust transcoding profile mid-stream.
