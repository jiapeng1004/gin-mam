/// 应用启动入口：初始化 Flutter 绑定并挂载 Riverpod 根。
library;

import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:gin_mam_mobile/app.dart';
import 'package:gin_mam_mobile/core/providers.dart';
import 'package:shared_preferences/shared_preferences.dart';

/// 程序主入口。
Future<void> main() async {
  WidgetsFlutterBinding.ensureInitialized();
  final prefs = await SharedPreferences.getInstance();

  runApp(
    ProviderScope(
      overrides: [
        sharedPreferencesProvider.overrideWithValue(prefs),
      ],
      child: const GinMamApp(),
    ),
  );
}
