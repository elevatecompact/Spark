# WebSocket API

WebSocket connections provide real-time, bidirectional communication between clients and SPARK servers. They are used for live notifications, chat messaging, real-time content updates, and collaborative features.

## Connection Lifecycle

Clients establish a WebSocket connection by upgrading an HTTP request to wss://ws.sparkplatform.com. The initial handshake includes authentication credentials in the query string. Once connected, the server assigns a connectionId and sends a connection_ack message. Heartbeat ping/pong messages at 30-second intervals detect stale connections.

## Message Format

All WebSocket messages follow a JSON format with a type field indicating the message category, a payload field containing the message body, and an optional id field for request-response correlation. Message types include subscribe, unsubscribe, notification, event, error, and ack. Subscriptions follow the channel naming pattern {resource}:{event}:{scope} such as content:updated:user_123.

## Channels

System channels broadcast platform-wide events such as maintenance announcements. User channels deliver personalized notifications and updates. Content channels stream updates for specific content items such as live stream comments and viewer counts. Admin channels carry moderation events and system alerts.

## Reconnection

Clients implement exponential backoff reconnection starting at 1 second with a maximum delay of 60 seconds. The server supports state recovery through the last known eventId, allowing clients to reconnect without missing messages. Duplicate connections for the same user are detected and the older connection is closed.
