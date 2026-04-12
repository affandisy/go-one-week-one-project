<script lang="ts">
    import '../app.css';
    import { page } from '$app/stores';
    import { goto } from '$app/navigation';
    import { browser } from '$app/environment';

    // Svelte 5 Runes
    let { children } = $props();
    
    // Status loading diatur ke 'true' secara default untuk menahan render (mencegah flash)
    let isCheckingAuth = $state(true);

    $effect(() => {
        // Efek ini hanya berjalan di sisi browser, bukan di server (SSR)
        if (browser) {
            const token = localStorage.getItem('token');
            const currentPath = $page.url.pathname;
            
            // Daftar rute yang boleh diakses tanpa login
            const publicRoutes = ['/login', '/register', '/'];

            // 1. Jika TIDAK PUNYA token dan mencoba mengakses rute terproteksi
            if (!token && !publicRoutes.includes(currentPath)) {
                goto('/login');
            } 
            // 2. Jika SUDAH PUNYA token tapi mencoba mengakses halaman login/register
            else if (token && publicRoutes.includes(currentPath)) {
                goto('/dashboard');
                isCheckingAuth = false;
            } 
            // 3. Status aman, izinkan halaman dirender
            else {
                isCheckingAuth = false;
            }
        }
    });
</script>

{#if isCheckingAuth}
    <div class="min-h-screen flex flex-col items-center justify-center bg-slate-50">
        <div class="animate-spin rounded-full h-12 w-12 border-b-4 border-indigo-600 mb-4"></div>
        <p class="text-slate-500 font-bold text-sm">Memeriksa sesi...</p>
    </div>
{:else}
    {@render children()}
{/if}