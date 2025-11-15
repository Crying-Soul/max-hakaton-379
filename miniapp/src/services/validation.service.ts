import type { ValidationResult } from '../common/types/max-bridge';

const WEB_APP_DATA_PREFIX = 'WebAppData';
const DEFAULT_MAX_AGE_MINUTES = 60;

export class ValidationService  {
  static async validateInitData(
    botToken: string,
    maxAgeMinutes: number = DEFAULT_MAX_AGE_MINUTES,
  ): Promise<ValidationResult> {
    try {
      if (typeof window === 'undefined') {
        return {
          isValid: false,
          error: 'Not running in browser environment',
        };
      }

      const webApp = window.WebApp;

      if (!webApp) {
        return { isValid: false, error: 'WebApp bridge not found' };
      }

      const raw = webApp.initData;
      const initDataUnsafe = webApp.initDataUnsafe;

      if (!raw) {
        return {
          isValid: false,
          error: 'initData string not provided by WebApp',
          initData: initDataUnsafe,
        };
      }

      let decodedRaw: string;
      try {
        decodedRaw = decodeURIComponent(raw);
      } catch (e) {
        return {
          isValid: false,
          error: 'Failed to URL decode initData',
          initData: initDataUnsafe,
        };
      }

      const params = new URLSearchParams(decodedRaw);
      const receivedHash = params.get('hash');

      if (!receivedHash) {
        return {
          isValid: false,
          error: 'hash parameter missing in initData',
          initData: initDataUnsafe,
        };
      }

      const entries = Array.from(params.entries())
        .filter(([k]) => k !== 'hash')
        .sort((a, b) => a[0].localeCompare(b[0]));

      const dataPairs = entries.map(([key, value]) => {
        if (key === 'user' || key === 'chat') {
          try {
            const parsed = JSON.parse(value);
            const normalized = JSON.stringify(parsed);
            return `${key}=${normalized}`;
          } catch (e) {
            // Если не валидный JSON, используем как есть
          }
        }
        return `${key}=${value}`;
      });

      const dataCheckString = dataPairs.join('\n');

      if (!botToken) {
        return { isValid: false, error: 'botToken is required' };
      }

      if (!crypto?.subtle) {
        return {
          isValid: false,
          error: 'Web Crypto API (crypto.subtle) not available',
        };
      }

      const enc = new TextEncoder();

      const secretKey = await crypto.subtle.importKey(
        'raw',
        enc.encode(WEB_APP_DATA_PREFIX),
        { name: 'HMAC', hash: 'SHA-256' },
        false,
        ['sign'],
      );

      const secretKeyBytes = await crypto.subtle.sign(
        'HMAC',
        secretKey,
        enc.encode(botToken),
      );

      const verificationKey = await crypto.subtle.importKey(
        'raw',
        secretKeyBytes,
        { name: 'HMAC', hash: 'SHA-256' },
        false,
        ['sign'],
      );

      const signature = await crypto.subtle.sign(
        'HMAC',
        verificationKey,
        enc.encode(dataCheckString),
      );

      const expectedHash = this.arrayBufferToHex(signature);
      const isValid = this.compareHexHashes(receivedHash, expectedHash);

      const authDateParam = params.get('auth_date');
      let isFresh = true;

      if (authDateParam) {
        const authTimestamp = Number(authDateParam);
        if (Number.isFinite(authTimestamp) && authTimestamp > 0) {
          const authDateMs = authTimestamp * 1000;
          const nowMs = Date.now();
          const ageMs = nowMs - authDateMs;
          const maxAgeMs = maxAgeMinutes * 60 * 1000;
          isFresh = ageMs <= maxAgeMs;
        } else {
          isFresh = false;
        }
      }

      return {
        isValid: isValid && isFresh,
        isFresh,
        initData: initDataUnsafe,
        error: isValid && isFresh ? undefined : 
               !isValid ? 'Data validation failed: signature mismatch' :
               !isFresh ? 'Data validation failed: data expired' : 
               'Data validation failed',
      };
    } catch (err) {
      const errorMsg = err instanceof Error ? err.message : String(err);
      return {
        isValid: false,
        error: `Validation error: ${errorMsg}`,
      };
    }
  }

  private static arrayBufferToHex(buffer: ArrayBuffer): string {
    const bytes = new Uint8Array(buffer);
    return Array.from(bytes)
      .map((byte) => byte.toString(16).padStart(2, '0'))
      .join('')
      .toLowerCase();
  }

  private static hexToBytes(hex: string): Uint8Array {
    const cleanHex = hex.toLowerCase().replace(/[^0-9a-f]/g, '');
    if (cleanHex.length % 2 !== 0) {
      throw new Error('Invalid hex string: odd length');
    }

    const bytes = new Uint8Array(cleanHex.length / 2);
    for (let i = 0; i < bytes.length; i++) {
      const byte = parseInt(cleanHex.substring(i * 2, i * 2 + 2), 16);
      if (isNaN(byte)) {
        throw new Error('Invalid hex string: contains non-hex characters');
      }
      bytes[i] = byte;
    }
    return bytes;
  }

  private static compareHexHashes(received: string, expected: string): boolean {
    if (!received || !expected) return false;

    const receivedClean = received.toLowerCase().trim();
    const expectedClean = expected.toLowerCase().trim();

    if (receivedClean.length !== expectedClean.length) return false;

    try {
      const receivedBytes = this.hexToBytes(receivedClean);
      const expectedBytes = this.hexToBytes(expectedClean);

      let diff = 0;
      for (let i = 0; i < receivedBytes.length; i++) {
        diff |= receivedBytes[i] ^ expectedBytes[i];
      }
      return diff === 0;
    } catch {
      return false;
    }
  }
}