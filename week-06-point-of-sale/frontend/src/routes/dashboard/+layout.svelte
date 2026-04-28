<script lang="ts">
    import { page } from '$app/state';
    import { goto } from '$app/navigation';

    let { children } = $props();
    let currentPath = $derived(page.url.pathname);

    const menu = [
        { name: 'Ringkasan', path: '/dashboard', icon: '📊' },
        { name: 'Produk', path: '/dashboard/products', icon: '📦' },
        { name: 'Barang Masuk', path: '/dashboard/stock-in', icon: '📥' },
        { name: 'Laporan', path: '/dashboard/reports', icon: '📑' },
        { name: 'Pengguna', path: '/dashboard/users', icon: '👥' },
    ];

    function handleLogout() {
        localStorage.clear();
        goto('/login');
    }
</script>

<div class="flex h-screen bg-gray-100 overflow-hidden font-sans">
    <aside class="w-64 bg-gray-900 text-white flex flex-col shadow-xl">
        <div class="p-6 text-center border-b border-gray-800">
            <h2 class="text-2xl font-black tracking-tighter text-blue-400">POS KELONTONG</h2>
            <p class="text-xs text-gray-400 font-bold mt-1 uppercase">Panel Pemilik</p>
        </div>

        <nav class="flex-1 p-4 space-y-2 overflow-y-auto">
            {#each menu as item}
                <a href={item.path} 
                   class="flex items-center gap-3 px-4 py-3 rounded-xl font-bold transition-all
                   {currentPath === item.path ? 'bg-blue-600 text-white shadow-lg' : 'text-gray-400 hover:bg-gray-800 hover:text-gray-200'}">
                    <span class="text-xl">{item.icon}</span>
                    {item.name}
                </a>
            {/each}
        </nav>

        <div class="p-4 border-t border-gray-800">
            <button onclick={handleLogout} class="w-full py-3 bg-red-900/30 hover:bg-red-600 text-red-400 hover:text-white rounded-xl font-bold transition-all flex items-center justify-center gap-2">
                <span>🚪</span> Keluar
            </button>
        </div>
    </aside>

    <main class="flex-1 flex flex-col overflow-hidden">
        <header class="bg-white h-16 shadow-sm flex items-center px-8 justify-between z-10">
            <h3 class="font-black text-gray-700 text-lg uppercase tracking-wider">
                {menu.find(m => m.path === currentPath)?.name || 'Dashboard'}
            </h3>
            <div class="flex items-center gap-4">
                <div class="text-right hidden sm:block">
                    <p class="text-xs font-black text-gray-400 uppercase">Status Sistem</p>
                    <p class="text-sm font-bold text-green-600">Terhubung ke Database Lokal</p>
                </div>
            </div>
        </header>

        <section class="flex-1 overflow-y-auto p-8">
            {@render children()}
        </section>
    </main>
</div>