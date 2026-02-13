/**
 * Content API
 *
 * Handles all content-related API calls (fighting books, chapters, techniques).
 * These endpoints are publicly accessible without authentication.
 */

import apiClient from './client';
import type {
  FightingBook,
  Chapter,
  Technique,
  PaginatedResponse,
  PaginationParams,
} from '@/types/api';

/**
 * Get a paginated list of fighting books
 */
export const getFightingBooks = async (
  params: PaginationParams = {},
): Promise<PaginatedResponse<FightingBook>> => {
  const response = await apiClient.get<PaginatedResponse<FightingBook>>(
    '/api/fighting-books',
    {
      params: {
        page: params.page,
        page_size: params.page_size,
      },
    },
  );
  return response.data;
};

/**
 * Get all chapters for a fighting book
 */
export const getChapters = async (bookId: number): Promise<Chapter[]> => {
  const response = await apiClient.get<Chapter[]>(
    `/api/fighting-books/${bookId}/chapters`,
  );
  return response.data;
};

/**
 * Get all techniques for a chapter
 */
export const getTechniques = async (
  chapterId: number,
): Promise<Technique[]> => {
  const response = await apiClient.get<Technique[]>(
    `/api/chapters/${chapterId}/techniques`,
  );
  return response.data;
};
