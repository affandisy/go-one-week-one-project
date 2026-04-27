import { goto } from '$app/navigation';
import { browser } from '$app/environment';

const API_URL = 'http://localhost:3000/api/v1';

export async function apiFetch(endpoint: string, options: RequestInit & { data?: any } = {}) {
    let token = '';
    if (browser) {
        token = localStorage.getItem('token') || '';
    }

    const headers: Record<string, string> = {
        'Content-Type': 'application/json',
        ...(options.headers as Record<string, string>)
    };

    if (token) {
        headers['Authorization'] = `Bearer ${token}`;
    }

    const config: RequestInit = { ...options, headers };
    if (options.data) {
        config.body = JSON.stringify(options.data);
    }

    try {
        const response = await fetch(`${API_URL}${endpoint}`, config);
        
        // Interseptor 401 Unauthorized (Token Kedaluwarsa / Ditolak)
        if (response.status === 401) {
            if (browser) {
                localStorage.removeItem('token');
                localStorage.removeItem('username');
                localStorage.removeItem('role');
                goto('/login?error=session_expired');
            }
            throw new Error('Sesi tidak valid.');
        }

        const responseData = await response.json();

        if (!response.ok) {
            throw new Error(responseData.error || 'Terjadi kesalahan pada server');
        }

        return responseData;

    } catch (err: any) {
        // Deteksi jika server mati atau internet terputus
        if (err.message.includes('Failed to fetch') || err.message.includes('NetworkError')) {
            throw new Error('Koneksi terputus. Sistem berjalan dalam mode offline.');
        }
        throw err;
    }
}