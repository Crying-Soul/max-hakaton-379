import React, { useCallback, useState } from 'react';
import { useMap } from 'react-leaflet';
import { useWebApp } from '../../../hooks/useWebApp';

interface Props {
  onLocateUser?: () => Promise<void>;
  isLoadingLocation?: boolean;
}

const MapControls: React.FC<Props> = React.memo(
  ({ onLocateUser, isLoadingLocation = false }) => {
    const map = useMap();
    const [locating, setLocating] = useState(false);
    const { impactOccurred } = useWebApp();

    const handleZoomIn = useCallback(() => {
      map.zoomIn();
      impactOccurred('light');
    }, [map, impactOccurred]);

    const handleZoomOut = useCallback(() => {
      map.zoomOut();
      impactOccurred('light');
    }, [map, impactOccurred]);

    const handleLocate = useCallback(async () => {
      if (!onLocateUser || isLoadingLocation || locating) return;

      impactOccurred('medium');
      setLocating(true);

      try {
        await onLocateUser();
      } finally {
        setLocating(false);
      }
    }, [onLocateUser, isLoadingLocation, locating, impactOccurred]);

    const canLocate =
      typeof navigator !== 'undefined' && 'geolocation' in navigator;

    return (
      <div className="leaflet-top leaflet-right pointer-events-none">
        <div className="leaflet-control leaflet-bar !border-none !bg-transparent space-y-2 mr-2 mt-2 pointer-events-auto">
          <div className="bg-white/90 backdrop-blur rounded-xl shadow-lg overflow-hidden flex flex-col">
            <button
              onClick={handleZoomIn}
              className="w-10 h-10 flex items-center justify-center text-xl text-gray-700 hover:bg-gray-50 active:bg-gray-100 transition-colors"
              aria-label="Zoom in"
            >
              +
            </button>
            <div className="border-t border-gray-200/50" />
            <button
              onClick={handleZoomOut}
              className="w-10 h-10 flex items-center justify-center text-xl text-gray-700 hover:bg-gray-50 active:bg-gray-100 transition-colors"
              aria-label="Zoom out"
            >
              −
            </button>
          </div>

          <button
            className={`w-10 h-10 flex items-center justify-center rounded-xl shadow-lg transition-all ${
              canLocate
                ? 'bg-blue-500 text-white hover:bg-blue-600 active:bg-blue-700'
                : 'bg-white/90 text-gray-400'
            } ${
              isLoadingLocation || locating
                ? 'opacity-50 cursor-not-allowed'
                : ''
            }`}
            onClick={handleLocate}
            disabled={!canLocate || locating || isLoadingLocation}
            aria-label="Locate me"
          >
            {locating || isLoadingLocation ? (
              <div className="w-4 h-4 border-2 border-current border-t-transparent rounded-full animate-spin" />
            ) : (
              <svg className="w-4 h-4" fill="currentColor" viewBox="0 0 20 20">
                <path d="M5.05 4.05a7 7 0 119.9 9.9L10 18.9l-4.95-4.95a7 7 0 010-9.9zM10 11a2 2 0 100-4 2 2 0 000 4z" />
              </svg>
            )}
          </button>
        </div>
      </div>
    );
  },
);

MapControls.displayName = 'MapControls';
export default MapControls;
