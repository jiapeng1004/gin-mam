/// 应用启动入口：初始化 Flutter 绑定并挂载 Riverpod 根。
library;

import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:gin_mam_mobile/app.dart';

/// 程序主入口。
void main() {
  WidgetsFlutterBinding.ensureInitialized();
  runApp(
    const ProviderScope(
      child: GinMamApp(),
    ),
  );
}