/// 登录页：用户名/密码表单，成功后跳转待审列表。
library;

import 'package:antd_flutter_mobile/index.dart';
import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:gin_mam_mobile/app.dart';
import 'package:gin_mam_mobile/core/api/api_error.dart';
import 'package:gin_mam_mobile/core/providers.dart';

/// 路由 `/login` 对应的登录页。
class LoginPage extends ConsumerStatefulWidget {
  /// 创建登录页。
  const LoginPage({super.key});

  @override
  ConsumerState<LoginPage> createState() => _LoginPageState();
}

class _LoginPageState extends ConsumerState<LoginPage> {
  late final AntdInputController _usernameController;
  late final AntdInputController _passwordController;
  String? _errorText;

  @override
  void initState() {
    super.initState();
    _usernameController = AntdInputController()..text = 'admin';
    _passwordController = AntdInputController()..text = 'admin123';
  }

  @override
  void dispose() {
    _usernameController.dispose();
    _passwordController.dispose();
    super.dispose();
  }

  /// 提交登录表单。
  Future<void> _submit() async {
    final username = _usernameController.text.trim();
    final password = _passwordController.text;

    if (username.isEmpty || password.isEmpty) {
      setState(() => _errorText = '请输入用户名和密码');
      return;
    }

    setState(() => _errorText = null);

    try {
      await ref.read(authRepositoryProvider).login(
            username: username,
            password: password,
          );
      if (!mounted) return;
      Navigator.of(context).pushReplacementNamed(GinMamApp.pendingAuditRoute);
    } on ApiError catch (e) {
      final msg = e.errMsg.isNotEmpty ? e.errMsg : '登录失败';
      setState(() => _errorText = msg);
      if (mounted) {
        AntdToast.show(msg, toast: const AntdToast(type: AntdToastType.fail));
      }
    } catch (_) {
      const msg = '登录失败，请稍后重试';
      setState(() => _errorText = msg);
      if (mounted) {
        AntdToast.show(msg, toast: const AntdToast(type: AntdToastType.fail));
      }
    }
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      body: SafeArea(
        child: Center(
          child: SingleChildScrollView(
            padding: const EdgeInsets.symmetric(horizontal: 24),
            child: Column(
              mainAxisAlignment: MainAxisAlignment.center,
              crossAxisAlignment: CrossAxisAlignment.stretch,
              children: [
                Text(
                  '媒资审核',
                  textAlign: TextAlign.center,
                  style: Theme.of(context).textTheme.headlineSmall,
                ),
                const SizedBox(height: 32),
                const Text('用户名'),
                const SizedBox(height: 8),
                AntdInput(
                  controller: _usernameController,
                  placeholder: const Text('admin'),
                  keyboardType: TextInputType.text,
                ),
                const SizedBox(height: 16),
                const Text('密码'),
                const SizedBox(height: 8),
                AntdInput(
                  controller: _passwordController,
                  placeholder: const Text('******'),
                  obscureText: true,
                ),
                if (_errorText != null) ...[
                  const SizedBox(height: 12),
                  Text(
                    _errorText!,
                    style: TextStyle(color: Theme.of(context).colorScheme.error),
                  ),
                ],
                const SizedBox(height: 24),
                AntdButton(
                  key: const Key('login_submit_button'),
                  block: true,
                  behavior: HitTestBehavior.opaque,
                  onLoadingTap: _submit,
                  child: const Text('登录'),
                ),
              ],
            ),
          ),
        ),
      ),
    );
  }
}
