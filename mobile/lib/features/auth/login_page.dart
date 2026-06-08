/// 登录页占位：任务 22 将接入表单与鉴权流程。
library;

import 'package:flutter/material.dart';

/// 路由 `/login` 对应的登录占位页。
class LoginPage extends StatelessWidget {
  /// 创建登录占位页。
  const LoginPage({super.key});

  @override
  Widget build(BuildContext context) {
    return const Scaffold(
      body: Center(
        child: Text('登录（任务 22 完善）'),
      ),
    );
  }
}