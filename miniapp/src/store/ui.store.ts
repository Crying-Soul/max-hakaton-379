import { create } from 'zustand';
import { Activity } from '../common/types/domain';
import { MapEventsResponse } from '../common/types/api';

interface UIState {
  selected: Activity | null;
  isRegistered: boolean;
  isSessionModalOpen: boolean;
  lastFetchMeta: MapEventsResponse['meta'] | null;
  actions: {
    setSelected: (activity: Activity | null) => void;
    setRegistered: (value: boolean) => void;
    setSessionModalOpen: (value: boolean) => void;
    setLastFetchMeta: (meta: MapEventsResponse['meta'] | null) => void;
  };
}

export const useUIStore = create<UIState>((set) => ({
  selected: null,
  isRegistered: false,
  isSessionModalOpen: false,
  lastFetchMeta: null,
  actions: {
    setSelected: (selected) => set({ selected }),
    setRegistered: (isRegistered) => set({ isRegistered }),
    setSessionModalOpen: (isSessionModalOpen) => set({ isSessionModalOpen }),
    setLastFetchMeta: (lastFetchMeta) => set({ lastFetchMeta }),
  },
}));

// Селекторы
export const useSelectedActivity = () => useUIStore((state) => state.selected);
export const useIsSessionModalOpen = () => useUIStore((state) => state.isSessionModalOpen);
export const useUIActions = () => useUIStore((state) => state.actions);
export const useIsRegistered = () => useUIStore((state) => state.isRegistered);
export const useLastFetchMeta = () => useUIStore((state) => state.lastFetchMeta);