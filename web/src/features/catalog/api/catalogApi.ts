/**
 * 编目 API：对接后端 /api/v1/catalog/*（需 JWT）。
 */

import { del, get, post, put } from '../../../shared/api/client';

/** 编目树节点，对应后端 TreeNode */
export interface CatalogTreeNode {
  id: string;
  name: string;
  parentId?: string | null;
  children?: CatalogTreeNode[];
}

/** 编目详情，对应后端 CatalogVO */
export interface CatalogVO {
  id: string;
  name: string;
  parentId?: string | null;
  sortCode: number;
}

/** POST /catalog 创建请求体 */
export interface CatalogCreateRequest {
  name: string;
  parentId?: string | null;
  sortCode?: number;
}

/** PUT /catalog/:id 更新请求体 */
export interface CatalogUpdateRequest {
  name?: string;
  parentId?: string | null;
  sortCode?: number;
}

/** 编目配置项，对应后端 ConfigVO */
export interface CatalogConfigItem {
  id: string;
  configKey: string;
  configValue: string;
}

/** PUT /catalog/:id/config 请求体 items 元素 */
export interface CatalogConfigInput {
  configKey: string;
  configValue: string;
}

/**
 * 获取编目树。
 * 后端：GET /api/v1/catalog/tree，响应：TreeNode[]
 */
export async function fetchCatalogTree(): Promise<CatalogTreeNode[]> {
  return get<CatalogTreeNode[]>('/catalog/tree');
}

/**
 * 创建编目节点。
 * 后端：POST /api/v1/catalog
 */
export async function createCatalog(body: CatalogCreateRequest): Promise<CatalogVO> {
  return post<CatalogVO>('/catalog', body);
}

/**
 * 更新编目节点。
 * 后端：PUT /api/v1/catalog/:id
 */
export async function updateCatalog(id: string, body: CatalogUpdateRequest): Promise<CatalogVO> {
  return put<CatalogVO>(`/catalog/${id}`, body);
}

/**
 * 删除编目节点。
 * 后端：DELETE /api/v1/catalog/:id
 */
export async function deleteCatalog(id: string): Promise<void> {
  await del<unknown>(`/catalog/${id}`);
}

/**
 * 读取编目元数据配置列表。
 * 后端：GET /api/v1/catalog/:id/config
 */
export async function fetchCatalogConfig(catalogId: string): Promise<CatalogConfigItem[]> {
  return get<CatalogConfigItem[]>(`/catalog/${catalogId}/config`);
}

/**
 * 批量保存编目配置。
 * 后端：PUT /api/v1/catalog/:id/config，body：{ items: { configKey, configValue }[] }
 */
export async function putCatalogConfig(
  catalogId: string,
  items: CatalogConfigInput[],
): Promise<CatalogConfigItem[]> {
  return put<CatalogConfigItem[]>(`/catalog/${catalogId}/config`, { items });
}