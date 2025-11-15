import { useState, useCallback } from 'react';
import { GeolocationState } from '../common/types/domain';

export const useGeolocation = () => {
  const [state, setState] = useState<GeolocationState>({
    location: null,
    error: null,
    loading: false,
  });

  const requestLocation = useCallback(async (): Promise<
    [number, number] | null
  > => {
    if (!navigator.geolocation) {
      setState((prev) => ({
        ...prev,
        error: 'Геолокация не поддерживается вашим браузером',
        loading: false,
      }));
      return null;
    }

    setState((prev) => ({ ...prev, loading: true, error: null }));

    return new Promise((resolve) => {
      const options: PositionOptions = {
        enableHighAccuracy: true,
        timeout: 30000,
        maximumAge: 600000,
      };

      const onSuccess = (position: GeolocationPosition) => {
        const { latitude, longitude } = position.coords;
        const location: [number, number] = [latitude, longitude];

        setState({ location, error: null, loading: false });
        resolve(location);
      };

      const onError = (error: GeolocationPositionError) => {
        let errorMessage = 'Не удалось определить местоположение';

        switch (error.code) {
          case error.PERMISSION_DENIED:
            errorMessage = 'Доступ к геолокации запрещен.';
            break;
          case error.POSITION_UNAVAILABLE:
            errorMessage = 'Информация о местоположении недоступна';
            break;
          case error.TIMEOUT:
            errorMessage = 'Время ожидания определения местоположения истекло';
            break;
        }

        setState((prev) => ({ ...prev, error: errorMessage, loading: false }));
        resolve(null);
      };

      navigator.geolocation.getCurrentPosition(onSuccess, onError, options);
    });
  }, []);

  const clearLocation = useCallback(() => {
    setState({ location: null, error: null, loading: false });
  }, []);

  const clearError = useCallback(() => {
    setState((prev) => ({ ...prev, error: null }));
  }, []);

  return {
    location: state.location,
    error: state.error,
    loading: state.loading,
    requestLocation,
    clearLocation,
    clearError,
  };
};
