<script lang="ts">
    import { onMount } from 'svelte';
    import { apiFetch } from '$lib/api';
    import { goto } from '$app/navigation';

    // State Svelte 5 Runes
    let modules = $state<any[]>([]);
    let isLoading = $state(true);
    let errorMessage = $state('');
    let username = $state('');

    onMount(async () => {
        // Proteksi Halaman: Cek apakah user sudah login
        const token = localStorage.getItem('token');
        if (!token) {
            goto('/login');
            return;
        }

        username = localStorage.getItem('username') || 'Pelajar';
        await fetchDashboardData();
    });

    async function fetchDashboardData() {
        isLoading = true;
        errorMessage = '';
        try {
            // Mengambil data dari endpoint Phase 4 Backend yang sudah kita buat
            const res = await apiFetch('/progress');
            // Urutkan modul berdasarkan Level
            modules = (res.data || []).sort((a: any, b: any) => a.level_order - b.level_order);
        } catch (err: any) {
            errorMessage = err.message || 'Gagal memuat data pembelajaran.';
        } finally {
            isLoading = false;
        }
    }

    function handleLogout() {
        localStorage.removeItem('token');
        localStorage.removeItem('username');
        goto('/login');
    }

    function startLearning(moduleId: string) {
        goto(`/learn/${moduleId}`);
    }
</script>

<div class="min-h-screen bg-slate-50 font-sans">
    <nav class="bg-white border-b border-slate-200 px-4 py-4 sm:px-6 lg:px-8 flex justify-between items-center sticky top-0 z-10">
        <div class="flex items-center gap-2">
            <div class="w-8 h-8 bg-indigo-600 text-white rounded-lg flex items-center justify-center font-black text-xl">
                A
            </div>
            <span class="font-black text-xl text-slate-800 tracking-tight">LingoLearn</span>
        </div>
        <div class="flex items-center gap-4">
            <span class="text-sm font-bold text-slate-600 hidden sm:block">Halo, {username}!</span>
            <button onclick={handleLogout} class="text-sm font-bold text-red-500 hover:text-red-700 transition-colors">
                Keluar
            </button>
        </div>
    </nav>

    <main class="max-w-3xl mx-auto py-10 px-4 sm:px-6">
        <div class="mb-8 text-center sm:text-left">
            <h1 class="text-3xl sm:text-4xl font-black text-slate-800">Peta Belajar Anda</h1>
            <p class="text-slate-500 mt-2 font-medium">Selesaikan level saat ini dengan nilai minimal 80% untuk membuka level berikutnya.</p>
        </div>

        {#if errorMessage}
            <div class="p-4 bg-red-50 text-red-700 rounded-xl mb-6 border border-red-100 font-medium text-center">
                {errorMessage}
                <button onclick={fetchDashboardData} class="ml-2 underline font-bold">Coba Lagi</button>
            </div>
        {/if}

        {#if isLoading}
            <div class="space-y-4 animate-pulse">
                <div class="h-32 bg-slate-200 rounded-2xl w-full"></div>
                <div class="h-32 bg-slate-200 rounded-2xl w-full"></div>
                <div class="h-32 bg-slate-200 rounded-2xl w-full"></div>
            </div>
        {:else}
            <div class="space-y-5">
                {#each modules as mod}
                    <div class="relative bg-white rounded-2xl p-6 border transition-all duration-200
                        {mod.is_locked ? 'border-slate-200 shadow-sm opacity-75' : 'border-indigo-100 shadow-md hover:shadow-lg'}">
                        
                        <div class="flex flex-col sm:flex-row items-start sm:items-center justify-between gap-4">
                            
                            <div class="flex items-start gap-4">
                                <div class="w-12 h-12 rounded-full flex items-center justify-center font-black text-lg flex-shrink-0
                                    {mod.completed ? 'bg-green-100 text-green-600' : (mod.is_locked ? 'bg-slate-100 text-slate-400' : 'bg-indigo-100 text-indigo-600')}">
                                    {mod.level_order}
                                </div>
                                
                                <div>
                                    <h2 class="text-xl font-black text-slate-800 flex items-center gap-2">
                                        {mod.module_title}
                                        {#if mod.is_locked}
                                            <svg class="w-5 h-5 text-slate-400" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 15v2m-6 4h12a2 2 0 002-2v-6a2 2 0 00-2-2H6a2 2 0 00-2 2v6a2 2 0 002 2zm10-10V7a4 4 0 00-8 0v4h8z"></path></svg>
                                        {/if}
                                    </h2>
                                    
                                    <div class="mt-1 flex flex-wrap gap-2 sm:gap-4 text-sm font-medium">
                                        {#if mod.completed}
                                            <span class="text-green-600 bg-green-50 px-2.5 py-0.5 rounded-full border border-green-100">Selesai</span>
                                        {:else if !mod.is_locked}
                                            <span class="text-indigo-600 bg-indigo-50 px-2.5 py-0.5 rounded-full border border-indigo-100">Sedang Belajar</span>
                                        {:else}
                                            <span class="text-slate-500 bg-slate-50 px-2.5 py-0.5 rounded-full border border-slate-200">Terkunci</span>
                                        {/if}

                                        {#if mod.attempts > 0}
                                            <span class="text-slate-600">Skor Terbaik: <span class="font-black {mod.best_score >= 80 ? 'text-green-600' : 'text-orange-500'}">{mod.best_score}%</span></span>
                                        {/if}
                                    </div>
                                </div>
                            </div>

                            <div class="w-full sm:w-auto mt-2 sm:mt-0">
                                <button 
                                    disabled={mod.is_locked}
                                    onclick={() => startLearning(mod.module_id)}
                                    class="w-full sm:w-auto px-6 py-3 rounded-xl font-bold transition-colors
                                    {mod.is_locked 
                                        ? 'bg-slate-100 text-slate-400 cursor-not-allowed' 
                                        : 'bg-indigo-600 text-white hover:bg-indigo-700 shadow-sm'}">
                                    {mod.is_locked ? 'Terkunci' : (mod.completed ? 'Ulangi Belajar' : 'Mulai Belajar')}
                                </button>
                            </div>

                        </div>
                    </div>
                {/each}
            </div>
        {/if}
    </main>
</div>