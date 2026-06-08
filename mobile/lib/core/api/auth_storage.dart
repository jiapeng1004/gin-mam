/// 访问令牌持久化：读写 SharedPreferences，供 Dio 拦截器注入 Authorization。
library;

import 'package:gin_mam_mobile/core/constants/api_constants.dart';
import 'package:shared_preferences/shared_preferences.dart';

/// 封装认证令牌的本地存储。
class AuthStorage {
  final SharedPreferences _prefs;

  /// 使用已初始化的 [prefs] 实例构造。
  AuthStorage(this._prefs);

  /// 读取当前保存的访问令牌；未登录时返回 `null`。
  Future<String?> getToken() async {
    return _prefs.getString(authTokenKey);
  }

  /// 保存访问令牌（登录接口返回的 JWT 或同类字符串）。
  Future<void> setToken(String token) async {
    await _prefs.setString(authTokenKey, token);
  }

  /// 清除本地令牌（登出或收到 401 时使用）。
  Future<void> clearToken() async {
    await _prefs.remove(authTokenKey);
  }

  /// 是否已登录（脚手架：仅检查本地是否存在 token）。
  Future<bool> isAuthenticated() async {
    final token = await getToken();
    return token != null && token.isNotEmpty;
  }
}