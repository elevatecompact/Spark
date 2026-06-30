import 'models.dart';

class StreamClient {
  StreamClient(this._transport);
  final dynamic _transport;

  Future<Stream> get(String id) {
    return _transport.json<Stream>('GET', '/streams/$id',
        (data) => Stream.fromJson((data as Map).cast<String, dynamic>()));
  }

  Future<List<Stream>> live({int limit = 25}) {
    return _transport.jsonList<Stream>(
      'GET',
      '/streams/live',
      (data) => Stream.fromJson((data as Map).cast<String, dynamic>()),
      query: {'limit': '$limit'},
    );
  }

  Future<Stream> start({required String title, String? description, List<String> tags = const []}) {
    return _transport.json<Stream>(
      'POST',
      '/streams',
      (data) => Stream.fromJson((data as Map).cast<String, dynamic>()),
      body: {
        'title': title,
        if (description != null) 'description': description,
        'tags': tags,
      },
    );
  }

  Future<void> end(String id) async {
    await _transport.request('POST', '/streams/$id/end');
  }

  Future<String> ingestUrl(String id) async {
    final data = await _transport.request('POST', '/streams/$id/ingest-url');
    if (data is Map) {
      return (data['url'] as String?) ?? '';
    }
    return data?.toString() ?? '';
  }

  Future<String> playbackUrl(String id) async {
    final data = await _transport.request('GET', '/streams/$id/playback-url');
    if (data is Map) {
      return (data['url'] as String?) ?? '';
    }
    return data?.toString() ?? '';
  }
}
