// src/components/auth/SessionGateModal.tsx
import React, { useEffect, useRef } from 'react';
import {
  useSelectedActivity,
  useIsSessionModalOpen,
  useUIActions,
} from '../../store/ui.store';
import { useAuth } from '../../hooks/useAuth';

interface Props {
  onProceed: () => void;
}

const SessionGateModal: React.FC<Props> = ({ onProceed }) => {
  const isOpen = useIsSessionModalOpen();
  const { setSessionModalOpen } = useUIActions();
  const selected = useSelectedActivity();
  const { isAuthenticated, login, loading } = useAuth();

  const proceedButtonRef = useRef<HTMLButtonElement | null>(null);

  useEffect(() => {
    if (isOpen && proceedButtonRef.current) proceedButtonRef.current.focus();
  }, [isOpen]);

  useEffect(() => {
    if (isOpen && !isAuthenticated && window.WebApp) {
      const initData = window.WebApp.initData;
      if (initData) {
        login(initData);
      }
    }
  }, [isOpen, isAuthenticated, login]);

  if (!isOpen) return null;

  const close = () => setSessionModalOpen(false);

  const handleBackdropClick = (e: React.MouseEvent) => {
    if (e.target === e.currentTarget) close();
  };

  const handleConfirm = () => {
    setSessionModalOpen(false);
    onProceed();
  };

  return (
    <div
      className="fixed inset-0 z-[60] flex items-end justify-center p-4 sm:items-center sm:p-6"
      onClick={handleBackdropClick}
      role="dialog"
      aria-modal="true"
    >
      <div className="absolute inset-0 bg-black/40" aria-hidden="true" />

      <div className="relative bg-white rounded-2xl max-w-sm w-full mx-auto shadow-xl z-[70]">
        <div className="p-5 border-b border-gray-100 flex items-start justify-between">
          <h2 className="text-lg font-bold text-gray-900">
            {loading ? 'Проверка доступа...' : 'Участие в мероприятии'}
          </h2>
          <button
            onClick={close}
            className="w-7 h-7 rounded-full bg-gray-100 flex items-center justify-center hover:bg-gray-200"
            aria-label="Закрыть"
            disabled={loading}
          >
            <svg
              className="w-3.5 h-3.5 text-gray-600"
              fill="none"
              stroke="currentColor"
              viewBox="0 0 24 24"
            >
              <path
                strokeLinecap="round"
                strokeLinejoin="round"
                strokeWidth={2}
                d="M6 18L18 6M6 6l12 12"
              />
            </svg>
          </button>
        </div>

        <div className="p-5">
          {selected && (
            <div className="bg-gray-50 rounded-xl p-3 mb-4">
              <h3 className="font-semibold text-gray-900 text-sm mb-1">
                {selected.title}
              </h3>
              <p className="text-gray-600 text-xs">{selected.description}</p>
              <div className="flex items-center justify-between mt-2 text-xs text-gray-500">
                <span>
                  {new Date(selected.date).toLocaleDateString('ru-RU')}
                </span>
                <span>{selected.location}</span>
              </div>
            </div>
          )}

          <div className="text-center mb-4">
            <div className="w-12 h-12 bg-gradient-to-r from-blue-500 to-purple-600 rounded-xl flex items-center justify-center mx-auto mb-3">
              {loading ? (
                <div className="animate-spin rounded-full h-6 w-6 border-b-2 border-white"></div>
              ) : (
                <svg
                  className="w-6 h-6 text-white"
                  fill="none"
                  stroke="currentColor"
                  viewBox="0 0 24 24"
                >
                  <path
                    strokeLinecap="round"
                    strokeLinejoin="round"
                    strokeWidth={2}
                    d="M12 15v2m-6 4h12a2 2 0 002-2v-6a2 2 0 00-2-2H6a2 2 0 00-2 2v6a2 2 0 002 2zm10-10V7a4 4 0 00-8 0v4h8z"
                  />
                </svg>
              )}
            </div>
            <p className="text-gray-600 text-sm">
              {loading
                ? 'Проверяем авторизацию...'
                : 'Для участия в мероприятии требуется авторизация'}
            </p>
          </div>

          <div className="space-y-2">
            <div className="flex gap-2">
              <button
                onClick={close}
                className="flex-1 px-4 py-2.5 bg-gray-100 text-gray-700 font-medium rounded-lg hover:bg-gray-200 disabled:opacity-50"
                disabled={loading}
              >
                Отмена
              </button>
              <button
                ref={proceedButtonRef}
                onClick={handleConfirm}
                className="flex-1 px-4 py-2.5 bg-gradient-to-r from-green-600 to-green-700 text-white font-medium rounded-lg hover:from-green-700 hover:to-green-800 disabled:opacity-50"
                disabled={loading || !isAuthenticated}
              >
                {loading ? '...' : 'Участвовать'}
              </button>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
};

export default SessionGateModal;
