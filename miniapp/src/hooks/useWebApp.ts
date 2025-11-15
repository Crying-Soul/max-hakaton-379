import { useEffect, useState, useCallback } from 'react';
import { WebApp, WebAppEventMap } from '../common/types/max-bridge';

export const useWebApp = () => {
  const [webApp, setWebApp] = useState<WebApp | null>(null);
  const [isReady, setIsReady] = useState(false);
  const [isLoading, setIsLoading] = useState(true);

  useEffect(() => {
    const initWebApp = () => {
      const webAppInstance = window.WebApp;

      if (webAppInstance) {
        setWebApp(webAppInstance);
        setIsReady(true);

        // Инициализируем приложение
        webAppInstance.ready();

        console.log('MAX WebApp initialized:', {
          platform: webAppInstance.platform,
          version: webAppInstance.version,
          initData: webAppInstance.initDataUnsafe,
        });
      } else {
        console.warn('MAX WebApp bridge not found');
      }

      setIsLoading(false);
    };

    // Проверяем сразу и через таймаут
    initWebApp();
    const timeoutId = setTimeout(initWebApp, 500);

    return () => clearTimeout(timeoutId);
  }, []);

  // Основные методы
  const ready = useCallback(() => webApp?.ready(), [webApp]);
  const close = useCallback(() => webApp?.close(), [webApp]);

  const openLink = useCallback(
    (url: string) => {
      webApp?.openLink(url);
    },
    [webApp],
  );

  const openMaxLink = useCallback(
    (url: string) => {
      webApp?.openMaxLink(url);
    },
    [webApp],
  );

  // BackButton методы
  const showBackButton = useCallback(() => webApp?.BackButton.show(), [webApp]);
  const hideBackButton = useCallback(() => webApp?.BackButton.hide(), [webApp]);

  const onBackButtonClick = useCallback(
    (callback: () => void) => {
      if (!webApp?.BackButton) return () => {};

      webApp.BackButton.onClick(callback);
      return () => webApp.BackButton.offClick(callback);
    },
    [webApp],
  );

  // HapticFeedback методы
  const impactOccurred = useCallback(
    (style: 'light' | 'medium' | 'heavy' | 'rigid' | 'soft') => {
      webApp?.HapticFeedback.impactOccurred(style);
    },
    [webApp],
  );

  const notificationOccurred = useCallback(
    (type: 'error' | 'success' | 'warning') => {
      webApp?.HapticFeedback.notificationOccurred(type);
    },
    [webApp],
  );

  const selectionChanged = useCallback(() => {
    webApp?.HapticFeedback.selectionChanged();
  }, [webApp]);

  // ScreenCapture методы
  const enableScreenCapture = useCallback(async () => {
    return webApp?.ScreenCapture.enableScreenCapture();
  }, [webApp]);

  const disableScreenCapture = useCallback(async () => {
    return webApp?.ScreenCapture.disableScreenCapture();
  }, [webApp]);

  // Storage методы
  const setSecureStorageValue = useCallback(
    async (key: string, value: string) => {
      return webApp?.SecureStorage.setItem(key, value);
    },
    [webApp],
  );

  const getSecureStorageValue = useCallback(
    async (key: string) => {
      return webApp?.SecureStorage.getItem(key);
    },
    [webApp],
  );

  const setDeviceStorageValue = useCallback(
    async (key: string, value: string) => {
      return webApp?.DeviceStorage.setItem(key, value);
    },
    [webApp],
  );

  const getDeviceStorageValue = useCallback(
    async (key: string) => {
      return webApp?.DeviceStorage.getItem(key);
    },
    [webApp],
  );

  // Biometric методы
  const initBiometric = useCallback(async () => {
    return webApp?.BiometricManager.init();
  }, [webApp]);

  const requestBiometricAccess = useCallback(async () => {
    return webApp?.BiometricManager.requestAccess();
  }, [webApp]);

  // Share методы
  const shareContent = useCallback(
    async (content: { text: string; link: string }) => {
      return webApp?.shareContent(content);
    },
    [webApp],
  );

  const shareMaxContent = useCallback(
    async (content: { text: string; link: string }) => {
      return webApp?.shareMaxContent(content);
    },
    [webApp],
  );

  // Event методы
  const onEvent = useCallback(
    <K extends keyof WebAppEventMap>(
      eventName: K,
      callback: (data: WebAppEventMap[K]) => void,
    ) => {
      if (!webApp) return () => {};

      return webApp.onEvent(eventName, callback);
    },
    [webApp],
  );

  return {
    // Состояние
    webApp,
    isReady,
    isLoading,

    // Основные методы
    ready,
    close,
    openLink,
    openMaxLink,

    // BackButton
    showBackButton,
    hideBackButton,
    onBackButtonClick,

    // HapticFeedback
    impactOccurred,
    notificationOccurred,
    selectionChanged,

    // ScreenCapture
    enableScreenCapture,
    disableScreenCapture,

    // Storage
    setSecureStorageValue,
    getSecureStorageValue,
    setDeviceStorageValue,
    getDeviceStorageValue,

    // Biometric
    initBiometric,
    requestBiometricAccess,

    // Share
    shareContent,
    shareMaxContent,

    // Events
    onEvent,

    // Данные
    initData: webApp?.initDataUnsafe,
    initDataRaw: webApp?.initData,
    platform: webApp?.platform,
    version: webApp?.version,

    // Менеджеры
    BackButton: webApp?.BackButton,
    ScreenCapture: webApp?.ScreenCapture,
    HapticFeedback: webApp?.HapticFeedback,
    SecureStorage: webApp?.SecureStorage,
    DeviceStorage: webApp?.DeviceStorage,
    BiometricManager: webApp?.BiometricManager,
  };
};

// Специализированные хуки
export const useBackButton = (callback: () => void) => {
  const { onBackButtonClick } = useWebApp();

  useEffect(() => {
    const cleanup = onBackButtonClick(callback);
    return cleanup;
  }, [onBackButtonClick, callback]);
};

export const useWebAppEvent = <K extends keyof WebAppEventMap>(
  eventName: K,
  callback: (data: WebAppEventMap[K]) => void,
) => {
  const { onEvent } = useWebApp();

  useEffect(() => {
    const cleanup = onEvent(eventName, callback);
    return cleanup;
  }, [onEvent, eventName, callback]);
};

// Хук для ScreenCapture
export const useScreenCapture = () => {
  const { webApp } = useWebApp();
  const [isScreenCaptureEnabled, setIsScreenCaptureEnabled] = useState(false);

  useEffect(() => {
    if (webApp?.ScreenCapture) {
      setIsScreenCaptureEnabled(webApp.ScreenCapture.isScreenCaptureEnabled);
    }
  }, [webApp]);

  const enable = useCallback(async () => {
    if (!webApp?.ScreenCapture) return;
    const result = await webApp.ScreenCapture.enableScreenCapture();
    setIsScreenCaptureEnabled(result.isScreenCaptureEnabled);
    return result;
  }, [webApp]);

  const disable = useCallback(async () => {
    if (!webApp?.ScreenCapture) return;
    const result = await webApp.ScreenCapture.disableScreenCapture();
    setIsScreenCaptureEnabled(result.isScreenCaptureEnabled);
    return result;
  }, [webApp]);

  return {
    isScreenCaptureEnabled,
    enableScreenCapture: enable,
    disableScreenCapture: disable,
  };
};

// Хук для HapticFeedback
export const useHapticFeedback = () => {
  const { impactOccurred, notificationOccurred, selectionChanged } =
    useWebApp();

  return {
    impactOccurred,
    notificationOccurred,
    selectionChanged,
  };
};
