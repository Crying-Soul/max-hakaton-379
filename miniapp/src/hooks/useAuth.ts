// src/hooks/useAuth.ts
import { useState, useCallback, useEffect } from 'react';
import { AuthService } from '../services/auth.service';
import { User } from '../common/types/api';

interface UseAuthReturn {
  isAuthenticated: boolean;
  user: User | null;
  loading: boolean;
  error: string | null;
  login: (initData: string) => Promise<boolean>;
  logout: () => void;
  clearError: () => void;
}

export const useAuth = (): UseAuthReturn => {
  const [isAuthenticated, setIsAuthenticated] = useState(false);
  const [user, setUser] = useState<User | null>(null);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  // Инициализация при загрузке и подписка на изменения
  useEffect(() => {
    const initializeAuth = () => {
      const token = AuthService.getToken();
      const userData = AuthService.getUser();

      if (token && userData) {
        if (AuthService.isAuthenticated()) {
          setIsAuthenticated(true);
          setUser(userData);
        } else {
          AuthService.logout();
          setIsAuthenticated(false);
          setUser(null);
        }
      }
    };

    initializeAuth();

    // Подписываемся на изменения аутентификации
    const unsubscribe = AuthService.onAuthChange((authStatus) => {
      setIsAuthenticated(authStatus);
      if (authStatus) {
        setUser(AuthService.getUser());
      } else {
        setUser(null);
      }
    });

    return unsubscribe;
  }, []);

  const login = useCallback(async (initData: string): Promise<boolean> => {
    setLoading(true);
    setError(null);

    try {
      const result = await AuthService.login(initData);

      if (!result) {
        throw new Error('Authentication failed');
      }

      // Состояние обновится через подписку onAuthChange
      return true;
    } catch (err: any) {
      const errorMessage = err?.message ?? 'Authentication failed';
      setError(errorMessage);
      return false;
    } finally {
      setLoading(false);
    }
  }, []);

  const logout = useCallback(() => {
    AuthService.logout();
    // Состояние обновится через подписку onAuthChange
  }, []);

  const clearError = useCallback(() => {
    setError(null);
  }, []);

  return {
    isAuthenticated,
    user,
    loading,
    error,
    login,
    logout,
    clearError,
  };
};
