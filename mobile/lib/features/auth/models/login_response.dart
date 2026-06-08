/// 登录接口请求/响应模型，对齐 Web `authApi.ts`。
library;

/// 登录请求体。
class LoginRequest {
  /// 用户名。
  final String username;

  /// 密码。
  final String password;

  /// 构造登录请求。
  const LoginRequest({
    required this.username,
    required this.password,
  });

  /// 转为 JSON 请求体。
  Map<String, dynamic> toJson() => {
        'username': username,
        'password': password,
      };
}

/// 登录成功响应（camelCase）。
class LoginResponse {
  /// JWT 访问令牌。
  final String token;

  /// 过期时间（ISO 8601 字符串）。
  final String expiresAt;

  /// 构造登录响应。
  const LoginResponse({
    required this.token,
    required this.expiresAt,
  });

  /// 从 JSON 解析登录响应。
  factory LoginResponse.fromJson(Map<String, dynamic> json) {
    return LoginResponse(
      token: json['token'] as String,
      expiresAt: json['expiresAt'] as String,
    );
  }
}
