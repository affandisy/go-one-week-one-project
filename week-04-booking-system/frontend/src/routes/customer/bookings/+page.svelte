<script lang="ts">
    import { onMount, onDestroy } from 'svelte';
    import { apiFetch } from '$lib/api';
    import { goto } from '$app/navigation';

    // State Data
    let bookings = $state<any[]>([]);
    
    // UI State
    let isLoading = $state(true);
    let isActionLoading = $state(false);
    let message = $state('');
    let isError = $state(false);

    // State untuk Live Countdown
    let now = $state(new Date().getTime());
    let timerId: ReturnType<typeof setInterval>;

    onMount(async () => {
        // Cek autentikasi
        if (!localStorage.getItem('token')) {
            goto('/login');
            return;
        }

        await fetchMyBookings();

        // Jalankan detak jam setiap 1 detik untuk efek hitung mundur
        timerId = setInterval(() => {
            now = new Date().getTime();
        }, 1000);
    });

    onDestroy(() => {
        if (timerId) clearInterval(timerId); // Bersihkan memori saat pindah halaman
    });

    async function fetchMyBookings() {
        isLoading = true;
        try {
            const res = await apiFetch('/bookings/me');
            bookings = res.data || [];
        } catch (err: any) {
            isError = true;
            message = err.message || "Gagal mengambil data booking.";
        } finally {
            isLoading = false;
        }
    }

    async function cancelBooking(id: string) {
        if (!confirm('Yakin ingin membatalkan jadwal ini?')) return;
        
        isActionLoading = true;
        try {
            await apiFetch(`/bookings/${id}/cancel`, { method: 'PUT' });
            // Refresh data setelah sukses membatalkan
            await fetchMyBookings();
            alert('Booking berhasil dibatalkan.');
        } catch (err: any) {
            alert(err.message || 'Gagal membatalkan booking.');
        } finally {
            isActionLoading = false;
        }
    }

    // Fungsi Jembatan menuju Phase 3 (Pembayaran Simulasi)
    async function payBooking(id: string) {
        if (!confirm('Lanjutkan pembayaran sebesar tagihan? (Simulasi)')) return;

        isActionLoading = true;
        try {
            await apiFetch(`/bookings/${id}/pay`, { method: 'POST' });
            alert('Pembayaran berhasil! Status berubah menjadi PAID.');
            await fetchMyBookings();
        } catch (err: any) {
            alert(err.message || 'Gagal memproses pembayaran.');
        } finally {
            isActionLoading = false;
        }
    }

    async function downloadReceipt(id: string) {
        isActionLoading = true;
        try {
            const res = await apiFetch(`/bookings/${id}/receipt`, { method: 'GET' });
            
            // Membuka URL PDF di tab baru
            if (res.pdf_url) {
                window.open(res.pdf_url, '_blank');
            } else {
                alert('URL PDF tidak ditemukan.');
            }
        } catch (err: any) {
            alert(err.message || 'Gagal memuat e-receipt.');
        } finally {
            isActionLoading = false;
        }
    }

    // Helper untuk mengubah milidetik menjadi format Menit:Detik
    function formatTimeLeft(lockExpiry: string): string {
        const expiryTime = new Date(lockExpiry).getTime();
        const diff = expiryTime - now;

        if (diff <= 0) return "Waktu Habis";

        const minutes = Math.floor((diff % (1000 * 60 * 60)) / (1000 * 60));
        const seconds = Math.floor((diff % (1000 * 60)) / 1000);

        return `${minutes}m ${seconds}s`;
    }
</script>

<div class="max-w-5xl mx-auto py-12 px-4 sm:px-6 lg:px-8">
    <div class="flex justify-between items-end mb-8 border-b pb-4">
        <div>
            <h2 class="text-3xl font-extrabold text-gray-900">Tiket Saya</h2>
            <p class="mt-2 text-sm text-gray-600">Kelola jadwal main dan selesaikan pembayaran Anda di sini.</p>
        </div>
        <a href="/" class="text-blue-600 font-bold hover:text-blue-800 transition">
            + Booking Lapangan Baru
        </a>
    </div>

    {#if message}
        <div class="p-4 mb-6 rounded-lg shadow-sm border-l-4 {isError ? 'bg-red-50 text-red-700 border-red-500' : 'bg-green-50 text-green-700 border-green-500'}">
            {message}
        </div>
    {/if}

    {#if isLoading}
        <div class="text-center py-20 text-gray-500 font-medium animate-pulse">
            Memuat tiket Anda...
        </div>
    {:else if bookings.length === 0}
        <div class="text-center py-20 bg-white rounded-xl shadow-sm border border-gray-100">
            <svg class="mx-auto h-12 w-12 text-gray-300" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 7V3m8 4V3m-9 8h10M5 21h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v12a2 2 0 002 2z" />
            </svg>
            <h3 class="mt-2 text-sm font-bold text-gray-900">Belum Ada Tiket</h3>
            <p class="mt-1 text-sm text-gray-500">Anda belum pernah melakukan booking lapangan.</p>
        </div>
    {:else}
        <div class="space-y-4">
            {#each bookings as b}
                <div class="bg-white p-6 rounded-xl shadow-sm border border-gray-100 flex flex-col md:flex-row justify-between items-start md:items-center gap-6 hover:shadow-md transition-shadow">
                    
                    <div class="flex-1">
                        <div class="flex items-center gap-3 mb-2">
                            <span class="font-mono text-xs font-bold text-gray-500 bg-gray-100 px-2 py-1 rounded">
                                {b.booking_code}
                            </span>
                            
                            {#if b.status === 'locked'}
                                <span class="bg-orange-100 text-orange-800 text-xs font-black px-2.5 py-1 rounded-full animate-pulse">
                                    MENUNGGU PEMBAYARAN
                                </span>
                            {:else if b.status === 'paid'}
                               <button 
                                onclick={() => downloadReceipt(b.id)} 
                                disabled={isActionLoading}
                                class="w-full bg-green-600 text-white font-bold py-2.5 px-6 rounded-lg shadow hover:bg-green-700 flex items-center justify-center gap-2 transition-colors disabled:bg-green-400">
                                
                                {#if isActionLoading}
                                    <span class="animate-pulse">Mencetak...</span>
                                {:else}
                                    <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 16v1a3 3 0 003 3h10a3 3 0 003-3v-1m-4-4l-4 4m0 0l-4-4m4 4V4"></path></svg>
                                    E-Receipt (PDF)
                                {/if}
                            </button>
                            {:else if b.status === 'cancelled'}
                                <span class="bg-red-100 text-red-800 text-xs font-black px-2.5 py-1 rounded-full">
                                    DIBATALKAN
                                </span>
                            {:else}
                                <span class="bg-gray-100 text-gray-800 text-xs font-black px-2.5 py-1 rounded-full">
                                    EXPIRED
                                </span>
                            {/if}
                        </div>

                        <h3 class="text-xl font-bold text-gray-900">{b.court?.name}</h3>
                        <p class="text-gray-600 font-medium mt-1">
                            {new Date(b.booking_date).toLocaleDateString('id-ID', { weekday: 'long', year: 'numeric', month: 'long', day: 'numeric' })}
                        </p>
                        <p class="text-blue-600 font-black text-lg mt-1">
                            {b.start_time} - {b.end_time}
                        </p>
                    </div>

                    <div class="w-full md:w-auto flex flex-col items-end gap-3 border-t md:border-t-0 pt-4 md:pt-0">
                        <div class="text-right w-full">
                            <p class="text-sm text-gray-500 font-bold mb-1">Total Harga</p>
                            <p class="text-2xl font-black text-gray-900">
                                Rp {b.total_price.toLocaleString('id-ID')}
                            </p>
                        </div>

                        {#if b.status === 'locked'}
                            <div class="w-full bg-orange-50 border border-orange-200 rounded-lg p-3 text-center mb-2">
                                <p class="text-xs text-orange-800 font-bold mb-1">Sisa Waktu Bayar:</p>
                                <p class="text-lg font-black text-orange-600 font-mono">
                                    {formatTimeLeft(b.lock_expiry)}
                                </p>
                            </div>
                            
                            <div class="flex w-full gap-2">
                                <button onclick={() => cancelBooking(b.id)} disabled={isActionLoading}
                                    class="flex-1 bg-white border border-gray-300 text-gray-700 font-bold py-2 rounded-lg hover:bg-gray-50 disabled:opacity-50">
                                    Batal
                                </button>
                                <button onclick={() => payBooking(b.id)} disabled={isActionLoading || formatTimeLeft(b.lock_expiry) === 'Waktu Habis'}
                                    class="flex-1 bg-blue-600 text-white font-bold py-2 rounded-lg shadow hover:bg-blue-700 disabled:bg-gray-400">
                                    Bayar
                                </button>
                            </div>

                        {:else if b.status === 'paid'}
                            <button class="w-full bg-green-600 text-white font-bold py-2.5 px-6 rounded-lg shadow hover:bg-green-700 flex items-center justify-center gap-2">
                                <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 16v1a3 3 0 003 3h10a3 3 0 003-3v-1m-4-4l-4 4m0 0l-4-4m4 4V4"></path></svg>
                                E-Receipt (PDF)
                            </button>
                        {/if}
                    </div>

                </div>
            {/each}
        </div>
    {/if}
</div>