import 'models.dart';

class AuthClient {
  AuthClient(this._transport);

  final dynamic _transport;

  Future<AuthResponse> register({
    required String email,
    required String username,
    required String password,
    String? displayName,
  }) {
    return _transport.json<AuthResponse>('POST', '/auth/register', (data) {
      final result = AuthResponse.fromJson((data as Map).cast<String, dynamic>());
      return result;
    }, body: {
      'email': email,
      'username': username,
      'password': password,
      if (displayName != null) 'display_name': displayName,
    });
  }

  Future<AuthResponse> login({required String email, required String password}) {
    return _transport.json<AuthResponse>('POST', '/auth/login', (data) {
      return AuthResponse.fromJson((data as Map).cast<String, dynamic>());
    }, body: {'email': email, 'password': password});
  }

  Future<AuthResponse> refresh(String refreshToken) {
    return _transport.json<AuthResponse>('POST', '/auth/refresh', (data) {
      return AuthResponse.fromJson((data as Map).cast<String, dynamic>());
    }, body: {'refresh_token': refreshToken});
  }

  Future<void> logout() async {
    await _transport.request('POST', '/auth/logout');
  }

  Future<void> requestPasswordReset(String email) async {
    await _transport.request('POST', '/auth/password/reset/request', body: {'email': email});
  }

  Future<void> confirmPasswordReset({required String token, required String newPassword}) async {
    await _transport.request('POST', '/auth/password/reset/confirm', body: {
      'token': token,
      'new_password': newPassword,
    });
  }

  Future<void> verifyEmail(String token) async {
    await _transport.request('POST', '/auth/email/verify', body: {'token': token});
  }

  Future<User> me() {
    return _transport.json<User>('GET', '/users/me',
        (data) => User.fromJson((data as Map).cast<String, dynamic>()));
  }

  Future<User> updateProfile({
    String? displayName,
    String? bio,
    String? avatarUrl,
  }) {
    return _transport.json<User>(
      'PATCH',
      '/users/me',
      (data) => User.fromJson((data as Map).cast<String, dynamic>()),
      body: {
        if (displayName != null) 'display_name': displayName,
        if (bio != null) 'bio': bio,
        if (avatarUrl != null) 'avatar_url': avatarUrl,
      },
    );
  }

  Future<User> getUser(String userId) {
    return _transport.json<User>('GET', '/users/$userId',
        (data) => User.fromJson((data as Map).cast<String, dynamic>()));
  }

  Future<void> changePassword({required String currentPassword, required String newPassword}) async {
    await _transport.request('POST', '/users/me/password', body: {
      'current_password': currentPassword,
      'new_password': newPassword,
    });
  }
}
