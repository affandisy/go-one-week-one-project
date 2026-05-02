<script lang="ts">
    import { onMount } from 'svelte';
    import { apiFetch } from '$lib/api';

    // State untuk Pemilih Tanggal
    const today = new Date().toISOString().split('T')[0]; // Format YYYY-MM-DD
    let selectedDate = $state(today);
    
    // Sesuaikan State dengan API yang baru
    let summary = $state({ 
        total_transactions: 0, 
        total_omzet: 0, 
        total_profit: 0,
        cash_total: 0,
        qris_total: 0,
        date: '',
        products_sold: [] as any[]
    });
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
        } catch (err) {
            console.error(err);
        } finally {
            isLoading = false;
        }
    }

    const formatRp = (n: number) => new Intl.NumberFormat('id-ID', { style: 'currency', currency: 'IDR', minimumFractionDigits: 0 }).format(n);

    // FUNGSI EKSPOR KE CSV (EXCEL)
    function exportToExcel() {
        if (summary.products_sold.length === 0) {
            alert("Tidak ada data untuk diekspor pada tanggal ini.");
            return;
        }

        // Header CSV
        let csvContent = "data:text/csv;charset=utf-8,Nama Produk,Kuantitas Terjual,Total Penjualan (Rp)\n";

        // Isi Data CSV
        summary.products_sold.forEach(item => {
            const row = `"${item.name}",${item.quantity},${item.omzet}`;
            csvContent += row + "\n";
        });

        // Rekap Akhir di CSV
        csvContent += `\nTOTAL OMZET,,${summary.total_omzet}\n`;
        csvContent += `MARGIN KOTOR (PROFIT),,${summary.total_profit}\n`;

        // Proses Download Otomatis
        const encodedUri = encodeURI(csvContent);
        const link = document.createElement("a");
        link.setAttribute("href", encodedUri);
        link.setAttribute("download", `Laporan_POS_${summary.date}.csv`);
        document.body.appendChild(link);
        link.click();
        document.body.removeChild(link);
    }
</script>

<div class="max-w-4xl space-y-8">

    <!-- Header & Filter -->
    <div class="flex flex-col sm:flex-row justify-between items-start sm:items-center gap-4 border-b border-gray-200 pb-6">
        <div>
            <h2 class="text-2xl font-black text-gray-800 uppercase tracking-tight">Laporan Operasional</h2>
        </div>
        <div class="flex gap-3">
            <div class="bg-white p-2 rounded-xl shadow-sm border border-gray-200 flex items-center gap-2">
                <input type="date" bind:value={selectedDate} max={today} class="bg-gray-50 border-none rounded-lg px-4 py-2 font-bold cursor-pointer">
            </div>
            <!-- TOMBOL EKSPOR (PRD 5.5) -->
            <button onclick={exportToExcel} class="bg-green-600 hover:bg-green-700 text-white font-bold px-4 py-2 rounded-xl shadow transition-colors flex items-center gap-2">
                📥 Ekspor Excel
            </button>
        </div>
    </div>

    <!-- Papan Skor -->
    <div class="grid grid-cols-1 md:grid-cols-3 gap-6">
        <!-- Kotak Omzet -->
        <div class="bg-blue-600 p-6 rounded-3xl shadow-lg text-white">
            <p class="text-blue-200 font-bold text-sm uppercase">Total Omzet</p>
            <p class="text-4xl font-black mt-1">{formatRp(summary.total_omzet)}</p>
            <div class="mt-4 pt-4 border-t border-blue-500/50 flex justify-between text-sm font-bold">
                <span>Tunai: {formatRp(summary.cash_total)}</span>
                <span>QRIS: {formatRp(summary.qris_total)}</span>
            </div>
        </div>

        <!-- Kotak Keuntungan (PRD 5.5) -->
        <div class="bg-green-500 p-6 rounded-3xl shadow-lg text-white">
            <p class="text-green-100 font-bold text-sm uppercase">Margin Kotor (Profit)</p>
            <p class="text-4xl font-black mt-1">{formatRp(summary.total_profit)}</p>
            <p class="mt-4 text-green-100 text-sm font-medium">Berdasarkan Harga Jual - Modal</p>
        </div>

        <div class="bg-white p-6 rounded-3xl shadow-sm border border-gray-100 flex flex-col justify-center">
            <p class="text-gray-400 font-bold text-sm uppercase">Nota Dicetak</p>
            <p class="text-4xl font-black text-gray-800 mt-1">{summary.total_transactions} <span class="text-lg text-gray-400">Trx</span></p>
        </div>
    </div>

    <!-- Daftar Produk Terjual (PRD 5.5) -->
    <div class="bg-white rounded-3xl shadow-sm border border-gray-100 overflow-hidden">
        <div class="p-6 border-b border-gray-50">
            <h4 class="font-black text-gray-800 uppercase">Rincian Barang Terjual</h4>
        </div>
        {#if summary.products_sold.length === 0}
            <div class="p-8 text-center text-gray-400 font-bold">Belum ada transaksi pada tanggal ini.</div>
        {:else}
            <div class="overflow-x-auto">
                <table class="w-full text-left">
                    <thead class="bg-gray-50 text-gray-400 text-xs font-black uppercase">
                        <tr>
                            <th class="px-6 py-4">Nama Produk</th>
                            <th class="px-6 py-4">Kuantitas Terjual</th>
                            <th class="px-6 py-4 text-right">Subtotal Omzet</th>
                        </tr>
                    </thead>
                    <tbody class="divide-y divide-gray-50">
                        {#each summary.products_sold as item}
                            <tr class="hover:bg-gray-50/50">
                                <td class="px-6 py-4 font-bold text-gray-700">{item.name}</td>
                                <td class="px-6 py-4 font-black text-blue-600">{item.quantity}</td>
                                <td class="px-6 py-4 text-right font-bold text-gray-800">{formatRp(item.omzet)}</td>
                            </tr>
                        {/each}
                    </tbody>
                </table>
            </div>
        {/if}
    </div>
    
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