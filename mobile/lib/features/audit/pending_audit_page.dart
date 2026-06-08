/// 待审核列表页：下拉刷新，展示待我审流程实例。
library;

import 'package:antd_flutter_mobile/index.dart';
import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:gin_mam_mobile/app.dart';
import 'package:gin_mam_mobile/core/api/api_error.dart';
import 'package:gin_mam_mobile/core/providers.dart';
import 'package:gin_mam_mobile/features/audit/audit_detail_args.dart';
import 'package:gin_mam_mobile/features/audit/models/workflow_instance.dart';

/// 路由 `/audit/pending` 对应的待审核列表页。
class PendingAuditPage extends ConsumerWidget {
  /// 创建待审核列表页。
  const PendingAuditPage({super.key});

  /// 打开审核详情并在返回后刷新列表。
  Future<void> _openDetail(
    BuildContext context,
    WidgetRef ref,
    WorkflowInstanceVO item,
  ) async {
    await Navigator.of(context).pushNamed(
      GinMamApp.auditDetailRoute,
      arguments: AuditDetailArgs(
        instanceId: item.id,
        assetId: item.assetId,
      ),
    );
    await ref.read(pendingAuditListProvider.notifier).refresh();
  }

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final listAsync = ref.watch(pendingAuditListProvider);

    return Scaffold(
      appBar: AppBar(
        title: const Text('待我审核'),
        actions: [
          IconButton(
            tooltip: '退出登录',
            icon: const Icon(Icons.logout),
            onPressed: () async {
              await ref.read(authRepositoryProvider).logout();
              if (!context.mounted) return;
              Navigator.of(context).pushNamedAndRemoveUntil(
                GinMamApp.loginRoute,
                (_) => false,
              );
            },
          ),
        ],
      ),
      body: listAsync.when(
        loading: () => const Center(child: CircularProgressIndicator()),
        error: (error, _) {
          final msg = error is ApiError ? error.errMsg : '加载失败';
          return Center(
            child: Column(
              mainAxisAlignment: MainAxisAlignment.center,
              children: [
                Text(msg),
                const SizedBox(height: 16),
                AntdButton(
                  onTap: () =>
                      ref.read(pendingAuditListProvider.notifier).refresh(),
                  child: const Text('重试'),
                ),
              ],
            ),
          );
        },
        data: (items) {
          if (items.isEmpty) {
            return AntdPullToRefresh(
              onRefresh: () =>
                  ref.read(pendingAuditListProvider.notifier).refresh(),
              child: ListView(
                physics: const AlwaysScrollableScrollPhysics(),
                children: const [
                  SizedBox(height: 120),
                  Center(child: Text('暂无待审核项')),
                ],
              ),
            );
          }

          return AntdPullToRefresh(
            onRefresh: () =>
                ref.read(pendingAuditListProvider.notifier).refresh(),
            child: ListView.builder(
              physics: const AlwaysScrollableScrollPhysics(),
              itemCount: items.length,
              itemBuilder: (context, index) {
                final item = items[index];
                return ListTile(
                  title: Text('媒资 ${item.assetId}'),
                  subtitle: Text(
                    '层级 ${item.level}/${item.auditLevel} · 状态 ${item.auditStatus}',
                  ),
                  trailing: const Icon(Icons.chevron_right),
                  onTap: () => _openDetail(context, ref, item),
                );
              },
            ),
          );
        },
      ),
    );
  }
}
