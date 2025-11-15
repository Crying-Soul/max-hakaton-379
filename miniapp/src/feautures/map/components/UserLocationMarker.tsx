import React from 'react';
import { Marker } from 'react-leaflet';
import L from 'leaflet';

const createUserIcon = (): L.DivIcon =>
  new L.DivIcon({
    html: `
      <div style="position:relative; width:44px; height:44px;">
        <div style="width:20px;height:20px;border-radius:50%;background:#3B82F6;border:3px solid white;box-shadow:0 2px 8px rgba(0,0,0,0.3);position:absolute;left:50%;top:50%;transform:translate(-50%,-50%);"></div>
        <div style="width:44px;height:44px;border-radius:50%;background:#3B82F6;opacity:0.2;position:absolute;left:50%;top:50%;transform:translate(-50%,-50%);animation:rl-pulse 2s infinite;"></div>
      </div>
      <style>
        @keyframes rl-pulse {
          0% { transform: translate(-50%,-50%) scale(1); opacity: 0.2; }
          50% { transform: translate(-50%,-50%) scale(1.2); opacity: 0.1; }
          100% { transform: translate(-50%,-50%) scale(1); opacity: 0.2; }
        }
      </style>
    `,
    className: 'user-location-marker',
    iconSize: [44, 44],
    iconAnchor: [22, 22],
  });

interface Props {
  location: [number, number] | null;
}

const UserLocationMarker: React.FC<Props> = React.memo(({ location }) => {
  if (!location) return null;
  return (
    <Marker position={location} icon={createUserIcon()} zIndexOffset={1000} />
  );
});

UserLocationMarker.displayName = 'UserLocationMarker';
export default UserLocationMarker;
