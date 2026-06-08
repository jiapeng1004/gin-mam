/**
 * 工作流 API：对接后端 asset-workflow 与 workflow/def 相关 POST 接口。
 */

import { post } from '../../../shared/api/client';

/** 分页请求体 */
export interface PageRequest {
  page: number;
  pageSize: number;
}

/** 流程实例 VO */
export interface WorkflowInstanceVO {
  id: string;
  workflowDefId: string;
  assetId: string;
  level: number;
  auditLevel: number;
  auditStatus: number;
  createdBy: string;
  createdAt: string;
  updatedAt: string;
}

/** 实例分页结果 */
export interface WorkflowInstancePageResult {
  list: WorkflowInstanceVO[];
  total: number;
  page: number;
  pageSize: number;
}

/** 流程定义层级审核人 */
export interface WorkflowDefLevelUserVO {
  level: number;
  userId: string;
}

/** 流程定义 VO */
export interface WorkflowDefVO {
  id: string;
  name: string;
  auditLevel: number;
  status: number;
  levelUsers?: WorkflowDefLevelUserVO[];
  createdAt: string;
  updatedAt: string;
}

/** 流程定义分页结果 */
export interface WorkflowDefPageResult {
  list: WorkflowDefVO[];
  total: number;
  page: number;
  pageSize: number;
}

/** 提交审核请求 */
export interface SubmitWorkflowRequest {
  assetId: string;
  workflowDefId: string;
}

/** 审核请求 */
export interface AuditWorkflowRequest {
  instanceId: string;
  pass: boolean;
  remark?: string;
}

/** 创建流程定义请求 */
export interface CreateWorkflowDefRequest {
  name: string;
  auditLevel: number;
  levelUsers: WorkflowDefLevelUserVO[];
}

/**
 * 提交媒资进入审核流程。
 * @param body - 媒资 ID 与流程定义 ID
 * @returns 新建的流程实例
 */
export async function submitWorkflow(body: SubmitWorkflowRequest): Promise<WorkflowInstanceVO> {
  return post<WorkflowInstanceVO>('/asset-workflow/submit', body);
}

/**
 * 对指定实例执行通过或打回审核。
 * @param body - 实例 ID、是否通过及备注
 */
export async function auditWorkflow(body: AuditWorkflowRequest): Promise<void> {
  await post<void>('/asset-workflow/audit', body);
}

/**
 * 分页查询待当前用户审核的实例列表。
 */
export async function fetchAssignedToMe(page: number, pageSize: number): Promise<WorkflowInstancePageResult> {
  return post<WorkflowInstancePageResult>('/asset-workflow/assigned-to-me', { page, pageSize });
}

/**
 * 分页查询当前用户发起的实例列表。
 */
export async function fetchCreatedByMe(page: number, pageSize: number): Promise<WorkflowInstancePageResult> {
  return post<WorkflowInstancePageResult>('/asset-workflow/created-by-me', { page, pageSize });
}

/**
 * 分页查询当前用户已审核过的实例列表。
 */
export async function fetchAuditedByMe(page: number, pageSize: number): Promise<WorkflowInstancePageResult> {
  return post<WorkflowInstancePageResult>('/asset-workflow/audited-by-me', { page, pageSize });
}

/**
 * 分页查询全部流程实例（管理员视角）。
 */
export async function fetchAllInstances(page: number, pageSize: number): Promise<WorkflowInstancePageResult> {
  return post<WorkflowInstancePageResult>('/asset-workflow/all', { page, pageSize });
}

/**
 * 分页查询流程定义列表。
 */
export async function fetchWorkflowDefPage(page: number, pageSize: number): Promise<WorkflowDefPageResult> {
  return post<WorkflowDefPageResult>('/workflow/def/page', { page, pageSize });
}

/**
 * 创建新的流程定义。
 * @param body - 名称、审核层级与各级审核人
 */
export async function createWorkflowDef(body: CreateWorkflowDefRequest): Promise<WorkflowDefVO> {
  return post<WorkflowDefVO>('/workflow/def', body);
}