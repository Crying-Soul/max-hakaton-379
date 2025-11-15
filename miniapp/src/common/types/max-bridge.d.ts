// types/max-bridge.d.ts

// Основные типы данных
export interface WebAppUser {
  id: number;
  first_name: string;
  last_name?: string;
  username?: string;
  language_code?: string;
  photo_url?: string;
}

export interface WebAppChat {
  id: number;
  type: 'group' | 'supergroup' | 'channel' | 'bot';
  title?: string;
  username?: string;
  photo_url?: string;
}

export interface WebAppInitData {
  query_id?: string;
  auth_date: number;
  hash: string;
  start_param?: string;
  user?: WebAppUser;
  receiver?: WebAppUser;
  chat?: WebAppChat;
  chat_type?: 'sender' | 'private' | 'group' | 'supergroup' | 'channel';
  chat_instance?: string;
}

export interface ValidationResult {
  isValid: boolean;
  error?: string;
  initData?: WebAppInitData;
  isFresh?: boolean;
}

// BackButton интерфейс
export interface BackButton {
  readonly isVisible: boolean;
  show(): void;
  hide(): void;
  onClick(callback: () => void): void;
  offClick(callback: () => void): void;
}

// HapticFeedback интерфейс
export type ImpactStyle = 'light' | 'medium' | 'heavy' | 'rigid' | 'soft';
export type NotificationType = 'error' | 'success' | 'warning';

export interface HapticFeedback {
  impactOccurred(style: ImpactStyle): void;
  notificationOccurred(type: NotificationType): void;
  selectionChanged(): void;
}

// ScreenCapture интерфейс
export interface ScreenCapture {
  readonly isScreenCaptureEnabled: boolean;
  enableScreenCapture(): Promise<{ isScreenCaptureEnabled: boolean }>;
  disableScreenCapture(): Promise<{ isScreenCaptureEnabled: boolean }>;
}

// Storage интерфейсы
export interface Storage {
  setItem(key: string, value: string): Promise<void>;
  getItem(key: string): Promise<string | null>;
  removeItem(key: string): Promise<void>;
  clear(): Promise<void>;
}

// BiometricManager интерфейс
export type BiometricType = 'fingerprint' | 'face' | 'faceid' | 'unknown';

export interface BiometricManager {
  readonly isInited: boolean;
  readonly isBiometricAvailable: boolean;
  readonly isAccessRequested: boolean;
  readonly isAccessGranted: boolean;
  readonly isBiometricTokenSaved: boolean;
  readonly biometricType: BiometricType[];
  readonly deviceId: string | null;

  init(): Promise<void>;
  requestAccess(): Promise<void>;
  authenticate(): Promise<void>;
  updateBiometricToken(token: string): Promise<void>;
  openSettings(): void;
}

// Share интерфейсы
export interface ShareContent {
  text: string;
  link: string;
}

export interface ShareResult {
  status: 'shared' | 'cancelled';
}

// Code Reader интерфейс
export interface CodeReaderResult {
  value: string;
}

// Download File интерфейс
export interface DownloadFileResult {
  status: string;
}

// Event types
export interface WebAppEventMap {
  WebAppReady: never;
  WebAppClose: never;
  WebAppSetupBackButton: { isVisible: boolean };
  WebAppRequestPhone: never;
  WebAppSetupClosingBehavior: { needConfirmation: boolean };
  WebAppBackButtonPressed: never;
  WebAppOpenLink: { url: string };
  WebAppOpenMaxLink: { path: string };
  WebAppShare: { requestId: string; text: string; link: string };
  WebAppMaxShare: { requestId: string; text: string; link: string };
  WebAppSetupScreenCaptureBehavior: {
    requestId: string;
    isScreenCaptureEnabled: boolean;
  };
  WebAppHapticFeedbackImpact: {
    requestId: string;
    impactStyle: ImpactStyle;
    disableVibrationFallback: boolean;
  };
  WebAppHapticFeedbackNotification: {
    requestId: string;
    notificationType: NotificationType;
    disableVibrationFallback: boolean;
  };
  WebAppHapticFeedbackSelectionChange: {
    requestId: string;
    disableVibrationFallback: boolean;
  };
  WebAppOpenCodeReader: { requestId: string; fileSelect: boolean };
  WebAppDownloadFile: never;
  WebAppCopyText: never;
}

// Основной WebApp интерфейс
export interface WebApp {
  // Основные свойства
  readonly initData: string;
  readonly initDataUnsafe: WebAppInitData;
  readonly platform: 'ios' | 'android' | 'desktop' | 'web' | null;
  readonly version: string;

  // Менеджеры
  readonly BackButton: BackButton;
  readonly ScreenCapture: ScreenCapture;
  readonly HapticFeedback: HapticFeedback;
  readonly SecureStorage: Storage;
  readonly DeviceStorage: Storage;
  readonly BiometricManager: BiometricManager;

  // Состояния
  readonly isVerticalSwipesEnabled: boolean;

  // Методы событий
  onEvent<K extends keyof WebAppEventMap>(
    eventName: K,
    callback: (data: WebAppEventMap[K]) => void,
  ): () => void;

  offEvent<K extends keyof WebAppEventMap>(
    eventName: K,
    callback: (data: WebAppEventMap[K]) => void,
  ): void;

  // Основные методы
  ready(): void;
  close(): void;

  // Запросы данных
  requestContact(): Promise<{ phone: string }>;

  // Настройки UI
  enableClosingConfirmation(): void;
  disableClosingConfirmation(): void;

  // Навигация
  openLink(url: string): void;
  openMaxLink(url: string): void;

  // Шаринг
  shareContent(content: ShareContent): Promise<ShareResult>;
  shareMaxContent(content: ShareContent): Promise<ShareResult>;

  // Файлы
  downloadFile(url: string, fileName: string): Promise<DownloadFileResult>;

  // QR код
  openCodeReader(fileSelect?: boolean): Promise<CodeReaderResult>;

  // Яркость экрана
  requestScreenMaxBrightness(): Promise<void>;
  restoreScreenBrightness(): Promise<void>;

  // Свайпы
  enableVerticalSwipes(): Promise<{ allowVerticalSwipes: boolean }>;
  disableVerticalSwipes(): Promise<{ allowVerticalSwipes: boolean }>;
}

// Глобальное объявление
declare global {
  interface Window {
    WebApp: WebApp;
  }
}

export {};
