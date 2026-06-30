import 'dart:convert';

/// Thrown when a Spark REST call returns a non-2xx response.
class SparkApiException implements Exception {
  SparkApiException(this.statusCode, this.body);

  final int statusCode;
  final String body;

  String? get message {
    try {
      final dynamic parsed = jsonDecode(body);
      if (parsed is Map && parsed['message'] is String) {
        return parsed['message'] as String;
      }
      if (parsed is Map && parsed['error'] is String) {
        return parsed['error'] as String;
      }
    } catch (_) {
      // Body wasn't JSON; fall through.
    }
    return body.isEmpty ? 'HTTP $statusCode' : body;
  }

  @override
  String toString() => 'SparkApiException($statusCode): $message';
}
