// src/common/types/domain.ts
export interface Activity {
  id: string;
  title: string;
  description: string;
  date: string;
  durationHours: number;
  location: string;
  lat: number;
  lon: number;
  status: 'open' | 'closed';
  distanceText: string;
  distanceKm: number;
  categoryId: number;
  categoryName: string;
  organizerId: number;
  contacts: string;
  maxVolunteers: number;
  currentVolunteers: number;
  deeplink: string
  slotsLeft: number;
  createdAt: string;
  updatedAt: string;
}

export interface GeolocationState {
  location: [number, number] | null;
  error: string | null;
  loading: boolean;
}

export interface User {
  id: number;
  first_name: string;
  last_name?: string;
  username?: string;
  language_code?: string;
}