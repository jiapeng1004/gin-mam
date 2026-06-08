/**
 * 媒资 API：对接后端 /api/v1/asset/*（需 JWT）。
 */

import { apiClient, del, get, post, put } from '../../../shared/api/client';

/** 媒资视图，对应后端 AssetVO */
export interface AssetVO {
  id: string;
  title: string;
  type: string;
  status: number;
  catalogId?: string | null;
  description?: string;
  previewUrl?: string;
  storagePath?: string;
  metadata?: Record<string, string>;
  createdBy: string;
  createdAt: string;
  updatedAt: string;
}

/** POST /asset/page 请求体 */
export interface AssetPageRequest {
  page: number;
  pageSize: number;
  keyword?: string;
}

/** POST /asset/page 响应，对应 PageResult */
export interface AssetPageResponse {
  list: AssetVO[];
  total: number;
  page: number;
  pageSize: number;
}

/** POST /asset 创建请求体 */
export interface AssetCreateRequest {
  title: string;
  type: string;
  catalogId?: string | null;
  description?: string;
  storagePath: string;
  mimeType?: string;
  fileSize?: number;
  metadata?: Record<string, string>;
  createdBy?: string;
}

/** PUT /asset/:id 更新请求体 */
export interface AssetUpdateRequest {
  title?: string;
  catalogId?: string | null;
  description?: string;
  status?: number;
  metadata?: Record<string, string>;
}

/** POST /asset/upload/init 请求与响应 */
export interface UploadInitRequest {
  fileName: string;
  fileSize: number;
  mimeType?: string;
  chunkSize?: number;
}

export interface UploadInitResponse {
  uploadId: string;
  chunkSize: number;
}

/** POST /asset/upload/complete 请求体 */
export interface UploadCompleteRequest {
  uploadId: string;
  assetTitle: string;
  catalogId?: string | null;
  type: string;
}

/**
 * 分页查询媒资。
 * 后端：POST /api/v1/asset/page，body：page、pageSize、keyword?
 */
export async function fetchAssetPage(params: AssetPageRequest): Promise<AssetPageResponse> {
  return post<AssetPageResponse>('/asset/page', params);
}

/**
 * 按 ID 获取媒资详情（含 previewUrl）。
 * 后端：GET /api/v1/asset/:id
 */
export async function fetchAssetById(id: string): Promise<AssetVO> {
  return get<AssetVO>(`/asset/${id}`);
}

/**
 * 创建媒资记录。
 * 后端：POST /api/v1/asset
 */
export async function createAsset(body: AssetCreateRequest): Promise<AssetVO> {
  return post<AssetVO>('/asset', body);
}

/**
 * 更新媒资。
 * 后端：PUT /api/v1/asset/:id
 */
export async function updateAsset(id: string, body: AssetUpdateRequest): Promise<AssetVO> {
  return put<AssetVO>(`/asset/${id}`, body);
}

/**
 * 删除媒资（GORM 软删 deleted_at）。
 * 后端：DELETE /api/v1/asset/:id
 */
export async function deleteAsset(id: string): Promise<void> {
  await del<unknown>(`/asset/${id}`);
}

/**
 * 初始化分片上传会话。
 * 后端：POST /api/v1/asset/upload/init
 */
export async function initAssetUpload(body: UploadInitRequest): Promise<UploadInitResponse> {
  return post<UploadInitResponse>('/asset/upload/init', body);
}

/**
 * 上传单个分片（multipart）。
 * 后端：POST /api/v1/asset/upload/chunk，表单：uploadId、chunkIndex、file
 */
export async function uploadAssetChunk(
  uploadId: string,
  chunkIndex: number,
  chunk: Blob,
  fileName = 'chunk.bin',
): Promise<void> {
  const form = new FormData();
  form.append('uploadId', uploadId);
  form.append('chunkIndex', String(chunkIndex));
  form.append('file', chunk, fileName);
  await apiClient.post('/asset/upload/chunk', form, {
    headers: { 'Content-Type': 'multipart/form-data' },
  });
}

/**
 * 合并分片并创建媒资。
 * 后端：POST /api/v1/asset/upload/complete
 */
export async function completeAssetUpload(body: UploadCompleteRequest): Promise<AssetVO> {
  return post<AssetVO>('/asset/upload/complete', body);
}