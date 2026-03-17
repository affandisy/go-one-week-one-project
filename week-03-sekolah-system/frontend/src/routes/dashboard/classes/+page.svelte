<script lang="ts">
    import { onMount } from 'svelte';

    // ================= STATE MANAJEMEN =================
    // State untuk menampung data dari Backend
    let classes: any[] = [];
    let academicYears: any[] = []; // Untuk dropdown
    let studyPrograms: any[] = []; // Untuk dropdown
    
    // State untuk Form Tambah Kelas
    let newClassName = '';
    let selectedStudyProgram = '';
    let selectedAcademicYear = '';
    
    // State UI
    let isLoading = false;
    let isFetching = true;
    let message = '';
    let isError = false;

    // ================= LIFECYCLE =================
    onMount(async () => {
        // Saat halaman pertama kali dibuka, ambil semua data yang dibutuhkan
        await fetchAllData();
    });

    // ================= FUNGSI API (FETCH) =================
    async function fetchAllData() {
        isFetching = true;
        const token = localStorage.getItem('token');
        const headers = { 'Authorization': `Bearer ${token}` };

        try {
            // Kita bisa menjalankan beberapa request (GET) secara paralel (bersamaan)
            // agar loading halaman jauh lebih cepat!
            const [resClasses, resAY, resSP] = await Promise.all([
                fetch('http://localhost:3000/api/v1/classes', { headers }),
                fetch('http://localhost:3000/api/v1/academic-years', { headers }),
                // Anggap saja Anda sudah membuat endpoint GET /study-programs di backend
                fetch('http://localhost:3000/api/v1/study-programs', { headers }).catch(() => null) 
            ]);

            if (resClasses.ok) {
                const data = await resClasses.json();
                classes = data.data || [];
            }
            if (resAY.ok) {
                const data = await resAY.json();
                academicYears = data.data || [];
            }
            // Untuk Study Program, jika API belum ada, kita pakai data dummy sementara
            if (resSP && resSP.ok) {
                const data = await resSP.json();
                studyPrograms = data.data || [];
            } else {
                studyPrograms = [
                    { id: 'uuid-ipa', name: 'Ilmu Pengetahuan Alam (IPA)' },
                    { id: 'uuid-ips', name: 'Ilmu Pengetahuan Sosial (IPS)' }
                ];
            }
        } catch (err) {
            console.error("Gagal mengambil data master:", err);
        } finally {
            isFetching = false;
        }
    }

    // ================= FUNGSI API (POST) =================
    async function createClass() {
        if (!newClassName || !selectedStudyProgram || !selectedAcademicYear) {
            isError = true;
            message = "Nama Kelas, Program Studi, dan Tahun Ajaran wajib diisi!";
            return;
        }

        isLoading = true;
        message = '';
        const token = localStorage.getItem('token');

        try {
            const res = await fetch('http://localhost:3000/api/v1/classes', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                    'Authorization': `Bearer ${token}`
                },
                body: JSON.stringify({
                    name: newClassName,
                    study_program_id: selectedStudyProgram,
                    academic_year_id: selectedAcademicYear
                })
            });

            const data = await res.json();
            
            if (!res.ok) throw new Error(data.error);

            // Jika sukses:
            isError = false;
            message = "Kelas baru berhasil ditambahkan!";
            
            // 1. Bersihkan form
            newClassName = '';
            selectedStudyProgram = '';
            selectedAcademicYear = '';

            // 2. Ambil ulang data kelas dari backend agar tabel ter-update!
            await fetchAllData();

        } catch (err: any) {
            isError = true;
            message = err.message;
        } finally {
            isLoading = false;
        }
    }
</script>

<div class="p-8 max-w-6xl mx-auto">
    <div class="flex justify-between items-center mb-8">
        <div>
            <h2 class="text-3xl font-bold text-gray-800">Master Data: Kelas</h2>
            <p class="text-gray-500 mt-1">Kelola data rombongan belajar (rombel) sekolah.</p>
        </div>
    </div>

    {#if message}
        <div class="p-4 mb-6 rounded-lg shadow-sm border-l-4 {isError ? 'bg-red-50 text-red-700 border-red-500' : 'bg-green-50 text-green-700 border-green-500'} transition-all">
            {message}
        </div>
    {/if}

    <div class="grid grid-cols-1 lg:grid-cols-3 gap-8">
        
        <div class="lg:col-span-1">
            <div class="bg-white p-6 rounded-xl shadow-sm border border-gray-100 sticky top-6">
                <h3 class="text-lg font-bold text-gray-800 mb-4 border-b pb-2">Tambah Kelas Baru</h3>
                
                <div class="space-y-4">
                    <div>
                        <label for="class-name" class="block text-sm font-semibold text-gray-700 mb-1">Nama Kelas</label>
                        <input id="class-name" type="text" bind:value={newClassName} placeholder="Contoh: XII IPA 1" 
                            class="w-full border border-gray-300 rounded-lg p-2.5 focus:ring-2 focus:ring-blue-500 bg-gray-50 focus:bg-white transition" />
                    </div>

                    <div>
                        <label for="study-program" class="block text-sm font-semibold text-gray-700 mb-1">Program Studi</label>
                        <select id="study-program" bind:value={selectedStudyProgram} class="w-full border border-gray-300 rounded-lg p-2.5 focus:ring-2 focus:ring-blue-500 bg-gray-50 focus:bg-white transition">
                            <option value="">-- Pilih Program Studi --</option>
                            {#each studyPrograms as sp}
                                <option value={sp.id}>{sp.name}</option>
                            {/each}
                        </select>
                    </div>

                    <div>
                        <label for="academic-year" class="block text-sm font-semibold text-gray-700 mb-1">Tahun Ajaran</label>
                        <select id="academic-year" bind:value={selectedAcademicYear} class="w-full border border-gray-300 rounded-lg p-2.5 focus:ring-2 focus:ring-blue-500 bg-gray-50 focus:bg-white transition">
                            <option value="">-- Pilih Tahun Ajaran --</option>
                            {#each academicYears as ay}
                                <option value={ay.id}>{ay.name} {ay.is_active ? '(Aktif)' : ''}</option>
                            {/each}
                        </select>
                    </div>

                    <button on:click={createClass} disabled={isLoading} 
                        class="w-full mt-2 bg-blue-600 text-white font-bold px-4 py-2.5 rounded-lg shadow hover:bg-blue-700 disabled:bg-blue-300 transition-colors">
                        {isLoading ? 'Menyimpan...' : 'Simpan Kelas'}
                    </button>
                </div>
            </div>
        </div>

        <div class="lg:col-span-2">
            <div class="bg-white rounded-xl shadow-sm border border-gray-100 overflow-hidden">
                <div class="p-4 border-b border-gray-100 bg-gray-50 flex justify-between items-center">
                    <h3 class="font-bold text-gray-700">Daftar Kelas Terdaftar</h3>
                    <span class="bg-blue-100 text-blue-800 text-xs font-bold px-2.5 py-0.5 rounded-full">
                        Total: {classes.length} Kelas
                    </span>
                </div>

                {#if isFetching}
                    <div class="flex justify-center py-12">
                        <svg class="animate-spin h-8 w-8 text-blue-500" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24"><circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle><path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path></svg>
                    </div>
                {:else if classes.length === 0}
                    <div class="text-center py-12 text-gray-500">
                        Belum ada data kelas. Silakan tambahkan di form sebelah kiri.
                    </div>
                {:else}
                    <div class="overflow-x-auto">
                        <table class="min-w-full divide-y divide-gray-200">
                            <thead class="bg-white">
                                <tr>
                                    <th class="px-6 py-3 text-left text-xs font-bold text-gray-500 uppercase tracking-wider">Nama Kelas</th>
                                    <th class="px-6 py-3 text-left text-xs font-bold text-gray-500 uppercase tracking-wider">Program Studi</th>
                                    <th class="px-6 py-3 text-left text-xs font-bold text-gray-500 uppercase tracking-wider">Tahun Ajaran</th>
                                    <th class="px-6 py-3 text-left text-xs font-bold text-gray-500 uppercase tracking-wider">Wali Kelas</th>
                                </tr>
                            </thead>
                            <tbody class="divide-y divide-gray-100 bg-white">
                                {#each classes as cls}
                                    <tr class="hover:bg-gray-50 transition-colors">
                                        <td class="px-6 py-4 whitespace-nowrap font-bold text-gray-800">{cls.name}</td>
                                        <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-600">{cls.study_program?.name || '-'}</td>
                                        <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-600">
                                            {cls.academic_year?.name || '-'}
                                            {#if cls.academic_year?.is_active}
                                                <span class="ml-2 inline-flex items-center px-2 py-0.5 rounded text-xs font-medium bg-green-100 text-green-800">Aktif</span>
                                            {/if}
                                        </td>
                                        <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-500 italic">
                                            {cls.homeroom_teacher?.full_name || 'Belum ditugaskan'}
                                        </td>
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