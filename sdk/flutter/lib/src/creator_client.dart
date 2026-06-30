import 'models.dart';

class CreatorClient {
  CreatorClient(this._transport);
  final dynamic _transport;

  Future<Map<String, dynamic>> apply({Map<String, dynamic>? profile}) async {
    final data = await _transport.request('POST', '/creators/apply', body: profile ?? const {});
    return (data as Map).cast<String, dynamic>();
  }

  Future<Map<String, dynamic>> myProfile() async {
    final data = await _transport.request('GET', '/creators/me');
    return (data as Map).cast<String, dynamic>();
  }

  Future<Map<String, dynamic>> updateProfile(Map<String, dynamic> payload) async {
    final data = await _transport.request('PATCH', '/creators/me', body: payload);
    return (data as Map).cast<String, dynamic>();
  }

  Future<Map<String, dynamic>> getEarnings({DateTime? from, DateTime? to}) async {
    final data = await _transport.request('GET', '/creators/me/earnings', query: {
      if (from != null) 'from': from.toIso8601String(),
      if (to != null) 'to': to.toIso8601String(),
    });
    return (data as Map).cast<String, dynamic>();
  }

  Future<Map<String, dynamic>> getAnalytics() async {
    final data = await _transport.request('GET', '/creators/me/analytics');
    return (data as Map).cast<String, dynamic>();
  }
}
