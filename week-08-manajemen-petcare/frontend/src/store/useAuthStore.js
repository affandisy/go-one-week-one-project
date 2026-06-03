import { create } from 'zustand';
import { persist } from 'zustand/middleware';
import apiClient from '../utils/apiClient'; // Import klien HTTP yang baru dibuat

export const useAuthStore = create(
  persist(
    (set) => ({
      token: null,
      userRole: null,
      isAuthenticated: false,

      login: async (username, password) => {
        try {
          // Melakukan POST ke Golang Backend (/api/v1/auth/login)
          const response = await apiClient.post('/auth/login', {
            username,
            password
          });

          // Mengasumsikan Golang mengembalikan: { message: "...", data: { token: "..." } }
          // dan Anda perlu mem-parse JWT untuk mendapatkan Role (atau minta backend mengirimkannya)
          const token = response.data.data.token;
          
          // Cara sederhana men-decode JWT di frontend untuk mengambil peran (role)
          const payloadBase64 = token.split('.')[1];
          const decodedJson = atob(payloadBase64);
          const decodedPayload = JSON.parse(decodedJson);

          set({
            token: token,
            userRole: decodedPayload.role, // Diambil dari klaim "role" JWT
            isAuthenticated: true,
          });

          return { success: true };
        } catch (error) {
          // Error otomatis diekstrak oleh response interceptor
          return { success: false, message: error.message };
        }
      },

      logout: () => {
        set({
          token: null,
          userRole: null,
          isAuthenticated: false,
        });
      },
    }),
    {
      name: 'auth-storage',
    }
  )
);