import 'models.dart';

class AnalyticsClient {
  AnalyticsClient(this._transport);
  final dynamic _transport;

  Future<void> trackEvent({required String eventName, Map<String, dynamic>? properties, String? sessionId}) async {
    await _transport.request('POST', '/events/track', body: {
      'event_name': eventName,
      if (sessionId != null) 'session_id': sessionId,
      if (properties != null) 'properties': properties,
    });
  }

  Future<void> trackBatch(List<Map<String, dynamic>> events) async {
    await _transport.request('POST', '/events/batch', body: events);
  }

  Future<RealtimeMetrics> realtime() {
    return _transport.json<RealtimeMetrics>('GET', '/metrics/realtime',
        (data) => RealtimeMetrics.fromJson((data as Map).cast<String, dynamic>()));
  }

  Future<List<dynamic>> historical(String metric, {DateTime? from, DateTime? to}) async {
    final data = await _transport.request('GET', '/metrics/historical', query: {
      'metric': metric,
      if (from != null) 'from': from.toIso8601String(),
      if (to != null) 'to': to.toIso8601String(),
    });
    if (data is List) return data;
    return const [];
  }

  Future<Map<String, dynamic>> getCreatorDashboard(String userId) async {
    final data = await _transport.request('GET', '/dashboards/creator/$userId');
    return (data as Map).cast<String, dynamic>();
  }

  Future<Map<String, dynamic>> getViewerDashboard(String userId) async {
    final data = await _transport.request('GET', '/dashboards/viewer/$userId');
    return (data as Map).cast<String, dynamic>();
  }

  Future<Map<String, dynamic>> getAdminDashboard() async {
    final data = await _transport.request('GET', '/dashboards/admin');
    return (data as Map).cast<String, dynamic>();
  }
}
