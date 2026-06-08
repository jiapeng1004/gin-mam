/// 审核详情页路由参数：实例 ID 与媒资 ID。
library;

/// 传递给 [AuditDetailPage] 的路由参数。
class AuditDetailArgs {
  /// 流程实例 ID。
  final String instanceId;

  /// 媒资 ID。
  final String assetId;

  /// 构造审核详情路由参数。
  const AuditDetailArgs({
    required this.instanceId,
    required this.assetId,
  });
}
