/// 媒资仓储：按 ID 查询媒资详情。
library;

import 'package:gin_mam_mobile/core/api/api_client.dart';
import 'package:gin_mam_mobile/features/asset/models/asset.dart';

/// 媒资数据访问层。
class AssetRepository {
  final ApiClient _apiClient;

  /// 使用 [apiClient] 构造。
  AssetRepository(this._apiClient);

  /// 根据 ID 获取媒资详情。
  Future<AssetVO> getById(String id) async {
    final response = await _apiClient.get<Map<String, dynamic>>('/asset/$id');
    return AssetVO.fromJson(response);
  }
}
