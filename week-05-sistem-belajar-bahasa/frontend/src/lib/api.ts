// src/lib/api.ts
import { goto } from '$app/navigation';
import { browser } from '$app/environment';

const API_URL = 'http://localhost:3000/api/v1';

export async function apiFetch(endpoint: string, options: RequestInit & { data?: any } = {}) {
    let token = typeof window !== 'undefined' ? localStorage.getItem('token') || '' : '';

    const headers: Record<string, string> = {
        'Content-Type': 'application/json',
        ...(options.headers as Record<string, string>)
    };

    if (token) headers['Authorization'] = `Bearer ${token}`;

    const config: RequestInit = { ...options, headers };
    if (options.data) config.body = JSON.stringify(options.data);

    try {
        const response = await fetch(`${API_URL}${endpoint}`, config);
        
        // --- LOGIKA INTERSEPTOR (Penanganan Token Kedaluwarsa) ---
        if (response.status === 401) {
            if (browser) {
                // Hapus data sesi yang sudah tidak valid
                localStorage.removeItem('token');
                localStorage.removeItem('username');
                
                // Arahkan paksa ke login dengan pesan error di URL (optional)
                goto('/login?error=session_expired');
            }
            throw new Error('Sesi Anda telah berakhir. Silakan login kembali.');
        }

        const responseData = await response.json();

        if (!response.ok) {
            throw new Error(responseData.error || 'Terjadi kesalahan pada server');
        }

        return responseData;
    } catch (err: any) {
        throw err;
    }
}