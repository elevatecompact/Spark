import 'models.dart';

class WalletClient {
  WalletClient(this._transport);
  final dynamic _transport;

  Future<Wallet> me() {
    return _transport.json<Wallet>('GET', '/wallets/me',
        (data) => Wallet.fromJson((data as Map).cast<String, dynamic>()));
  }

  Future<Wallet> getByUser(String userId) {
    return _transport.json<Wallet>('GET', '/wallets/by-user',
        (data) => Wallet.fromJson((data as Map).cast<String, dynamic>()),
        query: {'user_id': userId});
  }

  Future<List<dynamic>> myBalances() async {
    final data = await _transport.request('GET', '/wallets/me/balances');
    if (data is List) return data;
    return const [];
  }

  Future<Map<String, dynamic>> deposit({required int amountCents, String currency = 'USD'}) async {
    final data = await _transport.request('POST', '/transactions/deposit', body: {
      'amount_cents': amountCents,
      'currency': currency,
    });
    return (data as Map).cast<String, dynamic>();
  }

  Future<Map<String, dynamic>> withdraw({required int amountCents, String currency = 'USD'}) async {
    final data = await _transport.request('POST', '/transactions/withdraw', body: {
      'amount_cents': amountCents,
      'currency': currency,
    });
    return (data as Map).cast<String, dynamic>();
  }

  Future<Map<String, dynamic>> transfer({required String toWalletId, required int amountCents}) async {
    final data = await _transport.request('POST', '/transactions/transfer', body: {
      'to_wallet_id': toWalletId,
      'amount_cents': amountCents,
    });
    return (data as Map).cast<String, dynamic>();
  }

  Future<Map<String, dynamic>> tip({required String toWalletId, required int amountCents}) async {
    final data = await _transport.request('POST', '/transactions/tip', body: {
      'to_wallet_id': toWalletId,
      'amount_cents': amountCents,
    });
    return (data as Map).cast<String, dynamic>();
  }

  Future<List<dynamic>> listTransactions({int limit = 50, String? cursor}) async {
    final data = await _transport.request('GET', '/transactions', query: {
      'limit': '$limit',
      if (cursor != null) 'cursor': cursor,
    });
    if (data is List) return data;
    return const [];
  }

  Future<Map<String, dynamic>> getTransaction(String id) async {
    final data = await _transport.request('GET', '/transactions/$id');
    return (data as Map).cast<String, dynamic>();
  }

  Future<Map<String, dynamic>> requestPayout({required int amountCents, required String method, String currency = 'USD'}) async {
    final data = await _transport.request('POST', '/payouts/request', body: {
      'amount_cents': amountCents,
      'method': method,
      'currency': currency,
    });
    return (data as Map).cast<String, dynamic>();
  }

  Future<List<dynamic>> listPayouts({int limit = 50, String? cursor}) async {
    final data = await _transport.request('GET', '/payouts', query: {
      'limit': '$limit',
      if (cursor != null) 'cursor': cursor,
    });
    if (data is List) return data;
    return const [];
  }

  Future<void> freeze() async {
    await _transport.request('POST', '/wallets/freeze');
  }

  Future<void> close() async {
    await _transport.request('POST', '/wallets/close');
  }
}
