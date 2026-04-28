<script lang="ts">
    import { onMount } from 'svelte';
    import { apiFetch } from '$lib/api';

    let stats = $state({ total_transactions: 0, total_omzet: 0 });
    let lowStockItems = $state<any[]>([]);
    let isLoading = $state(true);

    onMount(async () => {
        try {
            const [salesRes, stockRes] = await Promise.all([
                apiFetch('/reports/sales/daily'),
                apiFetch('/reports/stocks/low')
            ]);
            stats = salesRes.data;
            lowStockItems = stockRes.data;
        } finally {
            isLoading = false;
        }
    });

    const formatRp = (n: number) => new Intl.NumberFormat('id-ID', { style: 'currency', currency: 'IDR', minimumFractionDigits: 0 }).format(n);
</script>

<div class="space-y-8">
    <div class="grid grid-cols-1 md:grid-cols-3 gap-6">
        <div class="bg-white p-6 rounded-3xl shadow-sm border border-gray-100">
            <p class="text-gray-400 font-bold text-xs uppercase mb-1">Omzet Hari Ini</p>
            <p class="text-3xl font-black text-blue-600">{formatRp(stats.total_omzet)}</p>
        </div>
        <div class="bg-white p-6 rounded-3xl shadow-sm border border-gray-100">
            <p class="text-gray-400 font-bold text-xs uppercase mb-1">Transaksi Berhasil</p>
            <p class="text-3xl font-black text-gray-800">{stats.total_transactions} <span class="text-sm text-gray-400 font-bold">Nota</span></p>
        </div>
        <div class="bg-orange-50 p-6 rounded-3xl shadow-sm border border-orange-100">
            <p class="text-orange-600 font-bold text-xs uppercase mb-1">Peringatan Stok</p>
            <p class="text-3xl font-black text-orange-700">{lowStockItems.length} <span class="text-sm text-orange-500 font-bold">Barang Menipis</span></p>
        </div>
    </div>

    <div class="bg-white rounded-3xl shadow-sm border border-gray-100 overflow-hidden">
        <div class="p-6 border-b border-gray-50 flex justify-between items-center">
            <h4 class="font-black text-gray-800 uppercase tracking-tight">Perlu Restock Segera</h4>
            <a href="/dashboard/stock-in" class="text-blue-600 font-bold text-sm hover:underline">Tambah Stok &rarr;</a>
        </div>
        <div class="overflow-x-auto">
            <table class="w-full text-left">
                <thead class="bg-gray-50 text-gray-400 text-xs font-black uppercase">
                    <tr>
                        <th class="px-6 py-4">Nama Produk</th>
                        <th class="px-6 py-4">Sisa Stok</th>
                        <th class="px-6 py-4">Batas Minimum</th>
                        <th class="px-6 py-4">Status</th>
                    </tr>
                </thead>
                <tbody class="divide-y divide-gray-50">
                    {#each lowStockItems as item}
                        <tr class="hover:bg-gray-50/50">
                            <td class="px-6 py-4 font-bold text-gray-700">{item.name}</td>
                            <td class="px-6 py-4 font-black text-red-600">{item.stock} {item.unit}</td>
                            <td class="px-6 py-4 text-gray-500 font-medium">{item.min_stock}</td>
                            <td class="px-6 py-4">
                                <span class="px-3 py-1 bg-red-100 text-red-700 rounded-full text-xs font-black">KRITIS</span>
                            </td>
                        </tr>
                    {/each}
                </tbody>
            </table>
        </div>
    </div>
</div>