// src/store/map.store.ts
import { create } from 'zustand';
import { MapEventsParams } from '../common/types/api';

interface MapState {
  mapParams: MapEventsParams;
  lastFetchTime: Date | null;
  isUserEvents: boolean;
  actions: {
    setMapParams: (params: Partial<MapEventsParams>) => void;
    setCenter: (lat: number, lon: number) => void;
    setRadius: (radius_km: number) => void;
    setCategories: (categoryIds: number[]) => void;
    setLastFetchTime: (time: Date) => void;
    setIsUserEvents: (isUserEvents: boolean) => void;
    resetMapParams: () => void;
  };
}

const defaultParams: MapEventsParams = {
  lat: 59.935,
  lon: 30.325,
  radius_km: 10,
  limit: 50,
  offset: 0,
};

export const useMapStore = create<MapState>((set) => ({
  mapParams: defaultParams,
  lastFetchTime: null,
  isUserEvents: false,
  actions: {
    setMapParams: (params) => 
      set((state) => ({ 
        mapParams: { ...state.mapParams, ...params } 
      })),
    
    setCenter: (lat, lon) => 
      set((state) => ({ 
        mapParams: { ...state.mapParams, lat, lon } 
      })),
    
    setRadius: (radius_km) => 
      set((state) => ({ 
        mapParams: { ...state.mapParams, radius_km } 
      })),
    
    setCategories: (categoryIds) => 
      set((state) => ({ 
        mapParams: { ...state.mapParams, category_id: categoryIds } 
      })),
    
    setLastFetchTime: (time) => 
      set({ lastFetchTime: time }),
    
    setIsUserEvents: (isUserEvents) => 
      set({ isUserEvents }),
    
    resetMapParams: () => 
      set({ mapParams: defaultParams }),
  },
}));

// Селекторы
export const useMapParams = () => useMapStore((state) => state.mapParams);
export const useLastFetchTime = () => useMapStore((state) => state.lastFetchTime);
export const useIsUserEvents = () => useMapStore((state) => state.isUserEvents);
export const useMapActions = () => useMapStore((state) => state.actions);