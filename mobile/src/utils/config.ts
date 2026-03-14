/**
 * Application configuration
 *
 * For development, the API_URL defaults based on platform:
 * - Android emulator: 10.0.2.2 maps to host localhost
 * - iOS simulator: localhost works directly
 * - Physical device: use your machine's IP address
 * - Docker: use host.docker.internal
 */

interface Config {
  apiUrl: string;
  requestTimeout: number; // milliseconds
}

const getApiUrl = (): string => {
  if (__DEV__) {
    return 'http://10.0.2.2:8080'; // Default for Android emulator
  }
  return 'https://hema-lessons-api-564075903124.us-east1.run.app';
};

const config: Config = {
  apiUrl: getApiUrl(),
  requestTimeout: 30000, // 30 second timeout
};

export default config;
