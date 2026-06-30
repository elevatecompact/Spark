import 'models.dart';

class PaymentClient {
  PaymentClient(this._transport);
  final dynamic _transport;

  Future<PaymentIntent> createIntent({required int amountCents, String currency = 'USD', String? idempotencyKey}) {
    return _transport.json<PaymentIntent>(
      'POST',
      '/payment-intents',
      (data) => PaymentIntent.fromJson((data as Map).cast<String, dynamic>()),
      body: {
        'amount_cents': amountCents,
        'currency': currency,
        if (idempotencyKey != null) 'idempotency_key': idempotencyKey,
      },
    );
  }

  Future<PaymentIntent> getIntent(String id) {
    return _transport.json<PaymentIntent>('GET', '/payment-intents/$id',
        (data) => PaymentIntent.fromJson((data as Map).cast<String, dynamic>()));
  }

  Future<List<PaymentIntent>> listIntents({int limit = 50, String? cursor}) {
    return _transport.jsonList<PaymentIntent>(
      'GET',
      '/payment-intents',
      (data) => PaymentIntent.fromJson((data as Map).cast<String, dynamic>()),
      query: {
        'limit': '$limit',
        if (cursor != null) 'cursor': cursor,
      },
    );
  }

  Future<PaymentIntent> confirmIntent({required String intentId, required String paymentMethodId}) {
    return _transport.json<PaymentIntent>(
      'POST',
      '/payment-intents/$intentId/confirm',
      (data) => PaymentIntent.fromJson((data as Map).cast<String, dynamic>()),
      body: {'payment_method_id': paymentMethodId},
    );
  }

  Future<void> cancelIntent(String intentId) async {
    await _transport.request('POST', '/payment-intents/$intentId/cancel');
  }

  Future<void> refundIntent({required String intentId, int? amountCents}) async {
    await _transport.request('POST', '/payment-intents/$intentId/refund', body: {
      if (amountCents != null) 'amount_cents': amountCents,
    });
  }

  Future<Map<String, dynamic>> createPaymentMethod({required String processor, required String type, required String token, bool setAsDefault = false}) async {
    final data = await _transport.request('POST', '/payment-methods', body: {
      'processor': processor,
      'type': type,
      'token': token,
      'set_as_default': setAsDefault,
    });
    return (data as Map).cast<String, dynamic>();
  }

  Future<List<dynamic>> listPaymentMethods() async {
    final data = await _transport.request('GET', '/payment-methods');
    if (data is List) return data;
    return const [];
  }

  Future<void> setDefaultPaymentMethod(String id) async {
    await _transport.request('PATCH', '/payment-methods/$id');
  }

  Future<void> deletePaymentMethod(String id) async {
    await _transport.request('DELETE', '/payment-methods/$id');
  }

  Future<Map<String, dynamic>> processorStatus() async {
    final data = await _transport.request('GET', '/admin/processors/status');
    return (data as Map).cast<String, dynamic>();
  }
}
