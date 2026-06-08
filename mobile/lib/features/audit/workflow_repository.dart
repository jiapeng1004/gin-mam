/// 工作流仓储：待我审列表与审核提交。
library;

import 'package:gin_mam_mobile/core/api/api_client.dart';
import 'package:gin_mam_mobile/features/audit/models/workflow_instance.dart';

/// 工作流审核相关数据访问层。
class WorkflowRepository {
  final ApiClient _apiClient;

  /// 使用 [apiClient] 构造。
  WorkflowRepository(this._apiClient);

  /// 分页查询待当前用户审核的流程实例。
  Future<WorkflowInstancePageResult> assignedToMe({
    required int page,
    required int pageSize,
  }) async {
    final response = await _apiClient.post<Map<String, dynamic>>(
      '/asset-workflow/assigned-to-me',
      data: {'page': page, 'pageSize': pageSize},
    );
    return WorkflowInstancePageResult.fromJson(response);
  }

  /// 对指定实例执行通过或打回。
  Future<void> audit(AuditWorkflowRequest request) async {
    await _apiClient.post<void>(
      '/asset-workflow/audit',
      data: request.toJson(),
    );
  }
}
