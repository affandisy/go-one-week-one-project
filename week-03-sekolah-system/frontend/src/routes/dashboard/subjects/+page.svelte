<script lang="ts">
    import { onMount } from 'svelte';

    // State Data
    let subjects: any[] = [];
    let studyPrograms: any[] = [];
    
    // State Form
    let newCode = '';
    let newName = '';
    let selectedStudyProgram = '';
    
    // State UI
    let isLoading = false;
    let isFetching = true;
    let message = '';
    let isError = false;

    onMount(async () => {
        await fetchAllData();
    });

    async function fetchAllData() {
        isFetching = true;
        const token = localStorage.getItem('token');
        const headers = { 'Authorization': `Bearer ${token}` };

        try {
            const [resSubjects, resSP] = await Promise.all([
                fetch('http://localhost:3000/api/v1/subjects', { headers }),
                fetch('http://localhost:3000/api/v1/study-programs', { headers }).catch(() => null)
            ]);

            if (resSubjects.ok) {
                const data = await resSubjects.json();
                subjects = data.data || [];
            }
            if (resSP && resSP.ok) {
                const data = await resSP.json();
                studyPrograms = data.data || [];
            } else {
                // Dummy fallback jika endpoint belum ada
                studyPrograms = [
                    { id: 'uuid-ipa', name: 'Ilmu Pengetahuan Alam (IPA)' },
                    { id: 'uuid-ips', name: 'Ilmu Pengetahuan Sosial (IPS)' }
                ];
            }
        } catch (err) {
            console.error(err);
        } finally {
            isFetching = false;
        }
    }

    async function createSubject() {
        if (!newCode || !newName) {
            isError = true;
            message = "Kode dan Nama Mata Pelajaran wajib diisi!";
            return;
        }

        isLoading = true;
        message = '';
        const token = localStorage.getItem('token');

        try {
            const res = await fetch('http://localhost:3000/api/v1/subjects', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                    'Authorization': `Bearer ${token}`
                },
                body: JSON.stringify({
                    code: newCode,
                    name: newName,
                    // Kirim string kosong jika mapel umum (tidak ada prodi)
                    study_program_id: selectedStudyProgram === '' ? undefined : selectedStudyProgram
                })
            });

            const data = await res.json();
            if (!res.ok) throw new Error(data.error);

            isError = false;
            message = "Mata pelajaran berhasil ditambahkan!";
            
            newCode = '';
            newName = '';
            selectedStudyProgram = '';

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
    <div class="mb-8">
        <h2 class="text-3xl font-bold text-gray-800">Master Data: Mata Pelajaran</h2>
        <p class="text-gray-500 mt-1">Kelola kurikulum dan mata pelajaran sekolah.</p>
    </div>

    {#if message}
        <div class="p-4 mb-6 rounded-lg shadow-sm border-l-4 {isError ? 'bg-red-50 text-red-700 border-red-500' : 'bg-green-50 text-green-700 border-green-500'}">
            {message}
        </div>
    {/if}

    <div class="grid grid-cols-1 lg:grid-cols-3 gap-8">
        <div class="lg:col-span-1">
            <div class="bg-white p-6 rounded-xl shadow-sm border border-gray-100 sticky top-6">
                <h3 class="text-lg font-bold text-gray-800 mb-4 border-b pb-2">Tambah Mapel</h3>
                
                <div class="space-y-4">
                    <div>
                        <label for="subject-code" class="block text-sm font-semibold text-gray-700 mb-1">Kode Mapel</label>
                        <input id="subject-code" type="text" bind:value={newCode} placeholder="Contoh: MAT-LJT" 
                            class="w-full border border-gray-300 rounded-lg p-2.5 focus:ring-2 focus:ring-blue-500 bg-gray-50" />
                    </div>
                    <div>
                        <label for="subject-name" class="block text-sm font-semibold text-gray-700 mb-1">Nama Mata Pelajaran</label>
                        <input id="subject-name" type="text" bind:value={newName} placeholder="Contoh: Matematika Lanjut" 
                            class="w-full border border-gray-300 rounded-lg p-2.5 focus:ring-2 focus:ring-blue-500 bg-gray-50" />
                    </div>
                    <div>
                        <label for="subject-study-program" class="block text-sm font-semibold text-gray-700 mb-1">Program Studi (Opsional)</label>
                        <select id="subject-study-program" bind:value={selectedStudyProgram} class="w-full border border-gray-300 rounded-lg p-2.5 focus:ring-2 focus:ring-blue-500 bg-gray-50">
                            <option value="">-- Umum (Semua Jurusan) --</option>
                            {#each studyPrograms as sp}
                                <option value={sp.id}>{sp.name}</option>
                            {/each}
                        </select>
                        <p class="text-xs text-gray-400 mt-1">Kosongkan jika mapel wajib seperti Agama/PKn.</p>
                    </div>

                    <button on:click={createSubject} disabled={isLoading} 
                        class="w-full mt-2 bg-indigo-600 text-white font-bold px-4 py-2.5 rounded-lg shadow hover:bg-indigo-700 disabled:bg-indigo-300">
                        {isLoading ? 'Menyimpan...' : 'Simpan Mapel'}
                    </button>
                </div>
            </div>
        </div>

        <div class="lg:col-span-2">
            <div class="bg-white rounded-xl shadow-sm border border-gray-100 overflow-hidden">
                <div class="p-4 bg-gray-50 flex justify-between items-center border-b border-gray-100">
                    <h3 class="font-bold text-gray-700">Daftar Mata Pelajaran</h3>
                    <span class="bg-indigo-100 text-indigo-800 text-xs font-bold px-2.5 py-0.5 rounded-full">Total: {subjects.length} Mapel</span>
                </div>
                
                {#if isFetching}
                    <div class="text-center py-12 text-gray-400">Loading data...</div>
                {:else if subjects.length === 0}
                    <div class="text-center py-12 text-gray-500">Belum ada mata pelajaran terdaftar.</div>
                {:else}
                    <table class="min-w-full divide-y divide-gray-200">
                        <thead class="bg-white">
                            <tr>
                                <th class="px-6 py-3 text-left text-xs font-bold text-gray-500 uppercase tracking-wider">Kode</th>
                                <th class="px-6 py-3 text-left text-xs font-bold text-gray-500 uppercase tracking-wider">Nama Mapel</th>
                                <th class="px-6 py-3 text-left text-xs font-bold text-gray-500 uppercase tracking-wider">Program Studi</th>
                            </tr>
                        </thead>
                        <tbody class="divide-y divide-gray-100 bg-white">
                            {#each subjects as sub}
                                <tr class="hover:bg-gray-50">
                                    <td class="px-6 py-4 whitespace-nowrap font-mono text-sm text-indigo-600 font-bold">{sub.code}</td>
                                    <td class="px-6 py-4 whitespace-nowrap text-sm font-bold text-gray-800">{sub.name}</td>
                                    <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-600">
                                        {#if sub.study_program}
                                            <span class="bg-blue-50 text-blue-600 border border-blue-200 px-2 py-1 rounded text-xs">{sub.study_program.name}</span>
                                        {:else}
                                            <span class="bg-gray-100 text-gray-600 px-2 py-1 rounded text-xs">Umum</span>
                                        {/if}
                                    </td>
                                </tr>
                            {/each}
                        </tbody>
                    </table>
                {/if}
            </div>
        </div>
    </div>
</div>