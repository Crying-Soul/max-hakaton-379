// src/hooks/useMapEvents.ts
import { useState, useCallback, useRef } from 'react';
import {
  MapEvent,
  MapEventsParams,
  MapEventsResponse,
} from '../common/types/api';
import { apiService } from '../services/api.service';
import { AuthService } from '../services/auth.service';

interface UseMapEventsReturn {
  events: MapEvent[];
  loading: boolean;
  error: string | null;
  lastUpdated: Date | null;
  meta: MapEventsResponse['meta'] | null;
  fetchMapEvents: (params: MapEventsParams, force?: boolean) => Promise<void>;
  fetchUserMapEvents: (params: MapEventsParams) => Promise<void>;
  clearError: () => void;
  clearEvents: () => void;
}

export const useMapEvents = (): UseMapEventsReturn => {
  const [events, setEvents] = useState<MapEvent[]>([]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [lastUpdated, setLastUpdated] = useState<Date | null>(null);
  const [meta, setMeta] = useState<MapEventsResponse['meta'] | null>(null);

  const cacheRef = useRef<
    Map<
      string,
      { data: MapEvent[]; timestamp: number; meta: MapEventsResponse['meta'] }
    >
  >(new Map());
  const CACHE_TTL = 30000;

  const getCacheKey = (
    params: MapEventsParams,
    isUserEvents: boolean = false,
  ): string => {
    const user = AuthService.getUser();
    const base = isUserEvents ? `user_${user?.id}` : 'public';
    return `${base}_${params.lat}_${params.lon}_${params.radius_km}_${
      params.category_id?.join(',') || ''
    }`;
  };

  const fetchMapEvents = useCallback(
    async (params: MapEventsParams, force: boolean = false): Promise<void> => {
      const cacheKey = getCacheKey(params, false);
      const cached = cacheRef.current.get(cacheKey);

      if (!force && cached && Date.now() - cached.timestamp < CACHE_TTL) {
        setEvents(cached.data);
        setMeta(cached.meta);
        setLastUpdated(new Date(cached.timestamp));
        return;
      }

      setLoading(true);
      setError(null);

      try {
        const response = await apiService.getMapEvents(params);
        setEvents(response.data);
        setMeta(response.meta);
        setLastUpdated(new Date());

        cacheRef.current.set(cacheKey, {
          data: response.data,
          meta: response.meta,
          timestamp: Date.now(),
        });
      } catch (err: any) {
        const errorMessage = err.message || 'Failed to fetch events';
        setError(errorMessage);
        setEvents([]);
        setMeta(null);
      } finally {
        setLoading(false);
      }
    },
    [],
  );

  const fetchUserMapEvents = useCallback(
    async (params: MapEventsParams): Promise<void> => {
      const user = AuthService.getUser();
      if (!user || !AuthService.isAuthenticated()) {
        setError('User not authenticated');
        return;
      }

      const cacheKey = getCacheKey(params, true);

      setLoading(true);
      setError(null);

      try {
        const response = await apiService.getUserMapEvents(user.id, params);
        setEvents(response.data);
        setMeta(response.meta);
        setLastUpdated(new Date());

        cacheRef.current.set(cacheKey, {
          data: response.data,
          meta: response.meta,
          timestamp: Date.now(),
        });
      } catch (err: any) {
        const errorMessage = err.message || 'Failed to fetch user events';
        setError(errorMessage);

        // Если ошибка аутентификации, разлогиниваем пользователя
        if (err.message.includes('Authentication failed')) {
          AuthService.logout();
        }

        setEvents([]);
        setMeta(null);
      } finally {
        setLoading(false);
      }
    },
    [],
  );

  const clearError = useCallback(() => {
    setError(null);
  }, []);

  const clearEvents = useCallback(() => {
    setEvents([]);
    setMeta(null);
    cacheRef.current.clear();
  }, []);

  return {
    events,
    loading,
    error,
    lastUpdated,
    meta,
    fetchMapEvents,
    fetchUserMapEvents,
    clearError,
    clearEvents,
  };
};
