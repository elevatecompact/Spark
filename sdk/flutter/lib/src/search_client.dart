import 'models.dart';

class SearchClient {
  SearchClient(this._transport);
  final dynamic _transport;

  Future<SearchResult> search({required String query, String? type, int limit = 25}) {
    return _transport.json<SearchResult>(
      'GET',
      '/search',
      (data) => SearchResult.fromJson((data as Map).cast<String, dynamic>()),
      query: {
        'q': query,
        if (type != null) 'type': type,
        'limit': '$limit',
      },
    );
  }

  Future<List<String>> suggestions({required String query, int limit = 10}) {
    return _transport.jsonList<String>(
      'GET',
      '/search/suggestions',
      (data) => data as String,
      query: {'q': query, 'limit': '$limit'},
    );
  }

  Future<List<ContentItem>> byTag({required String tag, int limit = 25}) {
    return _transport.jsonList<ContentItem>(
      'GET',
      '/search/tags/$tag',
      (data) => ContentItem.fromJson((data as Map).cast<String, dynamic>()),
      query: {'limit': '$limit'},
    );
  }
}
