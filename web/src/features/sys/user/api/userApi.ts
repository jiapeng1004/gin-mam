/**
 * 系统用户 API：对接后端 GET /api/v1/sys/user/page（需 JWT）。
 */

import { get } from '../../../../shared/api/client';

/** 用户列表项，对应后端 UserVO */
export interface UserVO {
  id: string;
  username: string;
  nickname?: string;
  orgId?: string | null;
  status: number;
  createdAt: string;
  updatedAt: string;
}

/** 分页响应，对应后端 UserPageResult */
export interface UserPageResponse {
  list: UserVO[];
  total: number;
  page: number;
  pageSize: number;
}

/**
 * 分页查询系统用户。
 * @param page - 页码，从 1 开始
 * @param pageSize - 每页条数
 * @returns list、total、page、pageSize
 */
export async function fetchUserPage(page: number, pageSize: number): Promise<UserPageResponse> {
  return get<UserPageResponse>('/sys/user/page', { page, pageSize });
}