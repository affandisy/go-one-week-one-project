import axios from 'axios';
import { useAuthStore } from '../store/useAuthStore';

// 1. Inisiasi instance Axios
// Di produksi, ganti baseURL dengan environment variable (contoh: import.meta.env.VITE_API_URL)
const apiClient = axios.create({
  baseURL: 'http://localhost:3000/api/v1',
  headers: {
    'Content-Type': 'application/json',
  },
  timeout: 10000, // Timeout dalam 10 detik agar UI tidak hang jika server mati
});

// 2. REQUEST INTERCEPTOR: Menyisipkan Token JWT sebelum request dikirim
apiClient.interceptors.request.use(
  (config) => {
    // Mengambil token langsung dari state Zustand (tanpa hooks)
    const token = useAuthStore.getState().token;
    
    if (token) {
      config.headers['Authorization'] = `Bearer ${token}`;
    }
    
    return config;
  },
  (error) => {
    return Promise.reject(error);
  }
);

// 3. RESPONSE INTERCEPTOR: Menangani error secara global
apiClient.interceptors.response.use(
  (response) => {
    // Jika respons sukses, kembalikan datanya langsung
    return response;
  },
  (error) => {
    // Jika server mengembalikan status 401 (Unauthorized/Token Expired)
    if (error.response && error.response.status === 401) {
      console.warn('Sesi kedaluwarsa atau token tidak valid. Mengeluarkan pengguna...');
      
      // Panggil fungsi logout dari Zustand untuk membersihkan state dan localStorage
      useAuthStore.getState().logout();
      
      // Arahkan paksa kembali ke halaman login (opsional, karena ProtectedRoute biasanya sudah menangani ini)
      window.location.href = '/login';
    }

    // Ekstrak pesan error dari Golang agar mudah dibaca oleh komponen UI
    const errorMessage = error.response?.data?.error || 'Terjadi kesalahan pada peladen';
    return Promise.reject(new Error(errorMessage));
  }
);

export default apiClient;