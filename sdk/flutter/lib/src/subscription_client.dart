import 'models.dart';

class SubscriptionClient {
  SubscriptionClient(this._transport);
  final dynamic _transport;

  Future<List<SubscriptionPlan>> listPlans({String? creatorId, int limit = 50}) {
    return _transport.jsonList<SubscriptionPlan>(
      'GET',
      '/plans',
      (data) => SubscriptionPlan.fromJson((data as Map).cast<String, dynamic>()),
      query: {
        if (creatorId != null) 'creator_id': creatorId,
        'limit': '$limit',
      },
    );
  }

  Future<SubscriptionPlan> createPlan(Map<String, dynamic> payload) {
    return _transport.json<SubscriptionPlan>('POST', '/plans',
        (data) => SubscriptionPlan.fromJson((data as Map).cast<String, dynamic>()),
        body: payload);
  }

  Future<List<Subscription>> mySubscriptions() {
    return _transport.jsonList<Subscription>(
      'GET',
      '/subscriptions/me',
      (data) => Subscription.fromJson((data as Map).cast<String, dynamic>()),
    );
  }

  Future<Subscription> subscribe({required String planId}) {
    return _transport.json<Subscription>(
      'POST',
      '/subscriptions',
      (data) => Subscription.fromJson((data as Map).cast<String, dynamic>()),
      body: {'plan_id': planId},
    );
  }

  Future<Subscription> changePlan({required String subscriptionId, required String newPlanId}) {
    return _transport.json<Subscription>(
      'POST',
      '/subscriptions/$subscriptionId/change-plan',
      (data) => Subscription.fromJson((data as Map).cast<String, dynamic>()),
      body: {'new_plan_id': newPlanId},
    );
  }

  Future<void> cancel(String subscriptionId) async {
    await _transport.request('POST', '/subscriptions/$subscriptionId/cancel');
  }

  Future<void> reactivate(String subscriptionId) async {
    await _transport.request('POST', '/subscriptions/$subscriptionId/reactivate');
  }
}
