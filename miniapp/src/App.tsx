// App.tsx
import React, { useEffect, useState } from 'react';
import MapCanvas from './feautures/map/components/MapCanvas';
import MobileBottomSheet from './components/layout/MobileBottomSheet';
import { useUIStore } from './store/ui.store';
import { ValidationGuard } from './components/layout/ValidationGuard';
import { WebAppProvider } from './providers/WebAppProvider';
import SessionGateModal from './components/auth/SessionGateModal';
import { AuthService } from './services/auth.service';

const AppContent: React.FC = () => {
  const selected = useUIStore((s) => s.selected);
  const [isAuthenticated, setIsAuthenticated] = useState(false);

  useEffect(() => {
    // Следим за изменениями аутентификации
    const unsubscribe = AuthService.onAuthChange(setIsAuthenticated);
    return unsubscribe;
  }, []);

  const handleProceed = () => {
    if (selected) {
      const webApp = window.WebApp;

      if (webApp && selected.deeplink) {
        webApp.openMaxLink(selected.deeplink);
      } else {
        alert(`Переходим в активити ${selected.title}`);
      }
    }
  };

  if (!isAuthenticated) {
    return (
      <div className="flex items-center justify-center min-h-screen bg-gray-900">
        <div className="text-white text-center p-6 max-w-md">
          <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-white mx-auto mb-4"></div>
          <p className="text-gray-300">Инициализация аутентификации...</p>
        </div>
      </div>
    );
  }

  return (
    <div className="h-screen bg-gray-900 relative">
      <MapCanvas />
      <MobileBottomSheet />
      <SessionGateModal onProceed={handleProceed} />
    </div>
  );
};

const WebAppContent: React.FC<{ botToken: string }> = ({ botToken }) => {
  // Инициализация аутентификации при загрузке
  useEffect(() => {
    AuthService.initializeAuth();
  }, []);

  return (
    <ValidationGuard botToken={botToken}>
      <AppContent />
    </ValidationGuard>
  );
};

const App: React.FC = () => {
  const BOT_TOKEN = import.meta.env.VITE_BOT_TOKEN;

  if (!BOT_TOKEN) {
    return (
      <div className="flex items-center justify-center min-h-screen bg-gray-900">
        <div className="text-white text-center p-6 max-w-md">
          <div className="text-red-500 text-4xl mb-4">⚠️</div>
          <h2 className="text-xl font-bold mb-2">Bot token is undefined</h2>
          <p className="text-gray-300">
            Please set VITE_BOT_TOKEN in your environment variables
          </p>
        </div>
      </div>
    );
  }

  return (
    <WebAppProvider>
      <WebAppContent botToken={BOT_TOKEN} />
    </WebAppProvider>
  );
};

export default App;