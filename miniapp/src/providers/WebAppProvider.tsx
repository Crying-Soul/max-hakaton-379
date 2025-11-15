// src/providers/WebAppProvider.tsx
import React, { useEffect, useRef, useState } from 'react';
import { useWebApp, useBackButton } from '../hooks/useWebApp';
import { useAuth } from '../hooks/useAuth';

interface WebAppProviderProps {
  children: React.ReactNode;
}

export const WebAppProvider: React.FC<WebAppProviderProps> = ({ children }) => {
  const { webApp, isReady, isLoading, impactOccurred } = useWebApp();
  const { login, loading: authLoading, error: authError } = useAuth();
  const [error, setError] = useState<string | null>(null);
  const [authAttempted, setAuthAttempted] = useState(false);
  const mounted = useRef(true);

  useEffect(() => {
    return () => {
      mounted.current = false;
    };
  }, []);

  useBackButton(() => {
    console.log('Back button pressed');
    impactOccurred('light');
    webApp?.close();
  });

  useEffect(() => {
    if (!isReady || !webApp || authAttempted) return;

    (async () => {
      try {
        console.log('🚀 Initializing MAX WebApp...', {
          platform: webApp.platform,
          version: webApp.version,
          hasInitData: !!webApp.initData,
        });

        webApp.enableClosingConfirmation();

        const initData = webApp.initData;
        if (initData) {
          console.log('🔐 Attempting authentication with MAX WebApp initData');
          const success = await login(initData);
          if (!mounted.current) return;

          if (success) {
            console.log('✅ Authentication successful');
            impactOccurred('soft');
          } else {
            console.warn('❌ Authentication failed: login returned falsy');
            setError('Authentication failed. Running in limited mode.');
          }
        } else {
          console.warn('⚠️ No initData available - running in limited mode');
          setError('No initData - running in limited mode');
        }
      } catch (e: any) {
        console.error('❌ Authentication error:', e);
        setError(e?.message || 'Auth error');
      } finally {
        if (mounted.current) setAuthAttempted(true);
      }
    })();
  }, [isReady, webApp, login, authAttempted, impactOccurred]);

  if (isLoading || (authAttempted && authLoading)) {
    return (
      <div className="flex items-center justify-center min-h-screen bg-gray-900">
        <div className="text-white text-center">
          <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-white mx-auto mb-4"></div>
          <p>{authLoading && authAttempted ? 'Authenticating...' : 'Initializing MAX WebApp...'}</p>
        </div>
      </div>
    );
  }

  if (authError) {
    return (
      <div className="flex items-center justify-center min-h-screen bg-gray-900">
        <div className="text-white text-center p-6 max-w-md">
          <div className="text-yellow-500 text-4xl mb-4">⚠️</div>
          <h2 className="text-xl font-bold mb-2">Authentication Error</h2>
          <p className="text-gray-300 mb-4">{authError}</p>
          <p className="text-sm text-gray-400">Some features like personal events may not be available.</p>
        </div>
      </div>
    );
  }

  return (
    <>
      {error && (
        <div className="fixed top-4 left-1/2 transform -translate-x-1/2 z-50">
          <div className="bg-yellow-500 text-white px-4 py-2 rounded-lg shadow-lg max-w-md">
            <div className="flex items-center space-x-2">
              <span>⚠️</span>
              <span className="text-sm">{error}</span>
            </div>
          </div>
        </div>
      )}
      {children}
    </>
  );
};