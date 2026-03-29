<script lang="ts">
    import { onMount } from 'svelte';
    import { apiFetch } from '$lib/api';

    // State Data menggunakan Svelte 5 Runes
    let stats = $state<any>(null);
    let isLoading = $state(true);
    let errorMessage = $state('');

    onMount(async () => {
        try {
            const res = await apiFetch('/admin/dashboard-stats');
            stats = res.data;
        } catch (err: any) {
            errorMessage = err.message || 'Gagal memuat data statistik.';
        } finally {
            isLoading = false;
        }
    });

    // Helper format mata uang
    function formatIDR(amount: number) {
        return new Intl.NumberFormat('id-ID', { style: 'currency', currency: 'IDR', minimumFractionDigits: 0 }).format(amount || 0);
    }
</script>

<div class="max-w-7xl mx-auto">
    <div class="mb-8">
        <h2 class="text-3xl font-bold text-gray-800">Dashboard Utama</h2>
        <p class="text-gray-500 mt-1">Ringkasan performa bisnis lapangan Padel Anda.</p>
    </div>

    {#if errorMessage}
        <div class="p-4 mb-6 rounded-lg bg-red-50 text-red-700 border border-red-200">
            {errorMessage}
        </div>
    {/if}

    {#if isLoading}
        <div class="grid grid-cols-1 md:grid-cols-3 gap-6 animate-pulse">
            <div class="h-32 bg-gray-200 rounded-xl"></div>
            <div class="h-32 bg-gray-200 rounded-xl"></div>
            <div class="h-32 bg-gray-200 rounded-xl"></div>
        </div>
    {:else if stats}
        <div class="grid grid-cols-1 md:grid-cols-3 gap-6 mb-8">
            <div class="bg-gradient-to-br from-blue-500 to-blue-600 rounded-xl p-6 shadow-lg text-white">
                <div class="flex justify-between items-start">
                    <div>
                        <p class="text-blue-100 font-medium mb-1">Pendapatan Hari Ini</p>
                        <h3 class="text-3xl font-black">{formatIDR(stats.revenue.daily)}</h3>
                    </div>
                    <div class="p-3 bg-white/20 rounded-lg">
                        <svg class="w-6 h-6 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8c-1.657 0-3 .895-3 2s1.343 2 3 2 3 .895 3 2-1.343 2-3 2m0-8c1.11 0 2.08.402 2.599 1M12 8V7m0 1v8m0 0v1m0-1c-1.11 0-2.08-.402-2.599-1M21 12a9 9 0 11-18 0 9 9 0 0118 0z"></path></svg>
                    </div>
                </div>
            </div>

            <div class="bg-gradient-to-br from-indigo-500 to-indigo-600 rounded-xl p-6 shadow-lg text-white">
                <div class="flex justify-between items-start">
                    <div>
                        <p class="text-indigo-100 font-medium mb-1">Pendapatan Bulan Ini</p>
                        <h3 class="text-3xl font-black">{formatIDR(stats.revenue.monthly)}</h3>
                    </div>
                    <div class="p-3 bg-white/20 rounded-lg">
                        <svg class="w-6 h-6 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 7V3m8 4V3m-9 8h10M5 21h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v12a2 2 0 002 2z"></path></svg>
                    </div>
                </div>
            </div>

            <div class="bg-gradient-to-br from-emerald-500 to-emerald-600 rounded-xl p-6 shadow-lg text-white">
                <div class="flex justify-between items-start">
                    <div>
                        <p class="text-emerald-100 font-medium mb-1">Booking Aktif (Upcoming)</p>
                        <h3 class="text-3xl font-black">{stats.active_bookings} <span class="text-lg font-normal text-emerald-100">Jadwal</span></h3>
                        <p class="text-sm mt-2 text-emerald-100">Slot terpakai hari ini: <strong>{stats.today_occupancy_slots}</strong></p>
                    </div>
                    <div class="p-3 bg-white/20 rounded-lg">
                        <svg class="w-6 h-6 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M17 20h5v-2a3 3 0 00-5.356-1.857M17 20H7m10 0v-2c0-.656-.126-1.283-.356-1.857M7 20H2v-2a3 3 0 015.356-1.857M7 20v-2c0-.656.126-1.283.356-1.857m0 0a5.002 5.002 0 019.288 0M15 7a3 3 0 11-6 0 3 3 0 016 0zm6 3a2 2 0 11-4 0 2 2 0 014 0zM7 10a2 2 0 11-4 0 2 2 0 014 0z"></path></svg>
                    </div>
                </div>
            </div>
        </div>
    {/if}
</div>