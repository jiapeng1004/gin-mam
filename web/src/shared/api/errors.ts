/**
 * API 业务错误类型：与后端 httpx.ErrorVo（err_code / err_msg）在前端侧的映射。
 */

/**
 * 表示后端返回的业务错误（HTTP 4xx/5xx 且 body 为 ErrorVo）。
 */
export class ApiError extends Error {
  /** 业务错误码，对应 JSON 字段 err_code */
  readonly errCode: number;
  /** 错误说明，对应 JSON 字段 err_msg */
  readonly errMsg: string;

  /**
   * @param errCode - 业务错误码
   * @param errMsg - 错误说明
   */
  constructor(errCode: number, errMsg: string) {
    super(errMsg);
    this.name = 'ApiError';
    this.errCode = errCode;
    this.errMsg = errMsg;
  }
}
