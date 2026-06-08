// WorkflowRepository 单元测试：待我审分页与审核提交。

import 'dart:convert';

import 'package:dio/dio.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:gin_mam_mobile/core/api/api_client.dart';
import 'package:gin_mam_mobile/core/api/auth_storage.dart';
import 'package:gin_mam_mobile/features/audit/models/workflow_instance.dart';
import 'package:gin_mam_mobile/features/audit/workflow_repository.dart';
import 'package:shared_preferences/shared_preferences.dart';

import '../../helpers/stub_http_adapter.dart';

Map<String, dynamic> _readJsonBody(Object? data) {
  if (data is Map<String, dynamic>) return data;
  if (data is String) return jsonDecode(data) as Map<String, dynamic>;
  throw StateError('unexpected request body: $data');
}

void main() {
  group('WorkflowRepository', () {
    late AuthStorage authStorage;

    setUp(() async {
      SharedPreferences.setMockInitialValues({});
      final prefs = await SharedPreferences.getInstance();
      authStorage = AuthStorage(prefs);
    });

    ApiClient clientFor(Dio dio) => ApiClient(dio: dio, authStorage: authStorage);

    test('assignedToMe parses page result', () async {
      final dio = Dio(BaseOptions(baseUrl: 'http://test'));
      dio.httpClientAdapter = StubHttpAdapter((options) async {
        expect(options.path, '/asset-workflow/assigned-to-me');
        final body = _readJsonBody(options.data);
        expect(body['page'], 1);
        expect(body['pageSize'], 20);
        return jsonBody(
          {
            'list': [
              {
                'id': 'inst-1',
                'workflowDefId': 'def-1',
                'assetId': 'asset-1',
                'level': 1,
                'auditLevel': 2,
                'auditStatus': 0,
                'createdBy': 'user-1',
                'createdAt': '2026-01-01T00:00:00Z',
                'updatedAt': '2026-01-01T00:00:00Z',
              },
            ],
            'total': 1,
            'page': 1,
            'pageSize': 20,
          },
          200,
        );
      });

      final repo = WorkflowRepository(clientFor(dio));
      final page = await repo.assignedToMe(page: 1, pageSize: 20);

      expect(page.total, 1);
      expect(page.list, hasLength(1));
      expect(page.list.first.assetId, 'asset-1');
      expect(page.list.first.level, 1);
    });

    test('audit posts pass and remark', () async {
      Map<String, dynamic>? captured;
      final dio = Dio(BaseOptions(baseUrl: 'http://test'));
      dio.httpClientAdapter = StubHttpAdapter((options) async {
        expect(options.path, '/asset-workflow/audit');
        captured = _readJsonBody(options.data);
        return jsonBody(null, 200);
      });

      final repo = WorkflowRepository(clientFor(dio));
      await repo.audit(
        const AuditWorkflowRequest(
          instanceId: 'inst-1',
          pass: false,
          remark: '内容不合规',
        ),
      );

      expect(captured!['instanceId'], 'inst-1');
      expect(captured!['pass'], false);
      expect(captured!['remark'], '内容不合规');
    });
  });
}
