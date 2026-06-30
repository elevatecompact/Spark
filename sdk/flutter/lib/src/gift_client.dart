import 'models.dart';

class GiftClient {
  GiftClient(this._transport);
  final dynamic _transport;

  Future<Gift> send({required String toUserId, required int amountCents, String message = '', String? streamId}) {
    return _transport.json<Gift>(
      'POST',
      '/gifts',
      (data) => Gift.fromJson((data as Map).cast<String, dynamic>()),
      body: {
        'to_user_id': toUserId,
        'amount_cents': amountCents,
        'message': message,
        if (streamId != null) 'stream_id': streamId,
      },
    );
  }

  Future<List<Gift>> received({int limit = 50, String? cursor}) {
    return _transport.jsonList<Gift>(
      'GET',
      '/gifts/received',
      (data) => Gift.fromJson((data as Map).cast<String, dynamic>()),
      query: {
        'limit': '$limit',
        if (cursor != null) 'cursor': cursor,
      },
    );
  }

  Future<List<Gift>> sent({int limit = 50, String? cursor}) {
    return _transport.jsonList<Gift>(
      'GET',
      '/gifts/sent',
      (data) => Gift.fromJson((data as Map).cast<String, dynamic>()),
      query: {
        'limit': '$limit',
        if (cursor != null) 'cursor': cursor,
      },
    );
  }

  Future<List<Gift>> forStream(String streamId, {int limit = 50, String? cursor}) {
    return _transport.jsonList<Gift>(
      'GET',
      '/streams/$streamId/gifts',
      (data) => Gift.fromJson((data as Map).cast<String, dynamic>()),
      query: {
        'limit': '$limit',
        if (cursor != null) 'cursor': cursor,
      },
    );
  }
}
