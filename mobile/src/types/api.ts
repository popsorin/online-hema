/**
 * API Type Definitions
 *
 * Types matching the backend API responses
 */

// Content types
export interface Resource {
  id: number;
  author_id?: number;
  title: string;
  description: string;
  publication_year?: number;
  cover_image_url?: string;
  author_name?: string;
}

export interface Section {
  id: number;
  resource_id: number;
  parent_id?: number;
  kind: string;
  title: string;
  description: string;
  position: number;
}

export interface Item {
  id: number;
  section_id: number;
  kind: string;
  title: string;
  description: string;
  position: number;
  attributes?: Record<string, string>;
}

// Pagination types
export interface PaginatedResponse<T> {
  data: T[];
  page: number;
  page_size: number;
  total_count: number;
  total_pages: number;
}

export interface PaginationParams {
  page?: number;
  page_size?: number;
}
