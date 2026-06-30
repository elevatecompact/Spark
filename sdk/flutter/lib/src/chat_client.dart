import 'models.dart';

class ChatClient {
  ChatClient(this._transport);
  final dynamic _transport;

  Future<List<ChatMessage>> list(String roomId, {int limit = 50, String? cursor}) {
    return _transport.jsonList<ChatMessage>(
      'GET',
      '/chat/rooms/$roomId/messages',
      (data) => ChatMessage.fromJson((data as Map).cast<String, dynamic>()),
      query: {
        'limit': '$limit',
        if (cursor != null) 'cursor': cursor,
      },
    );
  }

  Future<ChatMessage> send({required String roomId, required String body, String? replyToId}) {
    return _transport.json<ChatMessage>(
      'POST',
      '/chat/rooms/$roomId/messages',
      (data) => ChatMessage.fromJson((data as Map).cast<String, dynamic>()),
      body: {
        'body': body,
        if (replyToId != null) 'reply_to': replyToId,
      },
    );
  }

  Future<void> delete({required String roomId, required String messageId}) async {
    await _transport.request('DELETE', '/chat/rooms/$roomId/messages/$messageId');
  }

  Future<void> join(String roomId) async {
    await _transport.request('POST', '/chat/rooms/$roomId/join');
  }

  Future<void> leave(String roomId) async {
    await _transport.request('POST', '/chat/rooms/$roomId/leave');
  }

  Future<void> react({required String roomId, required String messageId, required String emoji}) async {
    await _transport.request('POST', '/chat/rooms/$roomId/messages/$messageId/reactions', body: {
      'emoji': emoji,
    });
  }
}
