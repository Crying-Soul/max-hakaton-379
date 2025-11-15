// src/features/map/components/MarkersLayer.tsx
import React, { useEffect, useMemo, useCallback } from 'react';
import { Marker, Popup, useMap } from 'react-leaflet';
import L from 'leaflet';
import type { LatLngExpression } from 'leaflet';
import { useMapEvents } from '../../../hooks/useMapEvents';
import {
  useMapParams,
  useIsUserEvents,
  useMapActions,
} from '../../../store/map.store';
import { ActivityAdapterService } from '../../../services/activity-adapter.service';
import { useUIActions } from '../../../store/ui.store';
import { useAuth } from '../../../hooks/useAuth';
import { Activity } from '../../../common/types/domain';

const iconCache = new Map<string, L.DivIcon>();

const createCustomIcon = (status: string): L.DivIcon => {
  const cacheKey = status;
  if (iconCache.has(cacheKey)) return iconCache.get(cacheKey)!;

  const color = status === 'open' ? '#10B981' : '#EF4444';

  const icon = new L.DivIcon({
    html: `
      <div style="position:relative; width:32px; height:32px;">
        <div style="width:32px;height:32px;border-radius:50%;background:#fff;box-shadow:0 2px 8px rgba(0,0,0,0.15);border:2px solid #fff;display:flex;align-items:center;justify-content:center;position:absolute;left:50%;top:50%;transform:translate(-50%,-50%);">
          <div style="width:18px;height:18px;border-radius:50%;background:${color};"></div>
        </div>
      </div>
    `,
    className: 'custom-marker',
    iconSize: [32, 32],
    iconAnchor: [16, 32],
  });

  iconCache.set(cacheKey, icon);
  return icon;
};

const MarkersLayer: React.FC = React.memo(() => {
  const mapParams = useMapParams();
  const isUserEvents = useIsUserEvents();
  const { setSelected, setSessionModalOpen } = useUIActions();
  const { setLastFetchTime } = useMapActions();
  const { isAuthenticated } = useAuth();
  const {
    events,
    loading,
    error,
    fetchMapEvents,
    fetchUserMapEvents,
    lastUpdated,
  } = useMapEvents();
  
  const map = useMap();

  // Загрузка данных
  useEffect(() => {
    const fetchEvents = async () => {
      try {
        if (isUserEvents && isAuthenticated) {
          await fetchUserMapEvents(mapParams);
        } else {
          await fetchMapEvents(mapParams);
        }
        setLastFetchTime(new Date());
      } catch (err) {
        console.error('Error fetching events:', err);
      }
    };

    fetchEvents();
  }, [mapParams, isUserEvents, isAuthenticated]);

  // Преобразование событий в активности
  const activities = useMemo(() => {
    return ActivityAdapterService.mapEventsToActivities(events);
  }, [events]);

  // Создание маркеров с иконками
  const markers = useMemo(
    () =>
      activities.map((a) => ({
        ...a,
        position: [a.lat, a.lon] as LatLngExpression,
        icon: createCustomIcon(a.status),
      })),
    [activities],
  );

  const handleMarkerClick = useCallback(
    (activity: Activity) => {
      setSelected(activity);
      map.flyTo([activity.lat, activity.lon], 16, { duration: 1 });
    },
    [setSelected, map],
  );

  const handleOpenActivity = useCallback(
    (activity: Activity) => {
      setSelected(activity);
      setSessionModalOpen(true);
    },
    [setSelected, setSessionModalOpen],
  );

  if (loading) {
    return (
      <div className="absolute top-4 left-1/2 transform -translate-x-1/2 bg-white px-4 py-2 rounded-lg shadow-lg z-[1000]">
        <div className="flex items-center space-x-2">
          <div className="animate-spin rounded-full h-4 w-4 border-b-2 border-blue-600"></div>
          <span className="text-sm text-gray-600">Загрузка мероприятий...</span>
        </div>
      </div>
    );
  }

  if (error) {
    return (
      <div className="absolute top-4 left-1/2 transform -translate-x-1/2 bg-red-50 px-4 py-2 rounded-lg shadow-lg z-[1000]">
        <span className="text-sm text-red-600">Ошибка: {error}</span>
      </div>
    );
  }

  return (
    <>
      {markers.map((m) => (
        <Marker
          key={m.id}
          position={m.position}
          icon={m.icon}
          eventHandlers={{
            click: () => handleMarkerClick(m as Activity),
          }}
          riseOnHover
        >
          <Popup
            closeButton={false}
            className="custom-popup"
            autoPan
            autoPanPadding={[20, 20]}
          >
            <div className="max-w-[260px] p-3">
              <div className="flex items-start justify-between mb-2">
                <h3 className="font-bold text-gray-900 text-sm pr-2">
                  {m.title}
                </h3>
                <span
                  className={`px-2 py-0.5 rounded-full text-[10px] font-medium ${
                    m.status === 'open'
                      ? 'bg-green-100 text-green-700'
                      : 'bg-red-100 text-red-700'
                  }`}
                >
                  {m.status === 'open' ? 'Открыто' : 'Закрыто'}
                </span>
              </div>

              <p className="text-gray-600 text-xs mb-3 line-clamp-2">
                {m.description}
              </p>

              <div className="flex items-center text-gray-500 text-[11px] mb-3">
                <svg
                  className="w-3.5 h-3.5 mr-1"
                  fill="currentColor"
                  viewBox="0 0 20 20"
                  aria-hidden
                >
                  <path d="M5.05 4.05a7 7 0 119.9 9.9L10 18.9l-4.95-4.95a7 7 0 010-9.9zM10 11a2 2 0 100-4 2 2 0 000 4z" />
                </svg>
                {m.distanceText ?? '—'} от вас
              </div>

              <button
                onClick={() => handleOpenActivity(m as Activity)}
                className="w-full py-2 rounded-xl bg-blue-600 text-white text-xs font-medium hover:bg-blue-700"
              >
                Открыть
              </button>
            </div>
          </Popup>
        </Marker>
      ))}
      
      {lastUpdated && (
        <div className="absolute bottom-4 left-1/2 transform -translate-x-1/2 bg-black/50 text-white px-3 py-1 rounded-full text-xs z-[1000]">
          Обновлено: {lastUpdated.toLocaleTimeString()}
        </div>
      )}
    </>
  );
});

MarkersLayer.displayName = 'MarkersLayer';
export default MarkersLayer;