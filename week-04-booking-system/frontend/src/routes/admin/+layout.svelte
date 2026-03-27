<script lang="ts">
    import { page } from '$app/stores';
    import { goto } from '$app/navigation';

    let { children } = $props();

    function handleLogout() {
        localStorage.removeItem('token');
        goto('/login');
    }
</script>

<div class="min-h-screen bg-gray-100 flex">
    <aside class="w-64 bg-white border-r border-gray-200 shadow-sm hidden md:block">
        <div class="p-6">
            <h1 class="text-xl font-black text-blue-600">Admin Padel</h1>
        </div>
        <nav class="mt-2 px-4 space-y-2">
            <a href="/admin/courts" 
               class="block px-4 py-2.5 rounded-lg font-medium transition-colors {$page.url.pathname.includes('/courts') ? 'bg-blue-50 text-blue-700' : 'text-gray-600 hover:bg-gray-50 hover:text-gray-900'}">
                Manajemen Lapangan
            </a>
            <a href="/admin/pricing" 
               class="block px-4 py-2.5 rounded-lg font-medium transition-colors {$page.url.pathname.includes('/pricing') ? 'bg-blue-50 text-blue-700' : 'text-gray-600 hover:bg-gray-50 hover:text-gray-900'}">
                Aturan Harga
            </a>
        </nav>
        <div class="absolute bottom-0 w-64 p-4 border-t border-gray-100">
            <button onclick={handleLogout} class="w-full text-left px-4 py-2 text-red-600 font-medium hover:bg-red-50 rounded-lg transition-colors">
                Logout
            </button>
        </div>
    </aside>

    <main class="flex-1 p-8 overflow-y-auto">
        {@render children()}
    </main>
</div>