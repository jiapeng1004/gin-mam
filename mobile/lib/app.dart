/// 应用根组件：MaterialApp、Ant Design Mobile 主题壳与命名路由。
library;

import 'package:antd_flutter_mobile/index.dart';
import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:gin_mam_mobile/core/providers.dart';
import 'package:gin_mam_mobile/features/audit/audit_detail_args.dart';
import 'package:gin_mam_mobile/features/audit/audit_detail_page.dart';
import 'package:gin_mam_mobile/features/audit/pending_audit_page.dart';
import 'package:gin_mam_mobile/features/auth/login_page.dart';

/// Gin MAM 移动端应用入口 Widget。
///
/// 使用 [antd_flutter_mobile](https://pub.dev/packages/antd_flutter_mobile) 提供
/// 设计令牌与浮层能力；该库处于 **alpha** 阶段，组件 API 可能变更，集成前请评估风险
/// （详见 `mobile/README.md`）。
class GinMamApp extends ConsumerWidget {
  /// 创建应用根组件。
  const GinMamApp({super.key});

  /// 登录页路由路径。
  static const String loginRoute = '/login';

  /// 待审核列表路由路径。
  static const String pendingAuditRoute = '/audit/pending';

  /// 审核详情路由路径。
  static const String auditDetailRoute = '/audit/detail';

  /// 根据 [settings] 构建命名路由。
  static Route<dynamic>? onGenerateRoute(RouteSettings settings) {
    switch (settings.name) {
      case loginRoute:
        return MaterialPageRoute<void>(
          settings: settings,
          builder: (_) => const LoginPage(),
        );
      case pendingAuditRoute:
        return MaterialPageRoute<void>(
          settings: settings,
          builder: (_) => const PendingAuditPage(),
        );
      case auditDetailRoute:
        final args = settings.arguments;
        if (args is! AuditDetailArgs) {
          return MaterialPageRoute<void>(
            settings: settings,
            builder: (_) => const Scaffold(
              body: Center(child: Text('缺少审核参数')),
            ),
          );
        }
        return MaterialPageRoute<void>(
          settings: settings,
          builder: (_) => AuditDetailPage(args: args),
        );
      default:
        return MaterialPageRoute<void>(
          settings: settings,
          builder: (_) => const LoginPage(),
        );
    }
  }

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final initialRoute = ref.watch(initialRouteProvider);
    final navigatorKey = ref.watch(navigatorKeyProvider);

    return initialRoute.when(
      loading: () => MaterialApp(
        title: 'Gin MAM',
        home: const Scaffold(
          body: Center(child: CircularProgressIndicator()),
        ),
      ),
      error: (error, stackTrace) => _buildApp(
        navigatorKey: navigatorKey,
        initialRoute: loginRoute,
      ),
      data: (route) => _buildApp(
        navigatorKey: navigatorKey,
        initialRoute: route,
      ),
    );
  }

  Widget _buildApp({
    required GlobalKey<NavigatorState> navigatorKey,
    required String initialRoute,
  }) {
    return MaterialApp(
      title: 'Gin MAM',
      navigatorKey: navigatorKey,
      navigatorObservers: [AntdLayer.observer],
      builder: (context, child) {
        return AntdTokenBuilder(
          builder: (context, token) => child ?? const SizedBox.shrink(),
        );
      },
      initialRoute: initialRoute,
      onGenerateRoute: onGenerateRoute,
    );
  }
}
