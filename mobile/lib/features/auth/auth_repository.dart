/// 认证仓储：封装登录 API 并在成功后持久化 JWT。
library;

import 'package:gin_mam_mobile/core/api/api_client.dart';
import 'package:gin_mam_mobile/core/api/auth_storage.dart';
import 'package:gin_mam_mobile/features/auth/models/login_response.dart';

/// 认证相关数据访问层。
class AuthRepository {
  final ApiClient _apiClient;
  final AuthStorage _authStorage;

  /// 使用 [apiClient] 与 [authStorage] 构造。
  AuthRepository({
    required ApiClient apiClient,
    required AuthStorage authStorage,
  })  : _apiClient = apiClient,
        _authStorage = authStorage;

  /// 用户名密码登录，成功后写入本地令牌。
  Future<LoginResponse> login({
    required String username,
    required String password,
  }) async {
    final response = await _apiClient.post<Map<String, dynamic>>(
      '/auth/login',
      data: LoginRequest(username: username, password: password).toJson(),
    );
    final loginResponse = LoginResponse.fromJson(response);
    await _authStorage.setToken(loginResponse.token);
    return loginResponse;
  }

  /// 清除本地令牌（登出）。
  Future<void> logout() async {
    await _authStorage.clearToken();
  }

  /// 是否已登录。
  Future<bool> isAuthenticated() => _authStorage.isAuthenticated();
}
