// AuthRepository 单元测试：登录成功写入 token、ErrorVo 透传。

import 'dart:convert';

import 'package:dio/dio.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:gin_mam_mobile/core/api/api_client.dart';
import 'package:gin_mam_mobile/core/api/api_error.dart';
import 'package:gin_mam_mobile/core/api/auth_storage.dart';
import 'package:gin_mam_mobile/core/constants/api_constants.dart';
import 'package:gin_mam_mobile/features/auth/auth_repository.dart';
import 'package:shared_preferences/shared_preferences.dart';

import '../../helpers/stub_http_adapter.dart';

Map<String, dynamic> _readJsonBody(Object? data) {
  if (data is Map<String, dynamic>) return data;
  if (data is String) return jsonDecode(data) as Map<String, dynamic>;
  throw StateError('unexpected request body: $data');
}

void main() {
  group('AuthRepository', () {
    late SharedPreferences prefs;
    late AuthStorage authStorage;

    setUp(() async {
      SharedPreferences.setMockInitialValues({});
      prefs = await SharedPreferences.getInstance();
      authStorage = AuthStorage(prefs);
    });

    test('login saves token on success', () async {
      final dio = Dio(BaseOptions(baseUrl: 'http://test'));
      dio.httpClientAdapter = StubHttpAdapter((options) async {
        expect(options.path, '/auth/login');
        expect(options.method, 'POST');
        final body = _readJsonBody(options.data);
        expect(body['username'], 'admin');
        expect(body['password'], 'admin123');
        return jsonBody(
          {'token': 'jwt-abc', 'expiresAt': '2026-12-31T00:00:00Z'},
          200,
        );
      });

      final repo = AuthRepository(
        apiClient: ApiClient(dio: dio, authStorage: authStorage),
        authStorage: authStorage,
      );

      final result = await repo.login(username: 'admin', password: 'admin123');

      expect(result.token, 'jwt-abc');
      expect(prefs.getString(authTokenKey), 'jwt-abc');
    });

    test('login throws ApiError on invalid credentials', () async {
      final dio = Dio(BaseOptions(baseUrl: 'http://test'));
      dio.httpClientAdapter = StubHttpAdapter((_) async {
        return jsonBody({'err_code': 40100, 'err_msg': '用户名或密码错误'}, 401);
      });

      final repo = AuthRepository(
        apiClient: ApiClient(dio: dio, authStorage: authStorage),
        authStorage: authStorage,
      );

      expect(
        () => repo.login(username: 'bad', password: 'bad'),
        throwsA(
          isA<ApiError>()
              .having((e) => e.errCode, 'errCode', 40100)
              .having((e) => e.errMsg, 'errMsg', '用户名或密码错误'),
        ),
      );
      expect(prefs.getString(authTokenKey), isNull);
    });
  });
}
