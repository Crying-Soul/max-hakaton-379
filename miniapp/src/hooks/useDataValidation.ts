// src/hooks/useDataValidation.ts
import { useEffect, useState } from 'react';
import type { ValidationResult } from '../common/types/max-bridge';
import { useWebApp } from './useWebApp';
import { ValidationService } from '../services/validation.service';

export const useDataValidation = (botToken: string, maxAgeMinutes = 60) => {
  const { webApp, initDataRaw, initData, isReady } = useWebApp();
  const [validation, setValidation] = useState<ValidationResult>({
    isValid: false,
    error: 'Waiting for WebApp bridge',
  });
  const [isValidating, setIsValidating] = useState(false);

  useEffect(() => {
    let cancelled = false;

    const runValidation = async () => {
      if (!isReady || !webApp) {
        setValidation({
          isValid: false,
          error: 'WebApp bridge not ready',
        });
        return;
      }

      if (!initDataRaw) {
        setValidation({
          isValid: false,
          error: 'initData not provided by WebApp',
          initData: initData ?? undefined,
        });
        return;
      }

      if (!botToken) {
        setValidation({ isValid: false, error: 'Bot token is required' });
        return;
      }

      setIsValidating(true);
      try {
        const result = await ValidationService.validateInitData(botToken, maxAgeMinutes);
        if (cancelled) return;
        setValidation(result);
      } catch (err) {
        if (cancelled) return;
        setValidation({
          isValid: false,
          error: err instanceof Error ? `Validation failed: ${err.message}` : 'Validation failed',
        });
      } finally {
        if (!cancelled) setIsValidating(false);
      }
    };

    runValidation();

    return () => {
      cancelled = true;
    };
  }, [botToken, isReady, webApp, initDataRaw, maxAgeMinutes]);

  return {
    ...validation,
    isValidating,
  };
};