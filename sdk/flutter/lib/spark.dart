import 'dart:convert';
import 'package:http/http.dart' as http;

class Spark {
  final String baseUrl;
  String? _accessToken;
  final http.Client _client;

  Spark({
    this.baseUrl = 'https://api.spark.dev/api/v1',
    String? accessToken,
    http.Client? client,
  })  : _accessToken = accessToken,
        _client = client ?? http.Client();

  set accessToken(String token) => _accessToken = token;

  Map<String, String> get _headers => {
        'Content-Type': 'application/json',
        if (_accessToken != null) 'Authorization': 'Bearer $_accessToken',
      };

  Future<Map<String, dynamic>> _request(
    String method,
    String path, {
    Map<String, dynamic>? body,
  }) async {
    final uri = Uri.parse('$baseUrl$path');
    final request = http.Request(method, uri)..headers.addAll(_headers);
    if (body != null) {
      request.body = jsonEncode(body);
    }

    final streamed = await _client.send(request);
    final response = await http.Response.fromStream(streamed);

    if (response.statusCode < 200 || response.statusCode >= 300) {
      throw Exception('HTTP ${response.statusCode}: ${response.body}');
    }

    return jsonDecode(response.body) as Map<String, dynamic>;
  }

  Future<AuthResponse> register({
    required String email,
    required String username,
    required String password,
  }) async {
    final data = await _request('POST', '/auth/register', body: {
      'email': email,
      'username': username,
      'password': password,
    });
    final result = AuthResponse.fromJson(data);
    if (result.accessToken.isNotEmpty) {
      _accessToken = result.accessToken;
    }
    return result;
  }

  Future<AuthResponse> login({
    required String email,
    required String password,
  }) async {
    final data = await _request('POST', '/auth/login', body: {
      'email': email,
      'password': password,
    });
    final result = AuthResponse.fromJson(data);
    if (result.accessToken.isNotEmpty) {
      _accessToken = result.accessToken;
    }
    return result;
  }

  Future<User> me() async {
    final data = await _request('GET', '/users/me');
    return User.fromJson(data);
  }

  void dispose() {
    _client.close();
  }
}

class AuthResponse {
  final User user;
  final String accessToken;
  final String refreshToken;
  final int expiresIn;

  AuthResponse({
    required this.user,
    required this.accessToken,
    required this.refreshToken,
    required this.expiresIn,
  });

  factory AuthResponse.fromJson(Map<String, dynamic> json) => AuthResponse(
        user: User.fromJson(json['user']),
        accessToken: json['access_token'] as String,
        refreshToken: json['refresh_token'] as String,
        expiresIn: json['expires_in'] as int,
      );
}

class User {
  final String id;
  final String email;
  final String username;
  final String? displayName;
  final String? avatarUrl;
  final bool isCreator;
  final bool isVerified;
  final String createdAt;

  User({
    required this.id,
    required this.email,
    required this.username,
    this.displayName,
    this.avatarUrl,
    required this.isCreator,
    required this.isVerified,
    required this.createdAt,
  });

  factory User.fromJson(Map<String, dynamic> json) => User(
        id: json['id'] as String,
        email: json['email'] as String,
        username: json['username'] as String,
        displayName: json['display_name'] as String?,
        avatarUrl: json['avatar_url'] as String?,
        isCreator: json['is_creator'] as bool? ?? false,
        isVerified: json['is_verified'] as bool? ?? false,
        createdAt: json['created_at'] as String,
      );
}
