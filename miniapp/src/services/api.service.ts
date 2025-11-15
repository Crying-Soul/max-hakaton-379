// src/services/api.service.ts
import axios, { AxiosInstance, AxiosRequestConfig } from 'axios';
import {
  MapEventsResponse,
  AuthSessionRequest,
  AuthSessionResponse,
  MapEventsParams,
} from '../common/types/api';

const BASE_URL = import.meta.env.VITE_BASE_URL;

class ApiError extends Error {
  constructor(message: string, public status?: number, public data?: any) {
    super(message);
    this.name = 'ApiError';
  }
}

export class ApiService {
  private client: AxiosInstance;
  private token: string | null = null;

  constructor(baseURL = BASE_URL) {
    this.client = axios.create({
      baseURL,
      timeout: 30000,
      headers: {
        Accept: 'application/json',
        'Content-Type': 'application/json',
      },
    });

    // Перехватчик запросов
    this.client.interceptors.request.use((config) => {
      if (this.token) {
        config.headers.Authorization = `Bearer ${this.token}`;
      }
      if (config.data instanceof FormData) {
        delete config.headers['Content-Type'];
      }
      return config;
    });

    // Перехватчик ответов
    this.client.interceptors.response.use(
      (response) => response,
      (error) => {
        if (axios.isAxiosError(error)) {
          const status = error.response?.status;
          const data = error.response?.data;
          let message = error.message;

          if (data && typeof data === 'object') {
            message = data.error?.message || data.message || message;
          }

          return Promise.reject(new ApiError(message, status, data));
        }
        return Promise.reject(error);
      },
    );
  }

  setToken(token: string | null) {
    this.token = token;
  }

  clearToken() {
    this.token = null;
  }

  private async request<T = any>(config: AxiosRequestConfig): Promise<T> {
    const response = await this.client.request<T>(config);
    return response.data;
  }

  async createSession(
    authData: AuthSessionRequest,
  ): Promise<AuthSessionResponse> {
    const sessionParams = new URLSearchParams();

    sessionParams.append('session', String(authData.initData));

    return this.request<AuthSessionResponse>({
      method: 'GET',
      url: '/auth/session',
      params: sessionParams,
    });
  }

  async getMapEvents(params: MapEventsParams): Promise<MapEventsResponse> {
    const searchParams = new URLSearchParams();

    if (params.lat !== undefined)
      searchParams.append('lat', String(params.lat));
    if (params.lon !== undefined)
      searchParams.append('lon', String(params.lon));
    if (params.radius_km !== undefined)
      searchParams.append('radius_km', String(params.radius_km));

    searchParams.append('limit', String(params.limit || 50));
    searchParams.append('offset', String(params.offset || 0));

    const categories = params.categories || params.category_id;
    if (categories?.length) {
      categories.forEach((category) =>
        searchParams.append('categories', String(category)),
      );
    }

    return this.request<MapEventsResponse>({
      method: 'GET',
      url: '/map/events',
      params: searchParams,
    });
  }

  async getUserMapEvents(
    userId: number | string,
    params: MapEventsParams,
  ): Promise<MapEventsResponse> {
    const searchParams = new URLSearchParams();

    if (params.lat !== undefined)
      searchParams.append('lat', String(params.lat));
    if (params.lon !== undefined)
      searchParams.append('lon', String(params.lon));
    if (params.radius_km !== undefined)
      searchParams.append('radius_km', String(params.radius_km));

    searchParams.append('limit', String(params.limit || 50));
    searchParams.append('offset', String(params.offset || 0));

    const categories = params.categories || params.category_id;
    if (categories?.length) {
      categories.forEach((category) =>
        searchParams.append('categories', String(category)),
      );
    }

    return this.request<MapEventsResponse>({
      method: 'GET',
      url: `/map/users/${userId}/events`,
      params: searchParams,
    });
  }

  async healthCheck(timeout = 5000): Promise<boolean> {
    try {
      await this.request({ method: 'GET', url: '/health', timeout });
      return true;
    } catch (error: any) {
      if (error instanceof ApiError && error.status && error.status < 500) {
        return true;
      }
      return false;
    }
  }
}

export const apiService = new ApiService();
export default apiService;
