<script lang="ts">
    import { onMount } from 'svelte';
    import { apiFetch } from '$lib/api';
    import { goto } from '$app/navigation';

    // State Data
    let courts = $state<any[]>([]);
    let slots = $state<any[]>([]);
    
    // State Filter (Default hari ini)
    let selectedCourtId = $state('');
    let selectedDate = $state(new Date().toISOString().split('T')[0]);
    
    // UI State
    let isFetchingSlots = $state(false);
    let isBooking = $state(false);
    let message = $state('');
    let isError = $state(false);

    onMount(async () => {
        await fetchCourts();
    });

    // Rune $effect: Otomatis load jadwal jika Lapangan & Tanggal terisi
    $effect(() => {
        if (selectedCourtId && selectedDate) {
            fetchAvailability();
        } else {
            slots = [];
        }
    });

    async function fetchCourts() {
        try {
            // Catatan: Pastikan endpoint GET /courts di Golang Anda 
            // sudah dipindah ke rute publik agar bisa diakses tanpa login.
            const res = await apiFetch('/courts'); 
            courts = res.data || [];
            if (courts.length > 0) {
                selectedCourtId = courts[0].id;
            }
        } catch (err) {
            console.error("Gagal mengambil data lapangan:", err);
        }
    }

    async function fetchAvailability() {
        isFetchingSlots = true;
        try {
            const res = await apiFetch(`/availability?court_id=${selectedCourtId}&date=${selectedDate}`);
            slots = res.data || [];
        } catch (err) {
            console.error("Gagal mengambil jadwal:", err);
            slots = [];
        } finally {
            isFetchingSlots = false;
        }
    }

    async function bookSlot(startTime: string, endTime: string) {
        // Cek apakah user sudah login (punya token)
        const token = localStorage.getItem('token');
        if (!token) {
            alert('Silakan login terlebih dahulu untuk melakukan booking.');
            goto('/login');
            return;
        }

        if (!confirm(`Kunci slot jam ${startTime} - ${endTime}? Anda memiliki waktu 10 menit untuk membayar setelah ini.`)) {
            return;
        }

        isBooking = true;
        message = '';

        try {
            const res = await apiFetch('/bookings', {
                method: 'POST',
                data: {
                    court_id: selectedCourtId,
                    booking_date: selectedDate,
                    start_time: startTime,
                    end_time: endTime
                }
            });

            isError = false;
            message = res.message || 'Slot berhasil dikunci!';
            
            // Redirect ke halaman riwayat booking customer untuk bayar
            setTimeout(() => {
                goto('/customer/bookings');
            }, 1500);

        } catch (err: any) {
            isError = true;
            message = err.message;
            // Refresh ketersediaan karena mungkin slot baru saja diambil orang lain (Double Booking protection)
            fetchAvailability();
        } finally {
            isBooking = false;
        }
    }
</script>

<div class="min-h-screen bg-gray-50 pb-12">
    <div class="bg-blue-600 text-white pt-16 pb-24 px-4 sm:px-6 lg:px-8 text-center">
        <h1 class="text-4xl font-extrabold tracking-tight sm:text-5xl lg:text-6xl">
            Main Padel Kapan Saja
        </h1>
        <p class="mt-4 max-w-2xl mx-auto text-xl text-blue-100">
            Cek ketersediaan secara real-time dan booking lapangan favorit Anda dalam hitungan detik.
        </p>
    </div>

    <div class="max-w-5xl mx-auto px-4 sm:px-6 lg:px-8 -mt-12">
        <div class="bg-white rounded-xl shadow-lg p-6 border border-gray-100">
            
            <div class="flex flex-col md:flex-row gap-4 mb-8">
                <div class="flex-1">
                    <label for="selectedCourtId" class="block text-sm font-bold text-gray-700 mb-2">Pilih Lapangan</label>
                    <select id="selectedCourtId" bind:value={selectedCourtId} class="w-full border border-gray-300 rounded-lg p-3 focus:ring-2 focus:ring-blue-500 bg-gray-50 font-medium">
                        <option value="">-- Pilih Lapangan --</option>
                        {#each courts as court}
                            <option value={court.id}>{court.name} ({court.type})</option>
                        {/each}
                    </select>
                </div>
                <div class="flex-1">
                    <label for="selectedDate" class="block text-sm font-bold text-gray-700 mb-2">Tanggal Main</label>
                    <input id="selectedDate" type="date" bind:value={selectedDate} min={new Date().toISOString().split('T')[0]} 
                        class="w-full border border-gray-300 rounded-lg p-3 focus:ring-2 focus:ring-blue-500 bg-gray-50 font-medium" />
                </div>
            </div>

            {#if message}
                <div class="p-4 mb-6 rounded-lg shadow-sm border-l-4 font-medium {isError ? 'bg-red-50 text-red-700 border-red-500' : 'bg-green-50 text-green-700 border-green-500'}">
                    {message}
                </div>
            {/if}

            <div>
                <h3 class="text-lg font-bold text-gray-800 mb-4 border-b pb-2">Ketersediaan Jam</h3>
                
                {#if isFetchingSlots}
                    <div class="text-center py-12 text-gray-500 font-medium animate-pulse">
                        Mencari jadwal kosong...
                    </div>
                {:else if !selectedCourtId}
                    <div class="text-center py-12 text-gray-500">
                        Silakan pilih lapangan terlebih dahulu.
                    </div>
                {:else if slots.length === 0}
                    <div class="text-center py-12 text-gray-500">
                        Jadwal tidak tersedia untuk tanggal ini.
                    </div>
                {:else}
                    <div class="grid grid-cols-2 md:grid-cols-3 lg:grid-cols-4 gap-4">
                        {#each slots as slot}
                            <button 
                                disabled={slot.status !== 'available' || isBooking}
                                onclick={() => bookSlot(slot.start, slot.end)}
                                class="relative flex flex-col items-center justify-center p-4 rounded-xl border-2 transition-all 
                                    {slot.status === 'available' 
                                        ? 'border-blue-100 bg-white hover:border-blue-500 hover:shadow-md cursor-pointer group' 
                                        : 'border-gray-100 bg-gray-50 opacity-60 cursor-not-allowed'}">
                                
                                <span class="text-lg font-black text-gray-800">{slot.start} - {slot.end}</span>
                                
                                {#if slot.status === 'available'}
                                    <span class="mt-1 text-sm font-bold text-blue-600 group-hover:text-blue-700">
                                        Rp {slot.price.toLocaleString('id-ID')}
                                    </span>
                                    <span class="mt-2 px-2 py-0.5 bg-green-100 text-green-700 text-[10px] uppercase font-black rounded">
                                        Tersedia
                                    </span>
                                {:else}
                                    <span class="mt-1 text-sm font-medium text-gray-500 line-through">
                                        Rp {slot.price.toLocaleString('id-ID')}
                                    </span>
                                    <span class="mt-2 px-2 py-0.5 bg-red-100 text-red-700 text-[10px] uppercase font-black rounded">
                                        Sudah Dipesan
                                    </span>
                                {/if}
                            </button>
                        {/each}
                    </div>
                {/if}
            </div>

        </div>
    </div>
</div>