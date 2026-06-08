/// HTTP 客户端：基于 Dio，对齐 Web `apiClient` 的 2xx 解包与 ErrorVo 错误处理。
library;

import 'dart:ui';

import 'package:dio/dio.dart';
import 'package:gin_mam_mobile/core/api/api_error.dart';
import 'package:gin_mam_mobile/core/api/auth_storage.dart';
import 'package:gin_mam_mobile/core/constants/api_constants.dart';

/// 将 [DioException] 转为 [ApiError] 或带状态码的兜底错误（供测试与拦截器复用）。
ApiError mapDioExceptionToApiError(DioException error) {
  final status = error.response?.statusCode;
  final body = error.response?.data;

  if (body is Map<String, dynamic>) {
    final code = body['err_code'];
    if (code is int) {
      final msg = body['err_msg'];
      return ApiError(code, msg is String ? msg : (msg?.toString() ?? ''));
    }
  }

  final fallbackMsg = error.message ?? '请求失败';
  final fallbackCode = status ?? 0;
  return ApiError(fallbackCode, fallbackMsg);
}

/// Gin MAM 移动端 API 客户端。
class ApiClient {
  final Dio _dio;

  /// 令牌存储，供请求拦截器注入 Bearer 头。
  final AuthStorage authStorage;

  /// 收到 HTTP 401 且已清除本地令牌后的回调（通常跳转登录页）。
  final VoidCallback? onUnauthorized;

  /// 使用已有 [dio] 与 [authStorage] 构造（测试可注入自定义 Dio）。
  ApiClient({
    required Dio dio,
    required this.authStorage,
    this.onUnauthorized,
  }) : _dio = dio {
    _attachAuthInterceptors();
  }

  void _attachAuthInterceptors() {
    _dio.interceptors.add(
      InterceptorsWrapper(
        onRequest: (options, handler) async {
          final token = await authStorage.getToken();
          if (token != null && token.isNotEmpty) {
            options.headers['Authorization'] = 'Bearer $token';
          }
          handler.next(options);
        },
        onError: (error, handler) async {
          if (error.response?.statusCode == 401) {
            await authStorage.clearToken();
            onUnauthorized?.call();
          }
          handler.next(error);
        },
      ),
    );
  }

  /// 创建默认配置的客户端实例。
  factory ApiClient.create({
    required AuthStorage authStorage,
    VoidCallback? onUnauthorized,
    String baseUrl = defaultApiBaseUrl,
    Duration timeout = apiTimeout,
  }) {
    final dio = Dio(
      BaseOptions(
        baseUrl: baseUrl,
        connectTimeout: timeout,
        receiveTimeout: timeout,
        sendTimeout: timeout,
        headers: {'Accept': 'application/json'},
      ),
    );

    return ApiClient(
      dio: dio,
      authStorage: authStorage,
      onUnauthorized: onUnauthorized,
    );
  }

  /// 当前 Dio 实例（测试时替换 [httpClientAdapter]）。
  Dio get dio => _dio;

  /// GET 请求，成功时返回响应 JSON（Map / List / 标量），无额外 `code/data` 包装。
  Future<T> get<T>(
    String path, {
    Map<String, dynamic>? queryParameters,
  }) async {
    return _unwrap<T>(
      _dio.get<dynamic>(path, queryParameters: queryParameters),
    );
  }

  /// POST 请求，成功时返回响应 JSON 业务数据。
  Future<T> post<T>(
    String path, {
    Object? data,
    Map<String, dynamic>? queryParameters,
  }) async {
    return _unwrap<T>(
      _dio.post<dynamic>(
        path,
        data: data,
        queryParameters: queryParameters,
      ),
    );
  }

  /// PUT 请求，成功时返回响应 JSON 业务数据。
  Future<T> put<T>(
    String path, {
    Object? data,
    Map<String, dynamic>? queryParameters,
  }) async {
    return _unwrap<T>(
      _dio.put<dynamic>(
        path,
        data: data,
        queryParameters: queryParameters,
      ),
    );
  }

  /// DELETE 请求，成功时返回响应 JSON 业务数据。
  Future<T> delete<T>(
    String path, {
    Map<String, dynamic>? queryParameters,
  }) async {
    return _unwrap<T>(
      _dio.delete<dynamic>(path, queryParameters: queryParameters),
    );
  }

  Future<T> _unwrap<T>(Future<Response<dynamic>> call) async {
    try {
      final response = await call;
      return response.data as T;
    } on DioException catch (e) {
      throw mapDioExceptionToApiError(e);
    }
  }
}