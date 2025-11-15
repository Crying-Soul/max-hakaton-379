// src/hooks/useCategories.ts
import { useState, useCallback, useEffect } from 'react';

export interface Category {
  id: number;
  name: string;
  color: string;
  description?: string;
}

interface UseCategoriesReturn {
  categories: Category[];
  selectedCategories: number[];
  loading: boolean;
  error: string | null;
  toggleCategory: (categoryId: number) => void;
  selectCategories: (categoryIds: number[]) => void;
  clearSelection: () => void;
  clearError: () => void;
}

// Mock данные - в реальном приложении загружаются с API
const MOCK_CATEGORIES: Category[] = [
  { id: 1, name: 'Экология', color: '#10B981', description: 'Уборка территорий, посадка деревьев' },
  { id: 2, name: 'Социальная помощь', color: '#3B82F6', description: 'Помощь пожилым, детям' },
  { id: 3, name: 'Образование', color: '#8B5CF6', description: 'Обучение, репетиторство' },
  { id: 4, name: 'Медицина', color: '#EF4444', description: 'Медицинская помощь, донорство' },
  { id: 5, name: 'Культура', color: '#F59E0B', description: 'Мероприятия, выставки' },
  { id: 6, name: 'Спорт', color: '#06B6D4', description: 'Спортивные мероприятия' },
];

export const useCategories = (): UseCategoriesReturn => {
  const [categories, setCategories] = useState<Category[]>([]);
  const [selectedCategories, setSelectedCategories] = useState<number[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    // В реальном приложении здесь будет запрос к API
    const loadCategories = async () => {
      try {
        setLoading(true);
        // Имитация загрузки с API
        await new Promise(resolve => setTimeout(resolve, 500));
        setCategories(MOCK_CATEGORIES);
      } catch (err: any) {
        setError('Failed to load categories');
      } finally {
        setLoading(false);
      }
    };

    loadCategories();
  }, []);

  const toggleCategory = useCallback((categoryId: number) => {
    setSelectedCategories(prev => 
      prev.includes(categoryId)
        ? prev.filter(id => id !== categoryId)
        : [...prev, categoryId]
    );
  }, []);

  const selectCategories = useCallback((categoryIds: number[]) => {
    setSelectedCategories(categoryIds);
  }, []);

  const clearSelection = useCallback(() => {
    setSelectedCategories([]);
  }, []);

  const clearError = useCallback(() => {
    setError(null);
  }, []);

  return {
    categories,
    selectedCategories,
    loading,
    error,
    toggleCategory,
    selectCategories,
    clearSelection,
    clearError,
  };
};