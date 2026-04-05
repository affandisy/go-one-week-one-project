<script lang="ts">
    import { apiFetch } from '$lib/api';
    import { goto } from '$app/navigation';

    let username = $state('');
    let password = $state('');
    let isLoading = $state(false);
    let errorMessage = $state('');
    let successMessage = $state('');

    async function handleRegister(e: Event) {
        e.preventDefault();
        isLoading = true;
        errorMessage = '';
        successMessage = '';

        try {
            await apiFetch('/auth/register', {
                method: 'POST',
                data: { username, password }
            });
            
            successMessage = 'Registrasi berhasil! Mengarahkan ke halaman login...';
            
            // Jeda 1.5 detik agar user bisa membaca pesan sukses
            setTimeout(() => {
                goto('/login');
            }, 1500);

        } catch (err: any) {
            errorMessage = err.message || 'Gagal mendaftar. Silakan coba lagi.';
        } finally {
            isLoading = false;
        }
    }
</script>

<div class="min-h-screen flex items-center justify-center bg-slate-50 px-4">
    <div class="max-w-md w-full bg-white rounded-2xl shadow-sm border border-slate-100 p-8">
        <div class="text-center mb-8">
            <h1 class="text-3xl font-black text-slate-800">Mulai Belajar</h1>
            <p class="text-slate-500 mt-2">Buat akun untuk melacak progres Anda.</p>
        </div>

        <form onsubmit={handleRegister} class="space-y-5">
            {#if errorMessage}
                <div class="p-3 bg-red-50 text-red-700 rounded-lg text-sm font-medium border border-red-100">
                    {errorMessage}
                </div>
            {/if}

            {#if successMessage}
                <div class="p-3 bg-green-50 text-green-700 rounded-lg text-sm font-medium border border-green-100">
                    {successMessage}
                </div>
            {/if}

            <div>
                <label for="username" class="block text-sm font-bold text-slate-700 mb-1">Username</label>
                <input id="username" type="text" required bind:value={username}
                    class="w-full px-4 py-3 bg-slate-50 border border-slate-200 rounded-xl focus:ring-2 focus:ring-indigo-500 focus:bg-white transition-all font-medium" 
                    placeholder="Contoh: learner123">
            </div>

            <div>
                <label for="password" class="block text-sm font-bold text-slate-700 mb-1">Password</label>
                <input id="password" type="password" required bind:value={password} minlength="6"
                    class="w-full px-4 py-3 bg-slate-50 border border-slate-200 rounded-xl focus:ring-2 focus:ring-indigo-500 focus:bg-white transition-all font-medium" 
                    placeholder="Minimal 6 karakter">
            </div>

            <button type="submit" disabled={isLoading}
                class="w-full py-3.5 px-4 bg-indigo-600 hover:bg-indigo-700 text-white font-bold rounded-xl shadow-sm disabled:opacity-70 transition-colors">
                {isLoading ? 'Memproses...' : 'Daftar Sekarang'}
            </button>
        </form>

        <div class="mt-6 text-center text-sm text-slate-500 font-medium">
            Sudah punya akun? 
            <a href="/login" class="text-indigo-600 hover:text-indigo-800 font-bold">Masuk di sini</a>
        </div>
    </div>
</div>