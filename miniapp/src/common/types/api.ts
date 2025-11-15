// src/common/types/api.ts
export interface MapEvent {
  id: number;
  title: string;
  description: string;
  date: string;
  durationHours: number;
  location: string;
  locationLat: number;
  locationLon: number;
  categoryId: number;
  categoryName: string;
  organizerId: number;
  contacts: string;
  maxVolunteers: number;
  currentVolunteers: number;
  slotsLeft: number;
  distanceKm: number;
  status: 'open' | 'closed' | 'cancelled' | 'completed';
  createdAt: string;
  updatedAt: string;
}

export interface MapEventsResponse {
  data: MapEvent[];
  meta: {
    categories: string[] | null;
    count: number;
    lat: number;
    limit: number;
    lon: number;
    offset: number;
    radiusKm: number;
  };
}

export interface AuthSessionRequest {
  initData: string;
}

export interface AuthSessionResponse {
  token: string;
  expiresIn: number;
  user: User;
}

export interface MapEventsParams {
  lat: number;
  lon: number;
  radius_km: number;
  limit?: number;
  offset?: number;
  categories?:  number[];
  category_id?: number[];
}

export interface ApiError {
  code: string;
  message: string;
}

export interface User {
  id: number;
  first_name: string;
  last_name?: string;
  username?: string;
  language_code?: string;
}
