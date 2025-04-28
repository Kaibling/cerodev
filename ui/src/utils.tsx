import { apiRequest } from './api';
export const api = {
  start: (name: string) => apiRequest(`/api/containers/${name}/start`, { method: 'POST' }),
  stop: (name: string) => apiRequest(`/api/containers/${name}/stop`, { method: 'POST' }),
  delete: (name: string) => apiRequest(`/api/containers/${name}`, { method: 'DELETE' }),
};