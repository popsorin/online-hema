/**
 * Application configuration
 * 
 * Environment variables are loaded at build time.
 * For development, the API_URL can be configured via:
 * - Docker: Set in docker-compose.mobile.yml
 * - Local: Set in .env file or use default localhost
 */

interface Config {
  apiUrl: string;
  environment: 'development' | 'staging' | 'production';
  tokenRefreshThreshold: number; // seconds before expiry to refresh
  requestTimeout: number; // milliseconds
}

// Get environment or use sensible defaults
const getApiUrl = (): string => {
  // React Native doesn't have process.env at runtime like Node.js
  // For development, we'll use a configurable default
  // In production, this would be set via native build configs
  
  if (__DEV__) {
    // Android emulator: 10.0.2.2 maps to host localhost
    // iOS simulator: localhost works directly
    // Physical device: use your machine's IP address
    // Docker: use host.docker.internal
    return 'http://10.0.2.2:8080'; // Default for Android emulator
  }
  
  // Production API URL
  return 'https://hema-lessons-api-564075903124.us-east1.run.app';
};

const config: Config = {
  apiUrl: getApiUrl(),
  environment: __DEV__ ? 'development' : 'production',
  tokenRefreshThreshold: 60, // Refresh token 60 seconds before expiry
  requestTimeout: 30000, // 30 second timeout
};

export default config;

// Platform-specific API URLs for reference
export const PLATFORM_API_URLS = {
  androidEmulator: 'http://10.0.2.2:8080',
  iosSimulator: 'http://localhost:8080',
  docker: 'http://host.docker.internal:8080',
  // Set this to your machine's IP for physical device testing
  physicalDevice: 'http://192.168.1.X:8080',
};
