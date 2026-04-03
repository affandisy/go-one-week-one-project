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

    const response = await fetch(`${API_URL}${endpoint}`, config);
    const responseData = await response.json();

    if (!response.ok) throw new Error(responseData.error || 'Terjadi kesalahan server');
    return responseData;
}