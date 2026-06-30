// Spark Platform Flutter SDK
//
// This library provides a typed client for every public REST endpoint exposed
// by the Spark platform (identity, content, chat, recommendations,
// notifications, wallet, payment, subscription, streaming, gift, search,
// analytics, and admin).
//
// The `Spark` class is the entry point. Each domain has its own sub-client
// (e.g. `spark.identity`, `spark.content`) so applications can opt in to just
// the surfaces they need while still sharing a single HTTP client and
// authentication token.

import 'dart:async';
import 'dart:convert';

import 'package:http/http.dart' as http;

import 'src/analytics_client.dart';
import 'src/auth_client.dart';
import 'src/chat_client.dart';
import 'src/content_client.dart';
import 'src/creator_client.dart';
import 'src/discovery_client.dart';
import 'src/errors.dart';
import 'src/gift_client.dart';
import 'src/models.dart';
import 'src/notification_client.dart';
import 'src/payment_client.dart';
import 'src/recommendation_client.dart';
import 'src/search_client.dart';
import 'src/stream_client.dart';
import 'src/subscription_client.dart';
import 'src/wallet_client.dart';

export 'src/analytics_client.dart';
export 'src/auth_client.dart';
export 'src/chat_client.dart';
export 'src/content_client.dart';
export 'src/creator_client.dart';
export 'src/discovery_client.dart';
export 'src/errors.dart';
export 'src/gift_client.dart';
export 'src/models.dart';
export 'src/notification_client.dart';
export 'src/payment_client.dart';
export 'src/recommendation_client.dart';
export 'src/search_client.dart';
export 'src/stream_client.dart';
export 'src/subscription_client.dart';
export 'src/wallet_client.dart';

class Spark {
  Spark({
    this.baseUrl = 'https://api.spark.dev/api/v1',
    String? accessToken,
    http.Client? client,
  })  : _accessToken = accessToken,
        _client = client ?? http.Client() {
    _transport = _HttpTransport(_client, baseUrl, () => _accessToken);
    identity = AuthClient(_transport);
    content = ContentClient(_transport);
    chat = ChatClient(_transport);
    notifications = NotificationClient(_transport);
    recommendations = RecommendationClient(_transport);
    wallet = WalletClient(_transport);
    payment = PaymentClient(_transport);
    subscription = SubscriptionClient(_transport);
    stream = StreamClient(_transport);
    gifts = GiftClient(_transport);
    search = SearchClient(_transport);
    discovery = DiscoveryClient(_transport);
    analytics = AnalyticsClient(_transport);
    creator = CreatorClient(_transport);
  }

  final String baseUrl;
  String? _accessToken;
  final http.Client _client;
  late final _HttpTransport _transport;

  late final AuthClient identity;
  late final ContentClient content;
  late final ChatClient chat;
  late final NotificationClient notifications;
  late final RecommendationClient recommendations;
  late final WalletClient wallet;
  late final PaymentClient payment;
  late final SubscriptionClient subscription;
  late final StreamClient stream;
  late final GiftClient gifts;
  late final SearchClient search;
  late final DiscoveryClient discovery;
  late final AnalyticsClient analytics;
  late final CreatorClient creator;

  String? get accessToken => _accessToken;
  set accessToken(String? token) {
    _accessToken = token;
  }

  /// Convenience wrapper for ad-hoc REST calls.
  Future<Map<String, dynamic>> request(
    String method,
    String path, {
    Map<String, dynamic>? body,
    Map<String, String>? query,
  }) async {
    return _transport.request(method, path, body: body, query: query);
  }

  void dispose() {
    _client.close();
  }
}

/// Internal HTTP transport used by every sub-client.
class _HttpTransport {
  _HttpTransport(this._client, this.baseUrl, this._tokenAccessor);

  final http.Client _client;
  final String baseUrl;
  final String? Function() _tokenAccessor;

  Map<String, String> get _headers => {
        'Content-Type': 'application/json',
        if (_tokenAccessor() != null) 'Authorization': 'Bearer ${_tokenAccessor()!}',
      };

  Future<dynamic> request(
    String method,
    String path, {
    Map<String, dynamic>? body,
    Map<String, String>? query,
  }) async {
    final base = Uri.parse(baseUrl);
    final uri = base.replace(
      path: base.path.endsWith('/') ? '${base.path}${path.replaceFirst('/', '')}' : '${base.path}$path',
      queryParameters: query,
    );
    final request = http.Request(method, uri)..headers.addAll(_headers);
    if (body != null) {
      request.body = jsonEncode(body);
    }
    final streamed = await _client.send(request);
    final response = await http.Response.fromStream(streamed);

    if (response.statusCode == 204 || response.body.isEmpty) {
      if (response.statusCode >= 200 && response.statusCode < 300) {
        return null;
      }
    }

    if (response.statusCode < 200 || response.statusCode >= 300) {
      throw SparkApiException(response.statusCode, response.body);
    }

    if (response.headers['content-type']?.contains('application/json') ?? false) {
      return jsonDecode(response.body);
    }
    return response.body;
  }

  Future<T> json<T>(
    String method,
    String path,
    T Function(dynamic json) parse, {
    Map<String, dynamic>? body,
    Map<String, String>? query,
  }) async {
    final data = await request(method, path, body: body, query: query);
    return parse(data);
  }

  Future<List<T>> jsonList<T>(
    String method,
    String path,
    T Function(dynamic json) parse, {
    Map<String, dynamic>? body,
    Map<String, String>? query,
  }) async {
    final data = await request(method, path, body: body, query: query);
    if (data is! List) {
      return const <T>[];
    }
    return data.map(parse).toList();
  }
}
