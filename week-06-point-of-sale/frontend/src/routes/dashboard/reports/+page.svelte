<script lang="ts">
    import { onMount } from 'svelte';
    import { apiFetch } from '$lib/api';

    // State untuk Pemilih Tanggal
    const today = new Date().toISOString().split('T')[0]; // Format YYYY-MM-DD
    let selectedDate = $state(today);
    
    // State Data Laporan
    let summary = $state({ total_transactions: 0, total_omzet: 0, date: '' });
    let isLoading = $state(false);

    // Reaktivitas: Otomatis memuat data baru saat selectedDate berubah
    $effect(() => {
        if (selectedDate) {
            fetchReport(selectedDate);
        }
    });

    async function fetchReport(dateStr: string) {
        isLoading = true;
        try {
            const res = await apiFetch(`/reports/sales/daily?date=${dateStr}`);
            summary = res.data;
        } catch (err: any) {
            console.error("Gagal memuat laporan", err);
        } finally {
            isLoading = false;
        }
    }

    const formatRp = (n: number) => new Intl.NumberFormat('id-ID', { style: 'currency', currency: 'IDR', minimumFractionDigits: 0 }).format(n);
</script>

<div class="max-w-4xl space-y-8">
    <div class="flex flex-col sm:flex-row justify-between items-start sm:items-center gap-4 border-b border-gray-200 pb-6">
        <div>
            <h2 class="text-2xl font-black text-gray-800 uppercase tracking-tight">Laporan Keuangan</h2>
            <p class="text-gray-500 font-medium">Evaluasi performa penjualan toko kelontong Anda.</p>
        </div>
        
        <div class="bg-white p-2 rounded-2xl shadow-sm border border-gray-200 flex items-center gap-3">
            <span class="pl-2 font-bold text-gray-500 text-sm">Pilih Tanggal:</span>
            <input type="date" bind:value={selectedDate} max={today} 
                class="bg-gray-50 border-none rounded-xl px-4 py-2 text-gray-800 font-bold focus:ring-2 focus:ring-blue-500 cursor-pointer">
        </div>
    </div>

    <div class="grid grid-cols-1 md:grid-cols-2 gap-6 relative">
        {#if isLoading}
            <div class="absolute inset-0 bg-white/60 backdrop-blur-sm z-10 flex items-center justify-center rounded-3xl">
                <div class="animate-spin rounded-full h-10 w-10 border-t-4 border-blue-600"></div>
            </div>
        {/if}

        <div class="bg-gradient-to-br from-blue-600 to-blue-800 p-8 rounded-3xl shadow-lg text-white">
            <p class="text-blue-200 font-bold text-sm uppercase tracking-wider mb-2">Total Omzet Bersih</p>
            <p class="text-5xl font-black">{formatRp(summary.total_omzet)}</p>
            <p class="mt-4 text-blue-100 font-medium text-sm">Berdasarkan data penjualan tanggal {summary.date}</p>
        </div>

        <div class="bg-white p-8 rounded-3xl shadow-sm border border-gray-100 flex flex-col justify-center">
            <p class="text-gray-400 font-bold text-sm uppercase tracking-wider mb-2">Total Nota Dicetak</p>
            <div class="flex items-baseline gap-2">
                <p class="text-5xl font-black text-gray-800">{summary.total_transactions}</p>
                <p class="font-bold text-gray-500 uppercase">Transaksi</p>
            </div>
            <div class="w-full bg-gray-100 h-2 rounded-full mt-6">
                <div class="bg-green-500 h-2 rounded-full" style="width: {Math.min(100, summary.total_transactions)}%;"></div>
            </div>
        </div>
    </div>

    <div class="bg-blue-50 border border-blue-100 p-6 rounded-2xl flex gap-4 items-start">
        <div class="text-3xl">💡</div>
        <div>
            <h4 class="font-black text-blue-800">Catatan Sistem</h4>
            <p class="text-blue-700 font-medium mt-1">
                Data omzet di atas sudah dipotong otomatis dengan diskon yang diberikan kasir. Total omzet mencakup pembayaran tunai dan non-tunai (QRIS).
            </p>
        </div>
    </div>
</div>