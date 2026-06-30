import 'models.dart';

class RecommendationClient {
  RecommendationClient(this._transport);
  final dynamic _transport;

  Future<Feed> home({required String userId, int limit = 50}) {
    return _transport.json<Feed>(
      'GET',
      '/feeds/home',
      (data) => Feed.fromJson((data as Map).cast<String, dynamic>()),
      query: {'userId': userId, 'limit': '$limit'},
    );
  }

  Future<Feed> trending({int limit = 50}) {
    return _transport.json<Feed>(
      'GET',
      '/feeds/trending',
      (data) => Feed.fromJson((data as Map).cast<String, dynamic>()),
      query: {'limit': '$limit'},
    );
  }

  Future<Feed> upNext({required String userId, required String contentId, int limit = 10}) {
    return _transport.json<Feed>(
      'GET',
      '/feeds/up-next/$contentId',
      (data) => Feed.fromJson((data as Map).cast<String, dynamic>()),
      query: {'userId': userId, 'limit': '$limit'},
    );
  }

  Future<Feed> similar({required String contentId, int limit = 10}) {
    return _transport.json<Feed>(
      'GET',
      '/feeds/similar/$contentId',
      (data) => Feed.fromJson((data as Map).cast<String, dynamic>()),
      query: {'limit': '$limit'},
    );
  }

  Future<Feed> creator({required String creatorId, int limit = 50}) {
    return _transport.json<Feed>(
      'GET',
      '/feeds/creator/$creatorId',
      (data) => Feed.fromJson((data as Map).cast<String, dynamic>()),
      query: {'limit': '$limit'},
    );
  }

  Future<void> recordClick({required String userId, required String contentId}) async {
    await _transport.request('POST', '/feedback/click', body: {
      'userId': userId,
      'contentId': contentId,
    });
  }

  Future<void> recordDismiss({required String userId, required String contentId}) async {
    await _transport.request('POST', '/feedback/dismiss', body: {
      'userId': userId,
      'contentId': contentId,
    });
  }
}
