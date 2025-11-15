// MapCanvas.tsx
import React, { useEffect, useState, useCallback } from "react";
import { MapContainer, TileLayer } from "react-leaflet";
import type { Map as LeafletMap } from "leaflet";
import "leaflet/dist/leaflet.css";

import MarkersLayer from "./MarkersLayer";
import UserLocationMarker from "./UserLocationMarker";
import MapControls from "./MapControls";
import { useGeolocation } from "../../../hooks/useGeo";

const DEFAULT_CENTER: [number, number] = [59.935, 30.325];
const TILE_URL = "https://{s}.basemaps.cartocdn.com/rastertiles/voyager/{z}/{x}/{y}{r}.png";

const MapCanvas: React.FC = () => {
  const [map, setMap] = useState<LeafletMap | null>(null);
  const [showLocationBanner, setShowLocationBanner] = useState(true);

  const { location: userLocation, error: locationError, loading: isLoadingLocation, requestLocation } =
    useGeolocation();

  useEffect(() => {
    void requestLocation();
  }, [requestLocation]);

  useEffect(() => {
    if (userLocation && map) {
      map.flyTo(userLocation, 16, { duration: 0.6 });
    }
  }, [userLocation, map]);

  useEffect(() => {
    if (userLocation && !isLoadingLocation) {
      const t = setTimeout(() => setShowLocationBanner(false), 3000);
      return () => clearTimeout(t);
    }
  }, [userLocation, isLoadingLocation]);

  const handleLocateUser = useCallback(async () => {
    if (isLoadingLocation) return;
    setShowLocationBanner(true);
    const coords = await requestLocation();
    if (coords && map) map.flyTo(coords, 16, { duration: 0.6 });
  }, [requestLocation, isLoadingLocation, map]);

  const center = userLocation ?? DEFAULT_CENTER;

  return (
    <div className="h-full w-full relative">
      {!map && (
        <div className="absolute inset-0 bg-gray-100 z-10 flex items-center justify-center">
          <div className="text-center">
            <div className="w-10 h-10 border-3 border-blue-500 border-t-transparent rounded-full animate-spin mx-auto mb-3" />
            <p className="text-gray-600 text-sm">Загрузка карты...</p>
          </div>
        </div>
      )}

      <MapContainer
        center={center}
        zoom={13}
        ref={(m: LeafletMap | null) => setMap(m)}
        style={{ height: "100%", width: "100%" }}
        zoomControl={false}
        preferCanvas={true}
        className="z-0"
      >
        <TileLayer url={TILE_URL} crossOrigin="anonymous" detectRetina maxZoom={19} minZoom={1} />
        <MarkersLayer />
        <UserLocationMarker location={userLocation} />
        <MapControls onLocateUser={handleLocateUser} isLoadingLocation={isLoadingLocation} />
      </MapContainer>

      {showLocationBanner && (isLoadingLocation || locationError) && (
        <div className="absolute top-3 left-3 right-3 z-40">
          <div className="bg-white/95 backdrop-blur rounded-xl px-3 py-2 shadow-lg">
            {isLoadingLocation && (
              <div className="flex items-center">
                <div className="w-3 h-3 border-2 border-blue-500 border-t-transparent rounded-full animate-spin mr-2" />
                <span className="text-blue-500 text-xs">Установка соединения...</span>
              </div>
            )}
            {locationError && (
              <div className="flex items-center">
                <svg className="w-3 h-3 text-red-500 mr-2" fill="currentColor" viewBox="0 0 20 20">
                  <path fillRule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zM8.707 7.293a1 1 0 00-1.414 1.414L8.586 10l-1.293 1.293a1 1 0 101.414 1.414L10 11.414l1.293 1.293a1 1 0 001.414-1.414L11.414 10l1.293-1.293a1 1 0 00-1.414-1.414L10 8.586 8.707 7.293z" clipRule="evenodd" />
                </svg>
                <span className="text-red-500 text-xs">{locationError}</span>
              </div>
            )}
          </div>
        </div>
      )}
    </div>
  );
};

export default React.memo(MapCanvas);
