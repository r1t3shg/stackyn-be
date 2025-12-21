// API configuration
export const API_BASE_URL = import.meta.env.VITE_API_BASE_URL || 'http://localhost:8080';

// Log the API base URL to help debug (always log in production too)
console.log('API Base URL:', API_BASE_URL);

export const API_ENDPOINTS = {
  apps: `${API_BASE_URL}/api/v1/apps`,
  deployments: `${API_BASE_URL}/api/v1/deployments`,
  health: `${API_BASE_URL}/health`,
} as const;


