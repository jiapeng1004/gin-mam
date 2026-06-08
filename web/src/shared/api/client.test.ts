/**
 * apiClient 与 REST 封装单测：验证无包装 2xx、ErrorVo 4xx 与 Authorization 头。
 */

import { AxiosError, type AxiosResponse, type InternalAxiosRequestConfig } from 'axios';
import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest';
import { AUTH_TOKEN_KEY } from '../../app/auth';
import { apiClient, get } from './client';
import { ApiError } from './errors';

type AdapterConfig = InternalAxiosRequestConfig;

function successAdapter(data: unknown) {
  return async (config: AdapterConfig): Promise<AxiosResponse> => ({
    data,
    status: 200,
    statusText: 'OK',
    headers: {},
    config,
  });
}

function errorAdapter(
  status: number,
  body: { err_code: number; err_msg: string },
) {
  return async (config: AdapterConfig): Promise<AxiosResponse> => {
    const error = new AxiosError('Request failed', String(status), config);
    error.response = {
      data: body,
      status,
      statusText: 'Error',
      headers: {},
      config,
    };
    throw error;
  };
}

describe('apiClient', () => {
  const originalAdapter = apiClient.defaults.adapter;

  beforeEach(() => {
    localStorage.clear();
  });

  afterEach(() => {
    apiClient.defaults.adapter = originalAdapter;
    vi.restoreAllMocks();
  });

  it('returns business payload on 200 without wrapper', async () => {
    apiClient.defaults.adapter = successAdapter({ id: 42, name: 'asset' });

    const result = await get<{ id: number; name: string }>('/assets/42');

    expect(result).toEqual({ id: 42, name: 'asset' });
  });

  it('throws ApiError with errCode from ErrorVo on 404', async () => {
    apiClient.defaults.adapter = errorAdapter(404, {
      err_code: 40401,
      err_msg: '资源不存在',
    });

    await expect(get('/missing')).rejects.toMatchObject({
      errCode: 40401,
      errMsg: '资源不存在',
    });
    await expect(get('/missing')).rejects.toBeInstanceOf(ApiError);
  });

  it('sends Authorization header when token is stored', async () => {
    localStorage.setItem(AUTH_TOKEN_KEY, 'test-jwt-token');

    apiClient.defaults.adapter = async (config: AdapterConfig) => {
      const auth = config.headers.get('Authorization');
      expect(auth).toBe('Bearer test-jwt-token');
      return {
        data: { ok: true },
        status: 200,
        statusText: 'OK',
        headers: {},
        config,
      };
    };

    await get('/me');
  });
});
