import 'models.dart';

class ContentClient {
  ContentClient(this._transport);
  final dynamic _transport;

  Future<ContentItem> get(String id) {
    return _transport.json<ContentItem>('GET', '/content/$id',
        (data) => ContentItem.fromJson((data as Map).cast<String, dynamic>()));
  }

  Future<ContentItem> create(Map<String, dynamic> payload) {
    return _transport.json<ContentItem>('POST', '/content',
        (data) => ContentItem.fromJson((data as Map).cast<String, dynamic>()),
        body: payload);
  }

  Future<ContentItem> update(String id, Map<String, dynamic> payload) {
    return _transport.json<ContentItem>('PATCH', '/content/$id',
        (data) => ContentItem.fromJson((data as Map).cast<String, dynamic>()),
        body: payload);
  }

  Future<void> delete(String id) async {
    await _transport.request('DELETE', '/content/$id');
  }

  Future<List<ContentItem>> listByCreator(String creatorId, {int limit = 25, String? cursor}) {
    return _transport.jsonList<ContentItem>(
      'GET',
      '/content',
      (data) => ContentItem.fromJson((data as Map).cast<String, dynamic>()),
      query: {
        'creator_id': creatorId,
        'limit': '$limit',
        if (cursor != null) 'cursor': cursor,
      },
    );
  }

  Future<List<ContentItem>> listByCategory(String category, {int limit = 25, String? cursor}) {
    return _transport.jsonList<ContentItem>(
      'GET',
      '/content',
      (data) => ContentItem.fromJson((data as Map).cast<String, dynamic>()),
      query: {
        'category': category,
        'limit': '$limit',
        if (cursor != null) 'cursor': cursor,
      },
    );
  }

  Future<void> like(String id) async {
    await _transport.request('POST', '/content/$id/like');
  }

  Future<void> unlike(String id) async {
    await _transport.request('DELETE', '/content/$id/like');
  }

  Future<void> recordView(String id, {int durationSec = 0}) async {
    await _transport.request('POST', '/content/$id/view', body: {
      'duration_sec': durationSec,
    });
  }
}
