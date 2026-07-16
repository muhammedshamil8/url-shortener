// Centralized configuration for the React frontend
const getApiBaseUrl = () => {
  const envUrl = import.meta.env.VITE_API_BASE_URL || '';
  
  // Safeguard: If we are running in the browser on a non-localhost domain,
  // we should never try to hit localhost for the API.
  if (typeof window !== 'undefined' && window.location) {
    const isLocalhost = window.location.hostname === 'localhost' || window.location.hostname === '127.0.0.1';
    if (!isLocalhost && envUrl.includes('localhost')) {
      return '';
    }
  }
  return envUrl;
};

export const API_BASE_URL = getApiBaseUrl();

