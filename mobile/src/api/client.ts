/**
 * API Client for HEMA Lessons Backend
 *
 * Handles all HTTP requests to the API with:
 * - Request/response configuration
 * - Error handling
 */

import axios, {AxiosInstance} from 'axios';
import config from '@/utils/config';

// Create axios instance with base configuration
const apiClient: AxiosInstance = axios.create({
  baseURL: config.apiUrl,
  timeout: config.requestTimeout,
  headers: {
    'Content-Type': 'application/json',
    Accept: 'application/json',
  },
});

export default apiClient;

// Type-safe error handler
export interface ApiError {
  message: string;
  status: number;
  data?: unknown;
}

export const handleApiError = (error: unknown): ApiError => {
  if (axios.isAxiosError(error)) {
    return {
      message: error.response?.data?.error || error.message,
      status: error.response?.status || 500,
      data: error.response?.data,
    };
  }
  return {
    message: 'An unexpected error occurred',
    status: 500,
  };
};
