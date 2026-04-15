<script lang="ts">
    import { apiFetch } from '$lib/api';
    import { goto } from '$app/navigation';
    import { page } from '$app/state';

    let username = $state('');
    let password = $state('');
    let isLoading = $state(false);
    let errorMessage = $state('');

    let moduleId = page.params.moduleId;
    const errorParam = page.url.searchParams.get('error');
    if (errorParam === 'session_expired') {
        errorMessage = 'Sesi Anda telah berakhir karena alasan keamanan. Silakan masuk kembali.';
    }

    async function handleLogin(e: Event) {
        e.preventDefault();
        isLoading = true;
        errorMessage = '';

        try {
            const res = await apiFetch('/auth/login', {
                method: 'POST',
                data: { username, password }
            });
            
            // Simpan token ke LocalStorage untuk sesi yang persisten
            localStorage.setItem('token', res.data.accessToken);
            localStorage.setItem('username', res.data.user.username);
            
            // Arahkan ke halaman utama pembelajaran
            goto('/dashboard');

        } catch (err: any) {
            errorMessage = err.message || 'Username atau password salah.';
        } finally {
            isLoading = false;
        }
    }
</script>

<div class="min-h-screen flex items-center justify-center bg-slate-50 px-4">
    <div class="max-w-md w-full bg-white rounded-2xl shadow-sm border border-slate-100 p-8">
        <div class="text-center mb-8">
            <h1 class="text-3xl font-black text-slate-800">Selamat Datang</h1>
            <p class="text-slate-500 mt-2">Masuk untuk melanjutkan pelajaran Anda.</p>
        </div>

        <form onsubmit={handleLogin} class="space-y-5">
            {#if errorMessage}
                <div class="p-3 bg-red-50 text-red-700 rounded-lg text-sm font-medium border border-red-100">
                    {errorMessage}
                </div>
            {/if}

            <div>
                <label for="username" class="block text-sm font-bold text-slate-700 mb-1">Username</label>
                <input id="username" type="text" required bind:value={username}
                    class="w-full px-4 py-3 bg-slate-50 border border-slate-200 rounded-xl focus:ring-2 focus:ring-indigo-500 focus:bg-white transition-all font-medium" 
                    placeholder="Masukkan username Anda">
            </div>

            <div>
                <label for="password" class="block text-sm font-bold text-slate-700 mb-1">Password</label>
                <input id="password" type="password" required bind:value={password}
                    class="w-full px-4 py-3 bg-slate-50 border border-slate-200 rounded-xl focus:ring-2 focus:ring-indigo-500 focus:bg-white transition-all font-medium" 
                    placeholder="••••••••">
            </div>

            <button type="submit" disabled={isLoading}
                class="w-full py-3.5 px-4 bg-indigo-600 hover:bg-indigo-700 text-white font-bold rounded-xl shadow-sm disabled:opacity-70 transition-colors">
                {isLoading ? 'Memeriksa...' : 'Masuk'}
            </button>
        </form>

        <div class="mt-6 text-center text-sm text-slate-500 font-medium">
            Belum punya akun? 
            <a href="/register" class="text-indigo-600 hover:text-indigo-800 font-bold">Daftar gratis</a>
        </div>
    </div>
</div>