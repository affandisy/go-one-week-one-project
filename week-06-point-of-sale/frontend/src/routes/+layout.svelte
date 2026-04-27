<script lang="ts">
    import '../app.css';
    import { page } from '$app/state';
    import { goto } from '$app/navigation';
    import { browser } from '$app/environment';

    let { children } = $props();
    let isCheckingAuth = $state(true);

    $effect(() => {
        if (browser) {
            const token = localStorage.getItem('token');
            const role = localStorage.getItem('role');
            const currentPath = page.url.pathname;
            
            const publicRoutes = ['/login'];

            if (!token && !publicRoutes.includes(currentPath)) {
                // Belum login, arahkan ke login
                goto('/login');
            } else if (token && publicRoutes.includes(currentPath)) {
                // Sudah login, arahkan ke dashboard/kasir berdasarkan role
                isCheckingAuth = false;
                if (role === 'cashier') {
                    goto('/pos'); // Kasir langsung ke layar transaksi
                } else {
                    goto('/dashboard'); // Owner/Admin ke dashboard laporan
                }
            } else {
                isCheckingAuth = false;
            }
        }
    });
</script>

{#if isCheckingAuth}
    <div class="min-h-screen flex items-center justify-center bg-gray-50">
        <div class="animate-spin rounded-full h-16 w-16 border-t-4 border-b-4 border-blue-600"></div>
    </div>
{:else}
    {@render children()}
{/if}