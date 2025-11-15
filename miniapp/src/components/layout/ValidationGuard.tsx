// src/components/ValidationGuard.tsx
import React from 'react';
import { useDataValidation } from '../../hooks/useDataValidation';
import { useWebApp } from '../../hooks/useWebApp';

interface ValidationGuardProps {
  botToken: string;
  children: React.ReactNode;
  fallback?: React.ReactNode;
  onValidationSuccess?: (initData: any) => void;
  onValidationError?: (error: string) => void;
  maxAgeMinutes?: number;
}

export const ValidationGuard: React.FC<ValidationGuardProps> = ({
  botToken,
  children,
  fallback,
  onValidationSuccess,
  onValidationError,
  maxAgeMinutes = 60,
}) => {
  const { isReady, isLoading } = useWebApp();
  const { isValid, error, isValidating, initData, isFresh } = useDataValidation(
    botToken,
    maxAgeMinutes,
  );

  React.useEffect(() => {
    if (!isReady) return;

    if (!isValidating) {
      if (isValid) {
        if (isFresh === false) {
          onValidationError?.('Validation expired');
        } else {
          onValidationSuccess?.(initData);
        }
      } else {
        onValidationError?.(error || 'Validation failed');
      }
    }
  }, [isValid, isFresh, error, isValidating, isReady, initData, onValidationSuccess, onValidationError]);

  if (isLoading) {
    return (
      <div className="flex items-center justify-center min-h-screen bg-gray-900">
        <div className="text-white text-center">
          <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-white mx-auto mb-4"></div>
          <p>Initializing MAX Bridge...</p>
        </div>
      </div>
    );
  }

  if (!isReady) {
    return (
      <div className="flex items-center justify-center min-h-screen bg-gray-900">
        <div className="text-white text-center p-6 max-w-md">
          <div className="text-yellow-500 text-4xl mb-4">⚠️</div>
          <h2 className="text-xl font-bold mb-2">MAX Bridge Not Found</h2>
          <p className="text-gray-300">This app is designed to work within MAX messenger. Please open it through MAX app.</p>
        </div>
      </div>
    );
  }

  if (isValidating) {
    return (
      <div className="flex items-center justify-center min-h-screen bg-gray-900">
        <div className="text-white text-center">
          <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-white mx-auto mb-4"></div>
          <p>Validating security data...</p>
        </div>
      </div>
    );
  }

  if (!isValid || isFresh === false) {
    if (fallback) return <>{fallback}</>;

    return (
      <div className="flex items-center justify-center min-h-screen bg-gray-900">
        <div className="text-white text-center p-6 max-w-md">
          <div className={`${isFresh === false ? 'text-yellow-400' : 'text-red-500'} text-4xl mb-4`}>
            {isFresh === false ? '⌛' : '⚠️'}
          </div>
          <h2 className="text-xl font-bold mb-2">
            {isFresh === false ? 'Session Expired' : 'Security Error'}
          </h2>
          <p className="text-gray-300 mb-4">
            {isFresh === false 
              ? 'Authentication data is expired. Please re-open the mini-app from MAX to refresh session.'
              : 'Data validation failed. The application cannot be trusted.'}
          </p>
          <p className="text-sm text-gray-400">{error}</p>
        </div>
      </div>
    );
  }

  return <>{children}</>;
};