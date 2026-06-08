/// 工作流实例模型，对齐 Web `workflowApi.ts` 中的 WorkflowInstanceVO。
library;

/// 流程实例视图对象。
class WorkflowInstanceVO {
  /// 实例 ID。
  final String id;

  /// 流程定义 ID。
  final String workflowDefId;

  /// 关联媒资 ID。
  final String assetId;

  /// 当前层级。
  final int level;

  /// 审核总层级数。
  final int auditLevel;

  /// 审核状态码。
  final int auditStatus;

  /// 创建人用户 ID。
  final String createdBy;

  /// 创建时间。
  final String createdAt;

  /// 更新时间。
  final String updatedAt;

  /// 构造流程实例 VO。
  const WorkflowInstanceVO({
    required this.id,
    required this.workflowDefId,
    required this.assetId,
    required this.level,
    required this.auditLevel,
    required this.auditStatus,
    required this.createdBy,
    required this.createdAt,
    required this.updatedAt,
  });

  /// 从 JSON 解析流程实例。
  factory WorkflowInstanceVO.fromJson(Map<String, dynamic> json) {
    return WorkflowInstanceVO(
      id: json['id'] as String,
      workflowDefId: json['workflowDefId'] as String,
      assetId: json['assetId'] as String,
      level: json['level'] as int,
      auditLevel: json['auditLevel'] as int,
      auditStatus: json['auditStatus'] as int,
      createdBy: json['createdBy'] as String,
      createdAt: json['createdAt'] as String,
      updatedAt: json['updatedAt'] as String,
    );
  }
}

/// 流程实例分页结果。
class WorkflowInstancePageResult {
  /// 当前页数据。
  final List<WorkflowInstanceVO> list;

  /// 总条数。
  final int total;

  /// 当前页码。
  final int page;

  /// 每页条数。
  final int pageSize;

  /// 构造分页结果。
  const WorkflowInstancePageResult({
    required this.list,
    required this.total,
    required this.page,
    required this.pageSize,
  });

  /// 从 JSON 解析分页结果。
  factory WorkflowInstancePageResult.fromJson(Map<String, dynamic> json) {
    final rawList = json['list'] as List<dynamic>;
    return WorkflowInstancePageResult(
      list: rawList
          .map((e) => WorkflowInstanceVO.fromJson(e as Map<String, dynamic>))
          .toList(),
      total: json['total'] as int,
      page: json['page'] as int,
      pageSize: json['pageSize'] as int,
    );
  }
}

/// 审核操作请求体。
class AuditWorkflowRequest {
  /// 流程实例 ID。
  final String instanceId;

  /// 是否通过（false 为打回）。
  final bool pass;

  /// 审核备注（可选）。
  final String? remark;

  /// 构造审核请求。
  const AuditWorkflowRequest({
    required this.instanceId,
    required this.pass,
    this.remark,
  });

  /// 转为 JSON 请求体。
  Map<String, dynamic> toJson() {
    return {
      'instanceId': instanceId,
      'pass': pass,
      if (remark != null && remark!.isNotEmpty) 'remark': remark,
    };
  }
}
