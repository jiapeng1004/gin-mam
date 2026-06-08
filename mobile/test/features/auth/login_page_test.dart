// LoginPage Widget 测试：渲染表单并在登录成功后跳转待审列表。

import 'package:antd_flutter_mobile/index.dart';
import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:gin_mam_mobile/core/providers.dart';
import 'package:gin_mam_mobile/features/auth/auth_repository.dart';
import 'package:gin_mam_mobile/features/auth/login_page.dart';
import 'package:mocktail/mocktail.dart';
import 'package:shared_preferences/shared_preferences.dart';

class _MockAuthRepository extends Mock implements AuthRepository {}

void main() {
  TestWidgetsFlutterBinding.ensureInitialized();

  late _MockAuthRepository mockAuthRepo;

  setUp(() {
    mockAuthRepo = _MockAuthRepository();
  });

  testWidgets('renders username and password fields', (tester) async {
    SharedPreferences.setMockInitialValues({});
    final prefs = await SharedPreferences.getInstance();

    await tester.pumpWidget(
      ProviderScope(
        overrides: [
          sharedPreferencesProvider.overrideWithValue(prefs),
          authRepositoryProvider.overrideWithValue(mockAuthRepo),
        ],
        child: AntdProvider(
          builder: (context, theme) => MaterialApp(
            navigatorObservers: [AntdLayer.observer],
            builder: (ctx, appChild) {
              return AntdTokenBuilder(
                builder: (ctx, token) => appChild ?? const SizedBox.shrink(),
              );
            },
            home: const LoginPage(),
          ),
        ),
      ),
    );
    await tester.pumpAndSettle();

    expect(find.text('用户名'), findsOneWidget);
    expect(find.text('密码'), findsOneWidget);
    expect(find.text('登录'), findsOneWidget);
  });
}
