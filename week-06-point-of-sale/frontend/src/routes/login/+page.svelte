<script lang="ts">
    import { apiFetch } from '$lib/api';
    import { goto } from '$app/navigation';
    import { page } from '$app/state';

    let username = $state('');
    let pin = $state('');
    let isLoading = $state(false);
    let errorMessage = $state('');

    // Tangkap error session expired dari URL
    const errorParam = page.url.searchParams.get('error');
    if (errorParam === 'session_expired') {
        errorMessage = 'Sesi Anda telah berakhir. Silakan masuk kembali.';
    }

    async function handleLogin(e: Event) {
        e.preventDefault();
        isLoading = true;
        errorMessage = '';

        try {
            const res = await apiFetch('/auth/login', {
                method: 'POST',
                data: { username, pin }
            });
            
            const { token, user } = res.data;
            
            // Simpan sesi
            localStorage.setItem('token', token);
            localStorage.setItem('username', user.username);
            localStorage.setItem('role', user.role);
            
            // Arahkan berdasarkan peran pengguna (Role)
            if (user.role === 'cashier') {
                goto('/pos'); // Layar Kasir Checkout
            } else {
                goto('/dashboard'); // Layar Manajemen/Pemilik
            }

        } catch (err: any) {
            errorMessage = err.message || 'Username atau PIN salah.';
        } finally {
            isLoading = false;
        }
    }
</script>

<div class="min-h-screen flex items-center justify-center bg-blue-900 p-4">
    <div class="max-w-md w-full bg-white rounded-3xl shadow-2xl overflow-hidden">
        <div class="p-8 text-center bg-blue-50 border-b border-blue-100">
            <div class="w-20 h-20 bg-blue-600 text-white rounded-full flex items-center justify-center mx-auto mb-4 shadow-lg">
                <svg class="w-10 h-10" fill="none" stroke="currentColor" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M16 11V7a4 4 0 00-8 0v4M5 9h14l1 12H4L5 9z"></path></svg>
            </div>
            <h1 class="text-3xl font-black text-gray-800">POS Cepat</h1>
            <p class="text-gray-500 font-medium mt-1">Sistem Kasir Toko Kelontong</p>
        </div>

        <div class="p-8">
            <form onsubmit={handleLogin} class="space-y-6">
                {#if errorMessage}
                    <div class="p-4 bg-red-50 border-l-4 border-red-500 text-red-700 font-bold rounded-r-lg">
                        {errorMessage}
                    </div>
                {/if}

                <div>
                    <label for="username" class="block text-sm font-bold text-gray-700 mb-2">Username Kasir / Admin</label>
                    <input id="username" type="text" required bind:value={username}
                        class="w-full px-5 py-4 bg-gray-50 border-2 border-gray-200 rounded-2xl focus:border-blue-500 focus:ring-0 text-lg font-medium transition-colors" 
                        placeholder="Contoh: admin" autocomplete="off">
                </div>

                <div>
                    <label for="pin" class="block text-sm font-bold text-gray-700 mb-2">PIN Keamanan</label>
                    <input id="pin" type="password" required bind:value={pin} inputmode="numeric"
                        class="w-full px-5 py-4 bg-gray-50 border-2 border-gray-200 rounded-2xl focus:border-blue-500 focus:ring-0 text-2xl tracking-widest text-center font-black transition-colors" 
                        placeholder="••••••">
                </div>

                <button type="submit" disabled={isLoading}
                    class="w-full py-4 px-6 bg-blue-600 hover:bg-blue-700 text-white font-black text-xl rounded-2xl shadow-xl hover:shadow-2xl active:scale-95 disabled:opacity-70 disabled:active:scale-100 transition-all">
                    {isLoading ? 'Memeriksa...' : 'MASUK'}
                </button>
            </form>
        </div>
    </div>
</div>