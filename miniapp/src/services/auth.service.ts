// src/services/auth.service.ts
import { apiService } from './api.service';
import { User } from '../common/types/api';

type SessionStorageValue = {
  token: string;
  expiresAt: number; // ms since epoch
  user: User;
};

export class AuthService {
  private static readonly STORAGE_KEY = 'max_auth_session';
  private static cache: SessionStorageValue | null = null;
  private static authListeners: ((isAuthenticated: boolean) => void)[] = [];

  private static readStorage(): SessionStorageValue | null {
    if (typeof window === 'undefined') return null;
    try {
      const raw = sessionStorage.getItem(this.STORAGE_KEY);
      if (!raw) return null;
      const parsed = JSON.parse(raw) as SessionStorageValue;
      if (!parsed.token || !parsed.expiresAt) return null;
      return parsed;
    } catch {
      return null;
    }
  }

  private static writeStorage(value: SessionStorageValue | null) {
    if (typeof window === 'undefined') return;
    try {
      if (value === null) {
        sessionStorage.removeItem(this.STORAGE_KEY);
        this.cache = null;
      } else {
        sessionStorage.setItem(this.STORAGE_KEY, JSON.stringify(value));
        this.cache = value;
      }
      this.notifyAuthChange();
    } catch {}
  }

  static async login(initData: string): Promise<{ token: string; user: User }> {
    try {
      const response = await apiService.createSession({ initData });

      const { token, expiresIn, user } = response;
      const expiresAt = Date.now() + (Number(expiresIn) || 0) * 1000;

      const payload = { token, expiresAt, user };
      this.writeStorage(payload);
      apiService.setToken(token);

      return { token, user };
    } catch (e) {
      this.clearAuth();
      return Promise.reject(e);
    }
  }

  static logout(): void {
    this.clearAuth();
    apiService.clearToken();
  }

  static getToken(): string | null {
    if (this.cache) {
      if (this.cache.expiresAt > Date.now()) return this.cache.token;
      this.clearAuth();
      return null;
    }

    const stored = this.readStorage();
    if (!stored) return null;
    if (stored.expiresAt <= Date.now()) {
      this.clearAuth();
      return null;
    }
    this.cache = stored;
    return stored.token;
  }

  static getUser(): User | null {
    if (this.cache?.user) return this.cache.user;
    const stored = this.readStorage();
    return stored?.user ?? null;
  }

  static isAuthenticated(): boolean {
    const token = this.getToken();
    if (!token) return false;

    const payload = this.decodeTokenPayload(token);
    if (!payload) return true;
    if (payload.exp) {
      const expMs = payload.exp * 1000;
      return expMs > Date.now();
    }
    return true;
  }

  static initializeAuth(): void {
    const stored = this.readStorage();
    if (stored && stored.expiresAt > Date.now()) {
      this.cache = stored;
      apiService.setToken(stored.token);
    } else {
      this.clearAuth();
    }
    this.notifyAuthChange();
  }

  // Для отслеживания изменений аутентификации
  static onAuthChange(
    callback: (isAuthenticated: boolean) => void,
  ): () => void {
    this.authListeners.push(callback);
    // Immediately call with current state
    callback(this.isAuthenticated());
    return () => {
      this.authListeners = this.authListeners.filter((cb) => cb !== callback);
    };
  }

  private static notifyAuthChange(): void {
    const isAuthenticated = this.isAuthenticated();
    this.authListeners.forEach((callback) => callback(isAuthenticated));
  }

  private static clearAuth(): void {
    this.writeStorage(null);
    apiService.clearToken();
  }

  private static decodeTokenPayload(token: string): any | null {
    try {
      const parts = token.split('.');
      if (parts.length < 2) return null;
      const raw = parts[1];
      const b64 = raw.replace(/-/g, '+').replace(/_/g, '/');
      const pad = b64.length % 4 === 0 ? '' : '='.repeat(4 - (b64.length % 4));
      const decoded = atob(b64 + pad);
      return JSON.parse(decoded);
    } catch {
      return null;
    }
  }
}
