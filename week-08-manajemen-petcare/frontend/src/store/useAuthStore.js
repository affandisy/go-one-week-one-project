import { create } from 'zustand';
import { persist } from 'zustand/middleware';

export const useAuthStore = create(
  persist(
    (set) => ({
      token: null,
      userRole: null,
      isAuthenticated: false,

      // Fungsi simulasi login (Nantinya diganti dengan pemanggilan API Axios/Fetch)
      login: async (username, password) => {
        // --- SIMULASI API CALL ---
        if (username === 'kasir' && password === 'rahasia') {
          set({
            token: 'mock-jwt-token-cashier-123',
            userRole: 'Cashier',
            isAuthenticated: true,
          });
          return { success: true };
        } else if (username === 'manajer' && password === 'rahasia') {
          set({
            token: 'mock-jwt-token-manager-456',
            userRole: 'Manager',
            isAuthenticated: true,
          });
          return { success: true };
        }
        
        return { success: false, message: 'Username atau password salah' };
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
      name: 'auth-storage', // Nama key di localStorage
    }
  )
);