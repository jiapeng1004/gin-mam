/// 应用根组件：MaterialApp、Ant Design Mobile 主题壳与命名路由。
library;

import 'package:antd_flutter_mobile/index.dart';
import 'package:flutter/material.dart';
import 'package:gin_mam_mobile/features/audit/pending_audit_page.dart';
import 'package:gin_mam_mobile/features/auth/login_page.dart';

/// Gin MAM 移动端应用入口 Widget。
///
/// 使用 [antd_flutter_mobile](https://pub.dev/packages/antd_flutter_mobile) 提供
/// 设计令牌与浮层能力；该库处于 **alpha** 阶段，组件 API 可能变更，集成前请评估风险
/// （详见 `mobile/README.md`）。
class GinMamApp extends StatelessWidget {
  /// 创建应用根组件。
  const GinMamApp({super.key});

  /// 登录页路由路径。
  static const String loginRoute = '/login';

  /// 待审核列表路由路径。
  static const String pendingAuditRoute = '/audit/pending';

  @override
  Widget build(BuildContext context) {
    return MaterialApp(
      title: 'Gin MAM',
      navigatorObservers: [AntdLayer.observer],
      builder: (context, child) {
        return AntdTokenBuilder(
          builder: (context, token) => child ?? const SizedBox.shrink(),
        );
      },
      initialRoute: loginRoute,
      routes: {
        loginRoute: (_) => const LoginPage(),
        pendingAuditRoute: (_) => const PendingAuditPage(),
      },
    );
  }
}