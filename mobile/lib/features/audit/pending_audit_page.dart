/// 待审核列表占位：任务 22+ 将对接审核 API。
library;

import 'package:flutter/material.dart';

/// 路由 `/audit/pending` 对应的待审核占位页。
class PendingAuditPage extends StatelessWidget {
  /// 创建待审核占位页。
  const PendingAuditPage({super.key});

  @override
  Widget build(BuildContext context) {
    return const Scaffold(
      body: Center(
        child: Text('待审核（任务 22 完善）'),
      ),
    );
  }
}