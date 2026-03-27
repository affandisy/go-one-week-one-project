<script lang="ts">
    import { apiFetch } from '$lib/api';
    import { goto } from '$app/navigation';

    // Svelte 5 Runes untuk state reaktif
    let email = $state('');
    let password = $state('');
    let errorMessage = $state('');
    let isLoading = $state(false);

    async function handleLogin(e: Event) {
        e.preventDefault(); // Mencegah form refresh halaman
        isLoading = true;
        errorMessage = '';

        try {
            const res = await apiFetch('/auth/login', {
                method: 'POST',
                data: { email, password }
            });
            
            // Simpan token ke LocalStorage
            localStorage.setItem('token', res.token);
            
            // Redirect ke halaman Admin (Untuk simulasi, kita arahkan ke courts)
            goto('/admin/courts');
        } catch (err: any) {
            errorMessage = err.message || 'Login gagal. Periksa kembali email dan password Anda.';
        } finally {
            isLoading = false;
        }
    }
</script>

<div class="min-h-screen flex items-center justify-center bg-gray-50 py-12 px-4 sm:px-6 lg:px-8">
    <div class="max-w-md w-full space-y-8 bg-white p-8 rounded-xl shadow-sm border border-gray-100">
        <div>
            <h2 class="mt-6 text-center text-3xl font-extrabold text-gray-900">Login Padel Booking</h2>
            <p class="mt-2 text-center text-sm text-gray-600">Masuk untuk mengelola lapangan</p>
        </div>
        
        <form class="mt-8 space-y-6" onsubmit={handleLogin}>
            {#if errorMessage}
                <div class="p-3 bg-red-50 text-red-700 rounded-md text-sm border border-red-200">
                    {errorMessage}
                </div>
            {/if}

            <div class="rounded-md shadow-sm space-y-4">
                <div>
                    <label for="email" class="block text-sm font-medium text-gray-700">Email Address</label>
                    <input id="email" type="email" required bind:value={email}
                        class="appearance-none relative block w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-blue-500 focus:border-blue-500 focus:z-10 sm:text-sm mt-1" 
                        placeholder="admin@padel.com">
                </div>
                <div>
                    <label for="password" class="block text-sm font-medium text-gray-700">Password</label>
                    <input id="password" type="password" required bind:value={password}
                        class="appearance-none relative block w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-blue-500 focus:border-blue-500 focus:z-10 sm:text-sm mt-1" 
                        placeholder="••••••••">
                </div>
            </div>

            <div>
                <button type="submit" disabled={isLoading}
                    class="group relative w-full flex justify-center py-2.5 px-4 border border-transparent text-sm font-medium rounded-md text-white bg-blue-600 hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500 disabled:bg-blue-300 transition-colors">
                    {isLoading ? 'Memproses...' : 'Sign In'}
                </button>
            </div>
        </form>
    </div>
</div>