// ApiClient 单元测试：2xx 解包、ErrorVo 4xx、Authorization 头。

import 'dart:convert';

import 'package:dio/dio.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:gin_mam_mobile/core/api/api_client.dart';
import 'package:gin_mam_mobile/core/api/api_error.dart';
import 'package:gin_mam_mobile/core/api/auth_storage.dart';
import 'package:gin_mam_mobile/core/constants/api_constants.dart';
import 'package:shared_preferences/shared_preferences.dart';

/// 测试用 HTTP 适配器：按路径返回固定 JSON 响应。
class StubHttpAdapter implements HttpClientAdapter {
  StubHttpAdapter(this.handler);

  final Future<ResponseBody> Function(RequestOptions options) handler;

  @override
  Future<ResponseBody> fetch(
    RequestOptions options,
    Stream<List<int>>? requestStream,
    Future<void>? cancelFuture,
  ) {
    return handler(options);
  }

  @override
  void close({bool force = false}) {}
}

ResponseBody jsonBody(Object data, int status) {
  return ResponseBody.fromString(
    jsonEncode(data),
    status,
    headers: {
      Headers.contentTypeHeader: [Headers.jsonContentType],
    },
  );
}

void main() {
  group('ApiClient', () {
    late SharedPreferences prefs;
    late AuthStorage authStorage;

    setUp(() async {
      SharedPreferences.setMockInitialValues({});
      prefs = await SharedPreferences.getInstance();
      authStorage = AuthStorage(prefs);
    });

    test('returns business payload on 200 without wrapper', () async {
      final dio = Dio(BaseOptions(baseUrl: 'http://test'));
      dio.httpClientAdapter = StubHttpAdapter((options) async {
        return jsonBody({'id': 42, 'name': 'asset'}, 200);
      });

      final client = ApiClient(dio: dio, authStorage: authStorage);
      final result = await client.get<Map<String, dynamic>>('/assets/42');

      expect(result, {'id': 42, 'name': 'asset'});
    });

    test('throws ApiError with errCode from ErrorVo on 404', () async {
      final dio = Dio(BaseOptions(baseUrl: 'http://test'));
      dio.httpClientAdapter = StubHttpAdapter((options) async {
        return jsonBody({'err_code': 40401, 'err_msg': '资源不存在'}, 404);
      });

      final client = ApiClient(dio: dio, authStorage: authStorage);

      expect(
        () => client.get('/missing'),
        throwsA(
          isA<ApiError>()
              .having((e) => e.errCode, 'errCode', 40401)
              .having((e) => e.errMsg, 'errMsg', '资源不存在'),
        ),
      );
    });

    test('sends Authorization header when token is stored', () async {
      await prefs.setString(authTokenKey, 'test-jwt-token');

      String? capturedAuth;
      final dio = Dio(BaseOptions(baseUrl: 'http://test'));
      dio.httpClientAdapter = StubHttpAdapter((options) async {
        capturedAuth = options.headers['Authorization'] as String?;
        return jsonBody({'ok': true}, 200);
      });

      final client = ApiClient(dio: dio, authStorage: authStorage);

      await client.get('/me');

      expect(capturedAuth, 'Bearer test-jwt-token');
    });
  });
}