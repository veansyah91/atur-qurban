import axios from "axios";
import { useAuthStore } from "@/stores/auth.store";
import { refreshAccessToken } from "@/services/auth.service";

// Instance axios dengan base URL dari environment
const api = axios.create({
  baseURL: process.env.NEXT_PUBLIC_API_URL,
  headers: {
    "Content-Type": "application/json",
  },
});

// Request interceptor: tambahkan Bearer token dari store
api.interceptors.request.use((config) => {
  const token = useAuthStore.getState().accessToken;
  if (token) {
    config.headers.Authorization = `Bearer ${token}`;
  }
  return config;
});

// Response interceptor: handle 401 → coba refresh token, kalau gagal logout
api.interceptors.response.use(
  (response) => response,
  async (error) => {
    const originalRequest = error.config;

    // Handle 401: Unauthorized
    if (error.response?.status === 401 && !originalRequest._retry) {
      originalRequest._retry = true;

      try {
        const refreshTokenStr = useAuthStore.getState().refreshToken;
        
        if (!refreshTokenStr) {
           throw new Error("No refresh token");
        }

        const data = await refreshAccessToken(refreshTokenStr);
        useAuthStore.getState().setAuth(data.user, data.access_token, data.refresh_token);

        // Retry request
        originalRequest.headers.Authorization = `Bearer ${data.access_token}`;
        return api(originalRequest);
      } catch (refreshError) {
        useAuthStore.getState().clearAuth();
        if (typeof window !== "undefined") {
           window.location.href = "/login";
        }
      }
    }
    return Promise.reject(error);
  }
);

export default api;
