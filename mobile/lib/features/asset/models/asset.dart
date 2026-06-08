/// 媒资视图模型，对齐 Web `assetApi.ts` 中的 AssetVO。
library;

/// 媒资视图对象。
class AssetVO {
  /// 媒资 ID。
  final String id;

  /// 标题。
  final String title;

  /// 类型（如 video、image）。
  final String type;

  /// 状态码。
  final int status;

  /// 描述。
  final String? description;

  /// 预览地址。
  final String? previewUrl;

  /// 创建时间。
  final String createdAt;

  /// 构造媒资 VO。
  const AssetVO({
    required this.id,
    required this.title,
    required this.type,
    required this.status,
    this.description,
    this.previewUrl,
    required this.createdAt,
  });

  /// 从 JSON 解析媒资。
  factory AssetVO.fromJson(Map<String, dynamic> json) {
    return AssetVO(
      id: json['id'] as String,
      title: json['title'] as String,
      type: json['type'] as String,
      status: json['status'] as int,
      description: json['description'] as String?,
      previewUrl: json['previewUrl'] as String?,
      createdAt: json['createdAt'] as String,
    );
  }
}
