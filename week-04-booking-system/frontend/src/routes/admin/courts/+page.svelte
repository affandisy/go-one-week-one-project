<script lang="ts">
    import { onMount } from 'svelte';
    import { apiFetch } from '$lib/api';

    // State Data
    let courts = $state<any[]>([]);
    
    // State Form
    let newName = $state('');
    let newType = $state('indoor');
    let newStatus = $state('active');
    
    // UI State
    let isLoading = $state(false);
    let message = $state('');
    let isError = $state(false);

    onMount(() => {
        fetchCourts();
    });

    async function fetchCourts() {
        try {
            const res = await apiFetch('/admin/courts');
            courts = res.data || [];
        } catch (err) {
            console.error("Gagal load courts:", err);
        }
    }

    async function createCourt(e: Event) {
        e.preventDefault();
        isLoading = true;
        message = '';

        try {
            await apiFetch('/admin/courts', {
                method: 'POST',
                data: {
                    name: newName,
                    type: newType,
                    status: newStatus
                }
            });

            isError = false;
            message = 'Lapangan berhasil ditambahkan!';
            
            // Reset form
            newName = '';
            newType = 'indoor';
            
            // Refresh data tabel
            await fetchCourts();
        } catch (err: any) {
            isError = true;
            message = err.message;
        } finally {
            isLoading = false;
        }
    }
</script>

<div class="max-w-6xl mx-auto">
    <div class="mb-8">
        <h2 class="text-3xl font-bold text-gray-800">Master Data: Lapangan Padel</h2>
        <p class="text-gray-500 mt-1">Kelola data fisik lapangan yang tersedia untuk di-booking.</p>
    </div>

    {#if message}
        <div class="p-4 mb-6 rounded-lg shadow-sm border-l-4 {isError ? 'bg-red-50 text-red-700 border-red-500' : 'bg-green-50 text-green-700 border-green-500'}">
            {message}
        </div>
    {/if}

    <div class="grid grid-cols-1 lg:grid-cols-3 gap-8">
        <div class="lg:col-span-1">
            <div class="bg-white p-6 rounded-xl shadow-sm border border-gray-100">
                <h3 class="text-lg font-bold text-gray-800 mb-4 border-b pb-2">Tambah Lapangan</h3>
                
                <form class="space-y-4" onsubmit={createCourt}>
                    <div>
                        <label for="courtName" class="block text-sm font-semibold text-gray-700 mb-1">Nama Lapangan</label>
                        <input id="courtName" type="text" required bind:value={newName} placeholder="Contoh: Court A Indoor" 
                            class="w-full border border-gray-300 rounded-lg p-2.5 focus:ring-2 focus:ring-blue-500" />
                    </div>

                    <div>
                        <label for="courtType" class="block text-sm font-semibold text-gray-700 mb-1">Tipe Lapangan</label>
                        <select id="courtType" bind:value={newType} class="w-full border border-gray-300 rounded-lg p-2.5 focus:ring-2 focus:ring-blue-500">
                            <option value="indoor">Indoor</option>
                            <option value="outdoor">Outdoor</option>
                        </select>
                    </div>

                    <div>
                        <label for="courtStatus" class="block text-sm font-semibold text-gray-700 mb-1">Status</label>
                        <select id="courtStatus" bind:value={newStatus} class="w-full border border-gray-300 rounded-lg p-2.5 focus:ring-2 focus:ring-blue-500">
                            <option value="active">Aktif (Bisa dibooking)</option>
                            <option value="maintenance">Maintenance</option>
                        </select>
                    </div>

                    <button type="submit" disabled={isLoading} 
                        class="w-full mt-2 bg-blue-600 text-white font-bold px-4 py-2.5 rounded-lg shadow hover:bg-blue-700 disabled:bg-blue-300 transition-colors">
                        {isLoading ? 'Menyimpan...' : 'Simpan Lapangan'}
                    </button>
                </form>
            </div>
        </div>

        <div class="lg:col-span-2">
            <div class="bg-white rounded-xl shadow-sm border border-gray-100 overflow-hidden">
                <div class="p-4 border-b border-gray-100 bg-gray-50">
                    <h3 class="font-bold text-gray-700">Daftar Lapangan Terdaftar</h3>
                </div>

                <div class="overflow-x-auto">
                    <table class="min-w-full divide-y divide-gray-200">
                        <thead class="bg-white">
                            <tr>
                                <th class="px-6 py-3 text-left text-xs font-bold text-gray-500 uppercase">Nama Lapangan</th>
                                <th class="px-6 py-3 text-left text-xs font-bold text-gray-500 uppercase">Tipe</th>
                                <th class="px-6 py-3 text-left text-xs font-bold text-gray-500 uppercase">Status</th>
                            </tr>
                        </thead>
                        <tbody class="divide-y divide-gray-100 bg-white">
                            {#each courts as court}
                                <tr class="hover:bg-gray-50 transition-colors">
                                    <td class="px-6 py-4 whitespace-nowrap font-bold text-gray-800">{court.name}</td>
                                    <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-600 capitalize">{court.type}</td>
                                    <td class="px-6 py-4 whitespace-nowrap text-sm">
                                        {#if court.status === 'active'}
                                            <span class="px-2 py-1 bg-green-100 text-green-800 rounded-full text-xs font-bold">Aktif</span>
                                        {:else}
                                            <span class="px-2 py-1 bg-yellow-100 text-yellow-800 rounded-full text-xs font-bold">Maintenance</span>
                                        {/if}
                                    </td>
                                </tr>
                            {:else}
                                <tr>
                                    <td colspan="3" class="px-6 py-8 text-center text-gray-500">Belum ada data lapangan.</td>
                                </tr>
                            {/each}
                        </tbody>
                    </table>
                </div>
            </div>
        </div>
    </div>
</div>