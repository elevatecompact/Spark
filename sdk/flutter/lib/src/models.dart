// Domain models shared by every Spark SDK client. The model classes map 1:1 to
// the Go domain types in the platform's services.

import 'dart:convert';

T? _as<T>(dynamic v) => v is T ? v : null;

DateTime? _parseDate(dynamic v) {
  if (v is String && v.isNotEmpty) {
    return DateTime.tryParse(v);
  }
  return null;
}

class User {
  final String id;
  final String email;
  final String username;
  final String? displayName;
  final String? avatarUrl;
  final bool isCreator;
  final bool isVerified;
  final DateTime? createdAt;

  User({
    required this.id,
    required this.email,
    required this.username,
    this.displayName,
    this.avatarUrl,
    required this.isCreator,
    required this.isVerified,
    this.createdAt,
  });

  factory User.fromJson(Map<String, dynamic> json) => User(
        id: json['id'] as String,
        email: json['email'] as String,
        username: json['username'] as String,
        displayName: json['display_name'] as String?,
        avatarUrl: json['avatar_url'] as String?,
        isCreator: json['is_creator'] as bool? ?? false,
        isVerified: json['is_verified'] as bool? ?? false,
        createdAt: _parseDate(json['created_at']),
      );
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
        user: User.fromJson((json['user'] as Map).cast<String, dynamic>()),
        accessToken: json['access_token'] as String,
        refreshToken: json['refresh_token'] as String,
        expiresIn: json['expires_in'] as int,
      );
}

class ContentItem {
  final String id;
  final String creatorId;
  final String title;
  final String? description;
  final String? thumbnailUrl;
  final String? streamUrl;
  final List<String> tags;
  final String category;
  final int durationSec;
  final int viewCount;
  final int likeCount;
  final DateTime? publishedAt;

  ContentItem({
    required this.id,
    required this.creatorId,
    required this.title,
    this.description,
    this.thumbnailUrl,
    this.streamUrl,
    required this.tags,
    required this.category,
    required this.durationSec,
    required this.viewCount,
    required this.likeCount,
    this.publishedAt,
  });

  factory ContentItem.fromJson(Map<String, dynamic> json) => ContentItem(
        id: json['id'] as String,
        creatorId: json['creator_id'] as String,
        title: json['title'] as String,
        description: json['description'] as String?,
        thumbnailUrl: json['thumbnail_url'] as String?,
        streamUrl: json['stream_url'] as String?,
        tags: ((json['tags'] as List?) ?? const []).map((e) => e.toString()).toList(),
        category: json['category'] as String? ?? '',
        durationSec: json['duration_sec'] as int? ?? 0,
        viewCount: json['view_count'] as int? ?? 0,
        likeCount: json['like_count'] as int? ?? 0,
        publishedAt: _parseDate(json['published_at']),
      );
}

class Recommendation {
  final String contentId;
  final double score;
  final String? reason;

  Recommendation({required this.contentId, required this.score, this.reason});

  factory Recommendation.fromJson(Map<String, dynamic> json) => Recommendation(
        contentId: json['content_id'] as String,
        score: (json['score'] as num?)?.toDouble() ?? 0.0,
        reason: json['reason'] as String?,
      );
}

class Feed {
  final String type;
  final String? userId;
  final List<Recommendation> items;
  final DateTime servedAt;

  Feed({required this.type, this.userId, required this.items, required this.servedAt});

  factory Feed.fromJson(Map<String, dynamic> json) => Feed(
        type: json['type'] as String,
        userId: json['user_id'] as String?,
        items: ((json['items'] as List?) ?? const [])
            .map((e) => Recommendation.fromJson((e as Map).cast<String, dynamic>()))
            .toList(),
        servedAt: _parseDate(json['served_at']) ?? DateTime.now().toUtc(),
      );
}

class ChatMessage {
  final String id;
  final String roomId;
  final String userId;
  final String body;
  final DateTime sentAt;

  ChatMessage({
    required this.id,
    required this.roomId,
    required this.userId,
    required this.body,
    required this.sentAt,
  });

  factory ChatMessage.fromJson(Map<String, dynamic> json) => ChatMessage(
        id: json['id'] as String,
        roomId: json['room_id'] as String,
        userId: json['user_id'] as String,
        body: json['body'] as String,
        sentAt: _parseDate(json['sent_at']) ?? DateTime.now().toUtc(),
      );
}

class Notification {
  final String id;
  final String type;
  final String title;
  final String body;
  final String channel;
  final DateTime? readAt;
  final DateTime createdAt;

  Notification({
    required this.id,
    required this.type,
    required this.title,
    required this.body,
    required this.channel,
    this.readAt,
    required this.createdAt,
  });

  factory Notification.fromJson(Map<String, dynamic> json) => Notification(
        id: json['id'] as String,
        type: json['type'] as String,
        title: json['title'] as String,
        body: json['body'] as String,
        channel: json['channel'] as String,
        readAt: _parseDate(json['read_at']),
        createdAt: _parseDate(json['created_at']) ?? DateTime.now().toUtc(),
      );
}

class Wallet {
  final String id;
  final String userId;
  final int balanceCents;
  final String currency;
  final String status;

  Wallet({
    required this.id,
    required this.userId,
    required this.balanceCents,
    required this.currency,
    required this.status,
  });

  factory Wallet.fromJson(Map<String, dynamic> json) => Wallet(
        id: json['id'] as String,
        userId: json['user_id'] as String,
        balanceCents: json['balance_cents'] as int,
        currency: json['currency'] as String,
        status: json['status'] as String,
      );
}

class PaymentIntent {
  final String id;
  final int amountCents;
  final String currency;
  final String status;
  final String? externalId;

  PaymentIntent({
    required this.id,
    required this.amountCents,
    required this.currency,
    required this.status,
    this.externalId,
  });

  factory PaymentIntent.fromJson(Map<String, dynamic> json) => PaymentIntent(
        id: json['id'] as String,
        amountCents: json['amount_cents'] as int,
        currency: json['currency'] as String,
        status: json['status'] as String,
        externalId: json['external_id'] as String?,
      );
}

class SubscriptionPlan {
  final String id;
  final String name;
  final int priceCents;
  final String currency;
  final String billingPeriod;
  final List<String> benefits;
  final bool isActive;

  SubscriptionPlan({
    required this.id,
    required this.name,
    required this.priceCents,
    required this.currency,
    required this.billingPeriod,
    required this.benefits,
    required this.isActive,
  });

  factory SubscriptionPlan.fromJson(Map<String, dynamic> json) => SubscriptionPlan(
        id: json['id'] as String,
        name: json['name'] as String,
        priceCents: json['price_cents'] as int,
        currency: json['currency'] as String,
        billingPeriod: json['billing_period'] as String,
        benefits: ((json['benefits'] as List?) ?? const []).map((e) => e.toString()).toList(),
        isActive: json['is_active'] as bool? ?? true,
      );
}

class Subscription {
  final String id;
  final String userId;
  final String planId;
  final String status;
  final DateTime currentPeriodStart;
  final DateTime currentPeriodEnd;

  Subscription({
    required this.id,
    required this.userId,
    required this.planId,
    required this.status,
    required this.currentPeriodStart,
    required this.currentPeriodEnd,
  });

  factory Subscription.fromJson(Map<String, dynamic> json) => Subscription(
        id: json['id'] as String,
        userId: json['user_id'] as String,
        planId: json['plan_id'] as String,
        status: json['status'] as String,
        currentPeriodStart: _parseDate(json['current_period_start']) ?? DateTime.now().toUtc(),
        currentPeriodEnd: _parseDate(json['current_period_end']) ?? DateTime.now().toUtc(),
      );
}

class Gift {
  final String id;
  final String fromUserId;
  final String toUserId;
  final int amountCents;
  final String message;
  final DateTime sentAt;

  Gift({
    required this.id,
    required this.fromUserId,
    required this.toUserId,
    required this.amountCents,
    required this.message,
    required this.sentAt,
  });

  factory Gift.fromJson(Map<String, dynamic> json) => Gift(
        id: json['id'] as String,
        fromUserId: json['from_user_id'] as String,
        toUserId: json['to_user_id'] as String,
        amountCents: json['amount_cents'] as int,
        message: json['message'] as String? ?? '',
        sentAt: _parseDate(json['sent_at']) ?? DateTime.now().toUtc(),
      );
}

class Stream {
  final String id;
  final String creatorId;
  final String title;
  final String status;
  final int viewerCount;
  final DateTime? startedAt;

  Stream({
    required this.id,
    required this.creatorId,
    required this.title,
    required this.status,
    required this.viewerCount,
    this.startedAt,
  });

  factory Stream.fromJson(Map<String, dynamic> json) => Stream(
        id: json['id'] as String,
        creatorId: json['creator_id'] as String,
        title: json['title'] as String,
        status: json['status'] as String,
        viewerCount: json['viewer_count'] as int? ?? 0,
        startedAt: _parseDate(json['started_at']),
      );
}

class SearchResult {
  final List<ContentItem> content;
  final List<User> creators;
  final int total;

  SearchResult({required this.content, required this.creators, required this.total});

  factory SearchResult.fromJson(Map<String, dynamic> json) {
    final contentList = ((json['content'] as List?) ?? const [])
        .map((e) => ContentItem.fromJson((e as Map).cast<String, dynamic>()))
        .toList();
    final creators = ((json['creators'] as List?) ?? const [])
        .map((e) => User.fromJson((e as Map).cast<String, dynamic>()))
        .toList();
    return SearchResult(
      content: contentList,
      creators: creators,
      total: json['total'] as int? ?? contentList.length + creators.length,
    );
  }
}

class RealtimeMetrics {
  final int activeViewers;
  final int giftsSent;
  final int chatMessages;
  final int newSubscribers;
  final DateTime timestamp;

  RealtimeMetrics({
    required this.activeViewers,
    required this.giftsSent,
    required this.chatMessages,
    required this.newSubscribers,
    required this.timestamp,
  });

  factory RealtimeMetrics.fromJson(Map<String, dynamic> json) => RealtimeMetrics(
        activeViewers: json['active_viewers'] as int? ?? 0,
        giftsSent: json['gifts_sent'] as int? ?? 0,
        chatMessages: json['chat_messages'] as int? ?? 0,
        newSubscribers: json['new_subscribers'] as int? ?? 0,
        timestamp: _parseDate(json['timestamp']) ?? DateTime.now().toUtc(),
      );
}

// Convenience re-export for callers that need to inspect raw payloads.
const jsonEncode = _sparkJsonEncode;
const jsonDecode = _sparkJsonDecode;

String _sparkJsonEncode(Object? value) => json.encode(value);
dynamic _sparkJsonDecode(String source) => json.decode(source);
