import 'models.dart';

class DiscoveryClient {
  DiscoveryClient(this._transport);
  final dynamic _transport;

  Future<List<ContentItem>> featured({String? category, int limit = 25}) {
    return _transport.jsonList<ContentItem>(
      'GET',
      '/discovery/featured',
      (data) => ContentItem.fromJson((data as Map).cast<String, dynamic>()),
      query: {
        if (category != null) 'category': category,
        'limit': '$limit',
      },
    );
  }

  Future<List<ContentItem>> newReleases({String? category, int limit = 25}) {
    return _transport.jsonList<ContentItem>(
      'GET',
      '/discovery/new',
      (data) => ContentItem.fromJson((data as Map).cast<String, dynamic>()),
      query: {
        if (category != null) 'category': category,
        'limit': '$limit',
      },
    );
  }

  Future<List<ContentItem>> byCategory({required String category, int limit = 25}) {
    return _transport.jsonList<ContentItem>(
      'GET',
      '/discovery/categories/$category',
      (data) => ContentItem.fromJson((data as Map).cast<String, dynamic>()),
      query: {'limit': '$limit'},
    );
  }

  Future<List<String>> categories() {
    return _transport.jsonList<String>(
      'GET',
      '/discovery/categories',
      (data) => data as String,
    );
  }
}
