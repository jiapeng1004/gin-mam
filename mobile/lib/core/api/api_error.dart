/// 后端业务错误类型，对应 httpx.ErrorVo：`err_code` / `err_msg`。
library;

/// 表示接口返回的业务错误（HTTP 4xx/5xx 且 body 为 ErrorVo 时抛出）。
class ApiError implements Exception {
  /// 业务错误码，对应 JSON 字段 `err_code`。
  final int errCode;

  /// 错误说明，对应 JSON 字段 `err_msg`。
  final String errMsg;

  /// 使用 [errCode] 与 [errMsg] 构造业务错误。
  ApiError(this.errCode, this.errMsg);

  @override
  String toString() => 'ApiError(errCode: $errCode, errMsg: $errMsg)';
}