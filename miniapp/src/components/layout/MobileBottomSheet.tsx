// src/components/MobileBottomSheet.tsx
import React from 'react';
import { useSelectedActivity, useUIActions } from '../../store/ui.store';

const MobileBottomSheet: React.FC = () => {
  const selected = useSelectedActivity();
  const { setSelected, setSessionModalOpen } = useUIActions();

  if (!selected) return null;

  const handleClose = () => setSelected(null);
  const handleOpenActivity = () => setSessionModalOpen(true);

  return (
    <div className="fixed bottom-0 left-0 right-0 z-50">
      <div
        className="absolute inset-0 bg-black/20"
        onClick={handleClose}
        aria-hidden="true"
      />

      <div className="relative bg-white rounded-t-2xl p-4 m-1 shadow-xl">
        <div
          role="button"
          tabIndex={0}
          onClick={handleClose}
          className="absolute top-2 left-1/2 transform -translate-x-1/2 w-12 h-1.5 bg-gray-300 rounded-full"
          aria-label="Закрыть лист"
        />

        <div className="pt-6 space-y-3">
          <div className="flex items-start justify-between">
            <div className="pr-6">
              <h2 className="font-bold text-gray-900 text-base">
                {selected.title}
              </h2>
              <p className="text-gray-600 text-sm mt-1">
                {selected.description}
              </p>
            </div>
            <span
              className={`px-2 py-1 rounded-full text-xs font-medium ${
                selected.status === 'open'
                  ? 'bg-green-100 text-green-800'
                  : 'bg-red-100 text-red-800'
              }`}
            >
              {selected.status === 'open' ? 'Открыто' : 'Закрыто'}
            </span>
          </div>

          <div className="space-y-2 text-sm text-gray-600">
            <div className="flex items-center">
              <svg className="w-4 h-4 mr-2" fill="currentColor" viewBox="0 0 20 20">
                <path fillRule="evenodd" d="M5.05 4.05a7 7 0 119.9 9.9L10 18.9l-4.95-4.95a7 7 0 010-9.9zM10 11a2 2 0 100-4 2 2 0 000 4z" clipRule="evenodd" />
              </svg>
              <span>{selected.distanceText}</span>
            </div>
            
            <div className="flex items-center">
              <svg className="w-4 h-4 mr-2" fill="currentColor" viewBox="0 0 20 20">
                <path fillRule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zm1-12a1 1 0 10-2 0v4a1 1 0 00.293.707l2.828 2.829a1 1 0 101.415-1.415L11 9.586V6z" clipRule="evenodd" />
              </svg>
              <span>{new Date(selected.date).toLocaleDateString('ru-RU')}</span>
            </div>

            <div className="flex items-center">
              <svg className="w-4 h-4 mr-2" fill="currentColor" viewBox="0 0 20 20">
                <path d="M13 6a3 3 0 11-6 0 3 3 0 016 0zM18 8a2 2 0 11-4 0 2 2 0 014 0zM14 15a4 4 0 00-8 0v3h8v-3z" />
              </svg>
              <span>{selected.currentVolunteers} / {selected.maxVolunteers} волонтеров</span>
            </div>
          </div>

          <div className="flex gap-2 pt-2">
            <button
              className="flex-1 px-4 py-2.5 rounded-lg bg-gray-100 text-gray-700 font-medium hover:bg-gray-200"
              onClick={handleClose}
            >
              Закрыть
            </button>
            <button
              className="flex-1 px-4 py-2.5 rounded-lg bg-blue-600 text-white font-medium hover:bg-blue-700"
              onClick={handleOpenActivity}
            >
              Участвовать
            </button>
          </div>
        </div>
      </div>
    </div>
  );
};

export default MobileBottomSheet;