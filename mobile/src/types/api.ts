/**
 * API Type Definitions
 *
 * Types matching the backend API responses
 */

// Content types
export interface FightingBook {
  id: number;
  sword_master_id: number;
  title: string;
  description: string;
  publication_year: number;
  cover_image_url: string | null;
  created_at: string;
  updated_at: string;
  sword_master_name: string;
}

export interface Chapter {
  id: number;
  fighting_book_id: number;
  chapter_number: number;
  title: string;
  description: string;
  created_at: string;
  updated_at: string;
}

export interface Technique {
  id: number;
  chapter_id: number;
  name: string;
  description: string;
  instructions: string;
  video_url: string | null;
  thumbnail_url: string | null;
  order_in_chapter: number;
  created_at: string;
  updated_at: string;
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

// Health check
export interface HealthResponse {
  status: string;
  timestamp: string;
}
