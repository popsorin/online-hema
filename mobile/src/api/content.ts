/**
 * Content API
 *
 * Handles all content-related API calls (resources, sections, items).
 * These endpoints are publicly accessible without authentication.
 */

import apiClient from './client';
import type {
  Resource,
  Section,
  Item,
  PaginatedResponse,
  PaginationParams,
} from '@/types/api';

/**
 * Get a paginated list of resources
 */
export const getResources = async (
  params: PaginationParams = {},
): Promise<PaginatedResponse<Resource>> => {
  const response = await apiClient.get<PaginatedResponse<Resource>>(
    '/api/resources',
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
 * Get all root sections for a resource
 */
export const getSections = async (resourceId: number): Promise<Section[]> => {
  const response = await apiClient.get<Section[]>(
    `/api/resources/${resourceId}/sections`,
  );
  return response.data;
};

/**
 * Get all items for a section
 */
export const getItems = async (sectionId: number): Promise<Item[]> => {
  const response = await apiClient.get<Item[]>(
    `/api/sections/${sectionId}/items`,
  );
  return response.data;
};
