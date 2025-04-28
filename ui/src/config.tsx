export const DEFAULT_API_BASE_URL = 'http://localhost:8088';
export const UI_VERSION = 'v0.0.1';

export const storage_prefix = 'cerodev';
export const storage_token = `${storage_prefix}_token`;
export const storage_api_url = `${storage_prefix}_api_url`;

export const get_base_api_url = () => {
  const url = localStorage.getItem(storage_api_url);
  if (url) {
    return url;
  } else {
    return DEFAULT_API_BASE_URL;
  }
};


