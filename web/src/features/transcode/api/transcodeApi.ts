/**
 * 转码 API：转码组分页 CRUD 与转码任务创建、分页、详情查询。
 */

import { del, get, post, put } from '../../../shared/api/client';

/** 转码组 VO */
export interface TranscodeGroupVO {
  id: string;
  name: string;
  message?: string;
  groupType: number;
  param?: string;
  strategyType: number;
  defaultFlag: number;
  availableFlag: number;
  createdBy: string;
  createdAt: string;
  updatedAt: string;
}

/** 转码组分页结果 */
export interface TranscodeGroupPageResult {
  list: TranscodeGroupVO[];
  total: number;
  page: number;
  pageSize: number;
}

/** 创建转码组请求 */
export interface CreateTranscodeGroupRequest {
  name: string;
  message?: string;
  groupType?: number;
  param?: string;
  strategyType?: number;
  defaultFlag?: number;
  availableFlag?: number;
}

/** 更新转码组请求（字段均可选） */
export interface UpdateTranscodeGroupRequest {
  name?: string;
  message?: string;
  groupType?: number;
  param?: string;
  strategyType?: number;
  defaultFlag?: number;
  availableFlag?: number;
}

/** 转码任务 VO */
export interface TranscodeTaskVO {
  id: string;
  assetId: string;
  assetFileId?: string | null;
  transcodeGroupId: string;
  profileId?: string | null;
  status: number;
  externalJobId?: string;
  outputPath?: string;
  errorMsg?: string;
  createdBy: string;
  createdAt: string;
  updatedAt: string;
}

/** 转码任务分页结果 */
export interface TranscodeTaskPageResult {
  list: TranscodeTaskVO[];
  total: number;
  page: number;
  pageSize: number;
}

/** 创建转码任务请求 */
export interface CreateTranscodeTaskRequest {
  assetId: string;
  assetFileId?: string;
  transcodeGroupId: string;
}

/**
 * 分页查询转码组。
 */
export async function fetchTranscodeGroupPage(
  page: number,
  pageSize: number,
): Promise<TranscodeGroupPageResult> {
  return post<TranscodeGroupPageResult>('/transcode/group/page', { page, pageSize });
}

/**
 * 创建转码组。
 */
export async function createTranscodeGroup(body: CreateTranscodeGroupRequest): Promise<TranscodeGroupVO> {
  return post<TranscodeGroupVO>('/transcode/group', body);
}

/**
 * 更新转码组。
 * @param id - 转码组 ID
 */
export async function updateTranscodeGroup(
  id: string,
  body: UpdateTranscodeGroupRequest,
): Promise<TranscodeGroupVO> {
  return put<TranscodeGroupVO>(`/transcode/group/${id}`, body);
}

/**
 * 删除转码组。
 * @param id - 转码组 ID
 */
export async function deleteTranscodeGroup(id: string): Promise<void> {
  await del<void>(`/transcode/group/${id}`);
}

/**
 * 创建转码任务。
 */
export async function createTranscodeTask(body: CreateTranscodeTaskRequest): Promise<TranscodeTaskVO> {
  return post<TranscodeTaskVO>('/transcode/task', body);
}

/**
 * 分页查询转码任务。
 */
export async function fetchTranscodeTaskPage(
  page: number,
  pageSize: number,
): Promise<TranscodeTaskPageResult> {
  return post<TranscodeTaskPageResult>('/transcode/task/page', { page, pageSize });
}

/**
 * 按 ID 查询转码任务详情。
 * @param id - 任务 ID
 */
export async function getTranscodeTask(id: string): Promise<TranscodeTaskVO> {
  return get<TranscodeTaskVO>(`/transcode/task/${id}`);
}