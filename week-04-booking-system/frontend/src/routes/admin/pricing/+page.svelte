<script lang="ts">
    import { onMount } from 'svelte';
    import { apiFetch } from '$lib/api';

    // State Data Master
    let courts = $state<any[]>([]);
    let pricingRules = $state<any[]>([]);
    
    // State Filter & Form
    let selectedCourtId = $state('');
    
    let newDayType = $state('weekday');
    let newStartTime = $state('06:00');
    let newEndTime = $state('18:00');
    let newBasePrice = $state(150000);
    let newMultiplier = $state(1.0);
    
    // UI State
    let isLoading = $state(false);
    let message = $state('');
    let isError = $state(false);

    // 1. Saat halaman dimuat, ambil daftar lapangan
    onMount(async () => {
        await fetchCourts();
    });

    // 2. Rune $effect: Bereaksi otomatis jika selectedCourtId berubah!
    $effect(() => {
        if (selectedCourtId) {
            fetchPricingRules(selectedCourtId);
        } else {
            pricingRules = [];
        }
    });

    async function fetchCourts() {
        try {
            const res = await apiFetch('/admin/courts');
            courts = res.data || [];
            // Auto-select lapangan pertama jika ada
            if (courts.length > 0) {
                selectedCourtId = courts[0].id;
            }
        } catch (err) {
            console.error("Gagal mengambil data lapangan:", err);
        }
    }

    async function fetchPricingRules(courtId: string) {
        try {
            const res = await apiFetch(`/admin/pricing-rules?court_id=${courtId}`);
            pricingRules = res.data || [];
        } catch (err) {
            console.error("Gagal mengambil data aturan harga:", err);
        }
    }

    async function createRule(e: Event) {
        e.preventDefault();
        if (!selectedCourtId) {
            isError = true;
            message = "Silakan pilih lapangan terlebih dahulu!";
            return;
        }

        isLoading = true;
        message = '';

        try {
            await apiFetch('/admin/pricing-rules', {
                method: 'POST',
                data: {
                    court_id: selectedCourtId,
                    day_type: newDayType,
                    start_time: newStartTime,
                    end_time: newEndTime,
                    base_price: Number(newBasePrice),
                    multiplier: Number(newMultiplier)
                }
            });

            isError = false;
            message = 'Aturan harga berhasil ditambahkan!';
            
            // Refresh tabel
            await fetchPricingRules(selectedCourtId);
        } catch (err: any) {
            isError = true;
            message = err.message;
        } finally {
            isLoading = false;
        }
    }

    // Fungsi bonus: Hapus aturan harga
    async function deleteRule(id: string) {
        if (!confirm('Apakah Anda yakin ingin menghapus aturan harga ini?')) return;
        
        try {
            await apiFetch(`/admin/pricing-rules/${id}`, { method: 'DELETE' });
            await fetchPricingRules(selectedCourtId); // Refresh data
        } catch (err) {
            alert('Gagal menghapus aturan harga');
        }
    }
</script>

<div class="max-w-7xl mx-auto">
    <div class="mb-8 flex flex-col md:flex-row md:items-end justify-between gap-4">
        <div>
            <h2 class="text-3xl font-bold text-gray-800">Aturan Harga (Pricing Rules)</h2>
            <p class="text-gray-500 mt-1">Atur harga berdasarkan hari dan jam sibuk (surge pricing).</p>
        </div>

        <div class="w-full md:w-72">
            <label for="courtSelect" class="block text-sm font-semibold text-gray-700 mb-1">Pilih Lapangan</label>
            <select id="courtSelect" bind:value={selectedCourtId} class="w-full border border-blue-300 bg-blue-50 text-blue-900 rounded-lg p-2.5 font-bold focus:ring-2 focus:ring-blue-500">
                <option value="">-- Pilih Lapangan --</option>
                {#each courts as court}
                    <option value={court.id}>{court.name} ({court.type})</option>
                {/each}
            </select>
        </div>
    </div>

    {#if message}
        <div class="p-4 mb-6 rounded-lg shadow-sm border-l-4 {isError ? 'bg-red-50 text-red-700 border-red-500' : 'bg-green-50 text-green-700 border-green-500'}">
            {message}
        </div>
    {/if}

    <div class="grid grid-cols-1 lg:grid-cols-3 gap-8">
        <div class="lg:col-span-1">
            <div class="bg-white p-6 rounded-xl shadow-sm border border-gray-100 sticky top-6">
                <h3 class="text-lg font-bold text-gray-800 mb-4 border-b pb-2">Tambah Harga Baru</h3>
                
                <form class="space-y-4" onsubmit={createRule}>
                    <div>
                        <label for="dayType" class="block text-sm font-semibold text-gray-700 mb-1">Tipe Hari</label>
                        <select id="dayType" bind:value={newDayType} class="w-full border border-gray-300 rounded-lg p-2.5 focus:ring-2 focus:ring-blue-500">
                            <option value="weekday">Senin - Jumat (Weekday)</option>
                            <option value="weekend">Sabtu - Minggu (Weekend)</option>
                        </select>
                    </div>

                    <div class="grid grid-cols-2 gap-4">
                        <div>
                            <label for="startTime" class="block text-sm font-semibold text-gray-700 mb-1">Jam Mulai</label>
                            <input id="startTime" type="time" required bind:value={newStartTime} 
                                class="w-full border border-gray-300 rounded-lg p-2.5 focus:ring-2 focus:ring-blue-500" />
                        </div>
                        <div>
                            <label for="endTime" class="block text-sm font-semibold text-gray-700 mb-1">Jam Selesai</label>
                            <input id="endTime" type="time" required bind:value={newEndTime} 
                                class="w-full border border-gray-300 rounded-lg p-2.5 focus:ring-2 focus:ring-blue-500" />
                        </div>
                    </div>

                    <div>
                        <label for="basePrice" class="block text-sm font-semibold text-gray-700 mb-1">Harga Dasar (Rp)</label>
                        <input id="basePrice" type="number" required min="0" step="1000" bind:value={newBasePrice} 
                            class="w-full border border-gray-300 rounded-lg p-2.5 focus:ring-2 focus:ring-blue-500" />
                    </div>

                    <div>
                        <label for="multiplier" class="block text-sm font-semibold text-gray-700 mb-1">Pengali Harga (Multiplier)</label>
                        <input id="multiplier" type="number" required min="1.0" step="0.1" bind:value={newMultiplier} 
                            class="w-full border border-gray-300 rounded-lg p-2.5 focus:ring-2 focus:ring-blue-500" />
                        <p class="text-xs text-gray-500 mt-1">1.0 = Harga Normal. 1.5 = Harga naik 50% (Prime Time).</p>
                    </div>

                    <button type="submit" disabled={isLoading || !selectedCourtId} 
                        class="w-full mt-2 bg-indigo-600 text-white font-bold px-4 py-2.5 rounded-lg shadow hover:bg-indigo-700 disabled:bg-indigo-300 transition-colors">
                        {isLoading ? 'Menyimpan...' : 'Simpan Harga'}
                    </button>
                </form>
            </div>
        </div>

        <div class="lg:col-span-2">
            <div class="bg-white rounded-xl shadow-sm border border-gray-100 overflow-hidden">
                <div class="p-4 border-b border-gray-100 bg-gray-50 flex justify-between items-center">
                    <h3 class="font-bold text-gray-700">Daftar Harga Lapangan Terpilih</h3>
                </div>

                {#if !selectedCourtId}
                    <div class="p-12 text-center text-gray-500 flex flex-col items-center">
                        <svg class="w-12 h-12 text-gray-300 mb-3" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 15l-2 5L9 9l11 4-5 2zm0 0l5 5M7.188 2.239l.777 2.897M5.136 7.965l-2.898-.777M13.95 4.05l-2.122 2.122m-5.657 5.656l-2.12 2.122"></path></svg>
                        Pilih lapangan di pojok kanan atas untuk melihat daftar harga.
                    </div>
                {:else}
                    <div class="overflow-x-auto">
                        <table class="min-w-full divide-y divide-gray-200">
                            <thead class="bg-white">
                                <tr>
                                    <th class="px-6 py-3 text-left text-xs font-bold text-gray-500 uppercase">Hari & Jam</th>
                                    <th class="px-6 py-3 text-left text-xs font-bold text-gray-500 uppercase">Harga Dasar</th>
                                    <th class="px-6 py-3 text-left text-xs font-bold text-gray-500 uppercase">Multiplier</th>
                                    <th class="px-6 py-3 text-right text-xs font-bold text-gray-500 uppercase">Harga Final</th>
                                    <th class="px-6 py-3 text-center text-xs font-bold text-gray-500 uppercase">Aksi</th>
                                </tr>
                            </thead>
                            <tbody class="divide-y divide-gray-100 bg-white">
                                {#each pricingRules as rule}
                                    <tr class="hover:bg-gray-50 transition-colors">
                                        <td class="px-6 py-4 whitespace-nowrap">
                                            <div class="font-bold text-gray-800 capitalize">{rule.day_type}</div>
                                            <div class="text-sm text-gray-500">{rule.start_time} - {rule.end_time}</div>
                                        </td>
                                        <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-600">
                                            Rp {rule.base_price.toLocaleString('id-ID')}
                                        </td>
                                        <td class="px-6 py-4 whitespace-nowrap text-sm">
                                            <span class="px-2 py-1 {rule.multiplier > 1.0 ? 'bg-orange-100 text-orange-800' : 'bg-gray-100 text-gray-800'} rounded text-xs font-bold">
                                                x {rule.multiplier}
                                            </span>
                                        </td>
                                        <td class="px-6 py-4 whitespace-nowrap text-sm font-black text-indigo-600 text-right">
                                            Rp {(rule.base_price * rule.multiplier).toLocaleString('id-ID')}
                                        </td>
                                        <td class="px-6 py-4 whitespace-nowrap text-center">
                                            <button onclick={() => deleteRule(rule.id)} class="text-red-500 hover:text-red-700 font-bold text-sm transition-colors">
                                                Hapus
                                            </button>
                                        </td>
                                    </tr>
                                {:else}
                                    <tr>
                                        <td colspan="5" class="px-6 py-8 text-center text-gray-500">Belum ada aturan harga untuk lapangan ini.</td>
                                    </tr>
                                {/each}
                            </tbody>
                        </table>
                    </div>
                {/if}
            </div>
        </div>
    </div>
</div>