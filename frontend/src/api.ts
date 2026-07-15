import axios from 'axios';
import { API_BASE_URL } from './config';

const api = axios.create({
  baseURL: API_BASE_URL,
  headers: {
    'Content-Type': 'application/json',
  },
});

// Request Interceptor to attach Access Token
api.interceptors.request.use(
  (config) => {
    const stored = localStorage.getItem('sn_user');
    if (stored) {
      const user = JSON.parse(stored);
      if (user && user.accessToken) {
        config.headers['Authorization'] = `Bearer ${user.accessToken}`;
      }
    }
    return config;
  },
  (error) => {
    return Promise.reject(error);
  }
);

// Queue mechanism for handling multiple simultaneous 401s during refresh
let isRefreshing = false;
let failedQueue: any[] = [];

const processQueue = (error: any, token: string | null = null) => {
  failedQueue.forEach((prom) => {
    if (token) {
      prom.resolve(token);
    } else {
      prom.reject(error);
    }
  });
  failedQueue = [];
};

// Response Interceptor for handling 401 token refreshes
api.interceptors.response.use(
  (response) => response,
  async (error) => {
    const originalRequest = error.config;

    if (error.response?.status === 401 && !originalRequest._retry) {
      if (isRefreshing) {
        return new Promise((resolve, reject) => {
          failedQueue.push({ resolve, reject });
        })
          .then((token) => {
            originalRequest.headers['Authorization'] = `Bearer ${token}`;
            return api(originalRequest);
          })
          .catch((err) => {
            return Promise.reject(err);
          });
      }

      originalRequest._retry = true;
      isRefreshing = true;

      const stored = localStorage.getItem('sn_user');
      if (stored) {
        const user = JSON.parse(stored);
        if (user && user.refreshToken) {
          try {
            const res = await axios.post(`${API_BASE_URL}/api/v1/auth/refresh`, {
              refresh_token: user.refreshToken,
            });

            const access_token = res.data.data.access_token;
            const refresh_token = res.data.data.refresh_token;

            const updatedUser = {
              ...user,
              accessToken: access_token,
              refreshToken: refresh_token,
            };
            localStorage.setItem('sn_user', JSON.stringify(updatedUser));

            // Notify React app that user state was updated
            window.dispatchEvent(new CustomEvent('sn-user-updated', { detail: updatedUser }));

            originalRequest.headers['Authorization'] = `Bearer ${access_token}`;
            processQueue(null, access_token);
            return api(originalRequest);
          } catch (refreshError) {
            processQueue(refreshError, null);
            localStorage.removeItem('sn_user');
            window.dispatchEvent(new Event('sn-user-logged-out'));
            return Promise.reject(refreshError);
          } finally {
            isRefreshing = false;
          }
        }
      }
    }

    return Promise.reject(error);
  }
);

export default api;
