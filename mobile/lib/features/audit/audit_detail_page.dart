/// 审核详情页：展示媒资信息，支持通过/打回并填写备注。
library;

import 'package:antd_flutter_mobile/index.dart';
import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:gin_mam_mobile/core/api/api_error.dart';
import 'package:gin_mam_mobile/core/providers.dart';
import 'package:gin_mam_mobile/features/asset/models/asset.dart';
import 'package:gin_mam_mobile/features/audit/audit_detail_args.dart';
import 'package:gin_mam_mobile/features/audit/models/workflow_instance.dart';

/// 路由 `/audit/detail` 对应的审核详情页。
class AuditDetailPage extends ConsumerStatefulWidget {
  /// 审核详情路由参数。
  final AuditDetailArgs args;

  /// 创建审核详情页。
  const AuditDetailPage({
    super.key,
    required this.args,
  });

  @override
  ConsumerState<AuditDetailPage> createState() => _AuditDetailPageState();
}

class _AuditDetailPageState extends ConsumerState<AuditDetailPage> {
  final _remarkController = AntdInputController();
  AssetVO? _asset;
  bool _loading = true;
  bool _submitting = false;
  String? _errorText;

  @override
  void initState() {
    super.initState();
    _loadAsset();
  }

  @override
  void dispose() {
    _remarkController.dispose();
    super.dispose();
  }

  /// 加载关联媒资详情。
  Future<void> _loadAsset() async {
    setState(() {
      _loading = true;
      _errorText = null;
    });
    try {
      final asset = await ref
          .read(assetRepositoryProvider)
          .getById(widget.args.assetId);
      if (!mounted) return;
      setState(() {
        _asset = asset;
        _loading = false;
      });
    } on ApiError catch (e) {
      if (!mounted) return;
      setState(() {
        _errorText = e.errMsg.isNotEmpty ? e.errMsg : '加载媒资失败';
        _loading = false;
      });
    } catch (_) {
      if (!mounted) return;
      setState(() {
        _errorText = '加载媒资失败';
        _loading = false;
      });
    }
  }

  /// 提交审核结果。
  Future<void> _submitAudit({required bool pass}) async {
    setState(() => _submitting = true);
    try {
      await ref.read(workflowRepositoryProvider).audit(
            AuditWorkflowRequest(
              instanceId: widget.args.instanceId,
              pass: pass,
              remark: _remarkController.text.trim(),
            ),
          );
      if (!mounted) return;
      AntdToast.show(
        pass ? '已通过' : '已打回',
        toast: AntdToast(type: pass ? AntdToastType.success : AntdToastType.normal),
      );
      Navigator.of(context).pop(true);
    } on ApiError catch (e) {
      if (!mounted) return;
      AntdToast.show(
        e.errMsg.isNotEmpty ? e.errMsg : '操作失败',
        toast: const AntdToast(type: AntdToastType.fail),
      );
    } catch (_) {
      if (!mounted) return;
      AntdToast.show('操作失败', toast: const AntdToast(type: AntdToastType.fail));
    } finally {
      if (mounted) {
        setState(() => _submitting = false);
      }
    }
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(title: const Text('审核详情')),
      body: _loading
          ? const Center(child: CircularProgressIndicator())
          : _errorText != null
              ? Center(
                  child: Column(
                    mainAxisAlignment: MainAxisAlignment.center,
                    children: [
                      Text(_errorText!),
                      const SizedBox(height: 16),
                      AntdButton(
                        onTap: _loadAsset,
                        child: const Text('重试'),
                      ),
                    ],
                  ),
                )
              : SingleChildScrollView(
                  padding: const EdgeInsets.all(16),
                  child: Column(
                    crossAxisAlignment: CrossAxisAlignment.stretch,
                    children: [
                      _InfoRow(label: '标题', value: _asset!.title),
                      _InfoRow(label: '类型', value: _asset!.type),
                      _InfoRow(label: '状态', value: '${_asset!.status}'),
                      if (_asset!.description != null &&
                          _asset!.description!.isNotEmpty)
                        _InfoRow(label: '描述', value: _asset!.description!),
                      const SizedBox(height: 16),
                      const Text('审核备注'),
                      const SizedBox(height: 8),
                      AntdTextArea(
                        controller: _remarkController,
                        placeholder: const Text('可选，打回时建议填写原因'),
                        minLines: 3,
                        maxLines: 5,
                      ),
                      const SizedBox(height: 24),
                      AntdButton(
                        block: true,
                        loading: _submitting,
                        onTap: _submitting ? null : () => _submitAudit(pass: true),
                        child: const Text('通过'),
                      ),
                      const SizedBox(height: 12),
                      AntdButton(
                        block: true,
                        color: AntdColor.danger,
                        fill: AntdButtonFill.outline,
                        loading: _submitting,
                        onTap: _submitting ? null : () => _submitAudit(pass: false),
                        child: const Text('打回'),
                      ),
                    ],
                  ),
                ),
    );
  }
}

/// 详情页键值行展示。
class _InfoRow extends StatelessWidget {
  final String label;
  final String value;

  const _InfoRow({required this.label, required this.value});

  @override
  Widget build(BuildContext context) {
    return Padding(
      padding: const EdgeInsets.only(bottom: 8),
      child: Row(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          SizedBox(
            width: 56,
            child: Text(
              label,
              style: Theme.of(context).textTheme.bodySmall,
            ),
          ),
          Expanded(child: Text(value)),
        ],
      ),
    );
  }
}
