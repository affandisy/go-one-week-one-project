<script lang="ts">
    import { onMount } from 'svelte';

    // State untuk Form
    let title = '';
    let content = '';
    let targetRoles: string[] = []; // Array untuk menampung role yang dicentang
    
    // State untuk UI
    let announcements: any[] = [];
    let isLoading = false;
    let isFetching = true;
    let message = '';
    let isError = false;

    // Daftar role yang bisa dipilih
    const availableRoles = ['GURU', 'MURID', 'TU', 'WALI_MURID'];

    onMount(() => {
        fetchAnnouncements();
    });

    // Fungsi mengambil daftar pengumuman dari Golang
    async function fetchAnnouncements() {
        const token = localStorage.getItem('token');
        try {
            const res = await fetch('http://localhost:3000/api/v1/announcements', {
                headers: { 'Authorization': `Bearer ${token}` }
            });
            const data = await res.json();
            if (res.ok) {
                announcements = data.data || [];
            }
        } catch (err) {
            console.error("Gagal mengambil pengumuman", err);
        } finally {
            isFetching = false;
        }
    }

    // Fungsi untuk Toggle Checkbox Role
    function toggleRole(role: string) {
        if (targetRoles.includes(role)) {
            targetRoles = targetRoles.filter(r => r !== role);
        } else {
            targetRoles = [...targetRoles, role];
        }
    }

    // Fungsi membuat pengumuman baru
    async function submitAnnouncement() {
        if (!title || !content || targetRoles.length === 0) {
            isError = true;
            message = "Judul, Konten, dan minimal 1 Target Role wajib diisi!";
            return;
        }

        isLoading = true;
        message = '';
        const token = localStorage.getItem('token');

        try {
            const res = await fetch('http://localhost:3000/api/v1/announcements', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                    'Authorization': `Bearer ${token}`
                },
                body: JSON.stringify({
                    title: title,
                    content: content,
                    target_roles: targetRoles
                })
            });

            const data = await res.json();
            
            if (!res.ok) throw new Error(data.error);

            isError = false;
            message = data.message;
            
            // Bersihkan form
            title = '';
            content = '';
            targetRoles = [];

            // Refresh daftar pengumuman agar yang baru langsung muncul!
            await fetchAnnouncements();

        } catch (err: any) {
            isError = true;
            message = err.message;
        } finally {
            isLoading = false;
        }
    }

    // Helper untuk format tanggal dari Golang (ISO string)
    function formatDate(dateString: string) {
        const options: Intl.DateTimeFormatOptions = { 
            year: 'numeric', month: 'long', day: 'numeric', 
            hour: '2-digit', minute: '2-digit' 
        };
        return new Date(dateString).toLocaleDateString('id-ID', options);
    }
</script>

<div class="p-8 max-w-5xl mx-auto">
    <h2 class="text-3xl font-bold text-gray-800 mb-8">Papan Pengumuman Sekolah</h2>

    {#if message}
        <div class="p-4 mb-6 rounded shadow-sm {isError ? 'bg-red-50 text-red-700 border-l-4 border-red-500' : 'bg-green-50 text-green-700 border-l-4 border-green-500'}">
            {message}
        </div>
    {/if}

    <div class="grid grid-cols-1 lg:grid-cols-3 gap-8">
        
        <div class="lg:col-span-1">
            <div class="bg-white p-6 rounded-xl shadow-sm border border-gray-100">
                <h3 class="text-lg font-bold text-gray-700 mb-4 flex items-center gap-2">
                    <svg class="w-5 h-5 text-blue-500" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M11 5.882V19.24a1.76 1.76 0 01-3.417.592l-2.147-6.15M18 13a3 3 0 100-6M5.436 13.683A4.001 4.001 0 017 6h1.832c4.1 0 7.625-1.234 9.168-3v14c-1.543-1.766-5.067-3-9.168-3H7a3.988 3.988 0 01-1.564-.317z"></path></svg>
                    Buat Pengumuman
                </h3>
                
                <div class="space-y-4">
                    <div>
                        <label for="announcement-title" class="block text-sm font-semibold text-gray-600 mb-1">Judul</label>
                        <input id="announcement-title" type="text" bind:value={title} placeholder="Contoh: Libur Semester Genap" 
                            class="w-full border-gray-300 rounded-lg p-2.5 bg-gray-50 focus:bg-white focus:ring-2 focus:ring-blue-500 transition" />
                    </div>

                    <div>
                        <label for="announcement-content" class="block text-sm font-semibold text-gray-600 mb-1">Konten / Isi</label>
                        <textarea id="announcement-content" bind:value={content} rows="4" placeholder="Tuliskan detail pengumuman di sini..." 
                            class="w-full border-gray-300 rounded-lg p-2.5 bg-gray-50 focus:bg-white focus:ring-2 focus:ring-blue-500 transition"></textarea>
                    </div>

                    <div>
                        <p class="block text-sm font-semibold text-gray-600 mb-2">Target Penerima</p>
                        <div class="flex flex-wrap gap-2">
                            {#each availableRoles as role}
                                <button type="button" 
                                    on:click={() => toggleRole(role)}
                                    class="px-3 py-1.5 text-sm rounded-full border transition-all duration-200 
                                    {targetRoles.includes(role) ? 'bg-blue-100 border-blue-400 text-blue-700 font-medium' : 'bg-white border-gray-200 text-gray-500 hover:bg-gray-50'}">
                                    {role}
                                </button>
                            {/each}
                        </div>
                    </div>

                    <button on:click={submitAnnouncement} disabled={isLoading} 
                        class="w-full mt-4 bg-blue-600 text-white font-semibold px-4 py-2.5 rounded-lg shadow-md hover:bg-blue-700 focus:ring-4 focus:ring-blue-300 disabled:opacity-50 transition">
                        {isLoading ? 'Menyiarkan...' : 'Siarkan Pengumuman'}
                    </button>
                </div>
            </div>
        </div>

        <div class="lg:col-span-2 space-y-4">
            {#if isFetching}
                <div class="flex justify-center items-center py-12 text-gray-400">
                    <svg class="animate-spin h-8 w-8 text-blue-500" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
                        <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
                        <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
                    </svg>
                </div>
            {:else if announcements.length === 0}
                <div class="bg-white p-8 rounded-xl shadow-sm border border-gray-100 text-center text-gray-500">
                    Belum ada pengumuman sekolah saat ini.
                </div>
            {:else}
                {#each announcements as ann}
                    <div class="bg-white p-6 rounded-xl shadow-sm border border-gray-100 hover:shadow-md transition">
                        <div class="flex justify-between items-start mb-2">
                            <h4 class="text-xl font-bold text-gray-800">{ann.title}</h4>
                            <span class="text-xs text-gray-400 bg-gray-50 px-2 py-1 rounded-md border">
                                {formatDate(ann.created_at)}
                            </span>
                        </div>
                        
                        <p class="text-gray-600 whitespace-pre-line mb-4 leading-relaxed">
                            {ann.content}
                        </p>
                        
                        <div class="flex items-center justify-between pt-4 border-t border-gray-50">
                            <div class="flex gap-2 items-center">
                                <span class="text-xs text-gray-500 font-medium">Untuk:</span>
                                {#each JSON.parse(ann.target_roles || '[]') as role}
                                    <span class="text-[10px] font-bold tracking-wider text-indigo-600 bg-indigo-50 px-2 py-0.5 rounded uppercase">
                                        {role}
                                    </span>
                                {/each}
                            </div>
                            <div class="text-xs text-gray-500">
                                Oleh: <span class="font-semibold text-gray-700">{ann.created_by?.full_name || 'Admin'}</span>
                            </div>
                        </div>
                    </div>
                {/each}
            {/if}
        </div>

    </div>
</div>