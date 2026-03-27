// Sesuaikan dengan port backend Golang Anda
const API_URL = 'http://localhost:3000/api/v1';

interface FetchOptions extends RequestInit {
    data?: any;
}

export async function apiFetch(endpoint: string, options: FetchOptions = {}) {
    let token = '';
    // SSR Check: SvelteKit berjalan di server dan browser. LocalStorage hanya ada di browser.
    if (typeof window !== 'undefined') {
        token = localStorage.getItem('token') || '';
    }

    const headers: Record<string, string> = {
        'Content-Type': 'application/json',
        ...(options.headers as Record<string, string>)
    };

    if (token) {
        headers['Authorization'] = `Bearer ${token}`;
    }

    const config: RequestInit = {
        ...options,
        headers,
    };

    if (options.data) {
        config.body = JSON.stringify(options.data);
    }

    try {
        const response = await fetch(`${API_URL}${endpoint}`, config);
        const responseData = await response.json();

        if (!response.ok) {
            throw new Error(responseData.error || 'Terjadi kesalahan pada server');
        }

        return responseData;
    } catch (error: any) {
        throw error;
    }
}