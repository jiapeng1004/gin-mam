/// Riverpod 全局 Provider：API 客户端、仓储与业务状态。
library;

import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:gin_mam_mobile/app.dart';
import 'package:gin_mam_mobile/core/api/api_client.dart';
import 'package:gin_mam_mobile/core/api/auth_storage.dart';
import 'package:gin_mam_mobile/features/asset/asset_repository.dart';
import 'package:gin_mam_mobile/features/audit/models/workflow_instance.dart';
import 'package:gin_mam_mobile/features/audit/workflow_repository.dart';
import 'package:gin_mam_mobile/features/auth/auth_repository.dart';
import 'package:shared_preferences/shared_preferences.dart';

/// 全局 NavigatorKey，供 401 回调跳转登录页。
final navigatorKeyProvider = Provider<GlobalKey<NavigatorState>>(
  (ref) => GlobalKey<NavigatorState>(),
);

/// 已初始化的 SharedPreferences（在 [main] 中注入）。
final sharedPreferencesProvider = Provider<SharedPreferences>((ref) {
  throw UnimplementedError('sharedPreferencesProvider 需在 main 中 override');
});

/// 认证令牌本地存储。
final authStorageProvider = Provider<AuthStorage>((ref) {
  return AuthStorage(ref.watch(sharedPreferencesProvider));
});

/// HTTP API 客户端（含 Bearer 注入与 401 处理）。
final apiClientProvider = Provider<ApiClient>((ref) {
  final authStorage = ref.watch(authStorageProvider);
  final navigatorKey = ref.watch(navigatorKeyProvider);
  return ApiClient.create(
    authStorage: authStorage,
    onUnauthorized: () {
      navigatorKey.currentState?.pushNamedAndRemoveUntil(
        GinMamApp.loginRoute,
        (_) => false,
      );
    },
  );
});

/// 认证仓储。
final authRepositoryProvider = Provider<AuthRepository>((ref) {
  return AuthRepository(
    apiClient: ref.watch(apiClientProvider),
    authStorage: ref.watch(authStorageProvider),
  );
});

/// 工作流仓储。
final workflowRepositoryProvider = Provider<WorkflowRepository>((ref) {
  return WorkflowRepository(ref.watch(apiClientProvider));
});

/// 媒资仓储。
final assetRepositoryProvider = Provider<AssetRepository>((ref) {
  return AssetRepository(ref.watch(apiClientProvider));
});

/// 启动时根据本地令牌决定初始路由。
final initialRouteProvider = FutureProvider<String>((ref) async {
  final authRepo = ref.watch(authRepositoryProvider);
  final isAuth = await authRepo.isAuthenticated();
  return isAuth ? GinMamApp.pendingAuditRoute : GinMamApp.loginRoute;
});

/// 待审核列表状态（支持下拉刷新）。
final pendingAuditListProvider =
    AsyncNotifierProvider<PendingAuditListNotifier, List<WorkflowInstanceVO>>(
  PendingAuditListNotifier.new,
);

/// 待审核列表 [AsyncNotifier]。
class PendingAuditListNotifier extends AsyncNotifier<List<WorkflowInstanceVO>> {
  static const int _pageSize = 50;

  @override
  Future<List<WorkflowInstanceVO>> build() async {
    return _fetch();
  }

  /// 重新拉取待审列表。
  Future<void> refresh() async {
    state = const AsyncLoading();
    state = await AsyncValue.guard(_fetch);
  }

  Future<List<WorkflowInstanceVO>> _fetch() async {
    final repo = ref.read(workflowRepositoryProvider);
    final page = await repo.assignedToMe(page: 1, pageSize: _pageSize);
    return page.list;
  }
}
