export const APP_CONFIG = {
  DEFAULT_MAP_CENTER: [59.935, 30.325] as [number, number],
  TILE_URL: 'https://{s}.basemaps.cartocdn.com/rastertiles/voyager/{z}/{x}/{y}{r}.png',
  DEFAULT_MAX_AGE_MINUTES: 60,
  WEB_APP_DATA_PREFIX: 'WebAppData',
} as const;