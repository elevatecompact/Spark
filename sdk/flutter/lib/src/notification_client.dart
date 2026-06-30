import 'models.dart';

class NotificationClient {
  NotificationClient(this._transport);
  final dynamic _transport;

  Future<List<Notification>> list({int limit = 50, String? cursor}) {
    return _transport.jsonList<Notification>(
      'GET',
      '/notifications',
      (data) => Notification.fromJson((data as Map).cast<String, dynamic>()),
      query: {
        'limit': '$limit',
        if (cursor != null) 'cursor': cursor,
      },
    );
  }

  Future<void> markRead(String id) async {
    await _transport.request('PATCH', '/notifications/$id/read');
  }

  Future<void> markAllRead() async {
    await _transport.request('POST', '/notifications/read-all');
  }

  Future<void> delete(String id) async {
    await _transport.request('DELETE', '/notifications/$id');
  }

  Future<void> registerDevice({required String platform, required String token}) async {
    await _transport.request('POST', '/notifications/devices', body: {
      'platform': platform,
      'token': token,
    });
  }

  Future<void> unregisterDevice(String id) async {
    await _transport.request('DELETE', '/notifications/devices/$id');
  }

  Future<Map<String, dynamic>> getPreferences() async {
    final data = await _transport.request('GET', '/notifications/preferences');
    return (data as Map).cast<String, dynamic>();
  }

  Future<void> updatePreferences(Map<String, dynamic> prefs) async {
    await _transport.request('PATCH', '/notifications/preferences', body: prefs);
  }

  Future<Map<String, dynamic>> deliveryStats() async {
    final data = await _transport.request('GET', '/notifications/admin/delivery-stats');
    return (data as Map).cast<String, dynamic>();
  }
}
