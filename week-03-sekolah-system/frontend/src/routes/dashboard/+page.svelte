<script lang="ts">
    import { onMount } from 'svelte';
    import { goto } from '$app/navigation';

    onMount(() => {
        const token = localStorage.getItem('token');
        if (!token) {
            goto('/login');
        }
    });

    function logout() {
        localStorage.removeItem('token');
        goto('/login');
    }
</script>

<div class="min-h-screen bg-gray-50 flex flex-col">
    <nav class="bg-blue-600 text-white p-4 flex justify-between items-center shadow-md">
        <h1 class="text-xl font-bold">Dashboard SMA</h1>
        <button on:click={logout} class="bg-red-500 px-4 py-2 rounded hover:bg-red-600">Logout</button>
    </nav>
    <main class="flex-1 p-8">
        <h2 class="text-3xl font-semibold text-gray-800">Selamat Datang di Sistem Akademik</h2>
        <p class="text-gray-600 mt-2">Ini adalah kerangka awal untuk Phase 1.</p>
        
        <div class="mt-12">
    <h3 class="text-xl font-bold text-gray-800 mb-4">Menu Guru</h3>
    <div class="flex gap-4">
        <a href="/dashboard/grades" class="bg-white border border-gray-200 p-4 rounded shadow-sm hover:shadow-md transition flex flex-col items-center justify-center w-48 text-blue-600">
            <svg class="w-8 h-8 mb-2" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z"></path></svg>
            <span class="font-semibold">Input Nilai</span>
        </a>
    </div>
</div>

        <div class="grid grid-cols-3 gap-4 mt-8">
            <div class="bg-white p-6 rounded shadow border-t-4 border-blue-500">
                <h3 class="font-bold text-gray-700">Total Siswa</h3>
                <p class="text-2xl font-black text-blue-600">0</p>
            </div>
            <div class="bg-white p-6 rounded shadow border-t-4 border-green-500">
                <h3 class="font-bold text-gray-700">Tahun Ajaran Aktif</h3>
                <p class="text-2xl font-black text-green-600">2025/2026</p>
            </div>
        </div>
    </main>
</div>