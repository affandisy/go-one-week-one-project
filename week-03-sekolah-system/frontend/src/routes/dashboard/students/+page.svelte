<script lang="ts">
    import { onMount } from 'svelte';

    // State Data
    let students: any[] = [];
    let classes: any[] = [];
    let academicYears: any[] = [];
    let userCandidates: any[] = []; // Daftar akun user yang belum jadi murid
    
    // State Form
    let selectedUserId = '';
    let newNis = '';
    let selectedClassId = '';
    let selectedEnrollmentYearId = '';
    
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
            const [resStudents, resClasses, resAY, resUsers] = await Promise.all([
                fetch('http://localhost:3000/api/v1/students', { headers }),
                fetch('http://localhost:3000/api/v1/classes', { headers }),
                fetch('http://localhost:3000/api/v1/academic-years', { headers }),
                // Asumsi: Endpoint ini memfilter user dengan role MURID yang belum didaftarkan
                fetch('http://localhost:3000/api/v1/users?role=MURID', { headers }).catch(() => null)
            ]);

            if (resStudents.ok) { const data = await resStudents.json(); students = data.data || []; }
            if (resClasses.ok) { const data = await resClasses.json(); classes = data.data || []; }
            if (resAY.ok) { const data = await resAY.json(); academicYears = data.data || []; }
            
            // Fallback Dummy User jika API /users belum selesai
            if (resUsers && resUsers.ok) {
                const data = await resUsers.json(); userCandidates = data.data || [];
            } else {
                userCandidates = [
                    { id: 'uuid-user-1', full_name: 'Budi Santoso', email: 'budi@siswa.com' },
                    { id: 'uuid-user-2', full_name: 'Siti Aminah', email: 'siti@siswa.com' }
                ];
            }
        } catch (err) {
            console.error(err);
        } finally {
            isFetching = false;
        }
    }

    async function createStudent() {
        if (!selectedUserId || !newNis || !selectedClassId || !selectedEnrollmentYearId) {
            isError = true;
            message = "Semua kolom wajib diisi!";
            return;
        }

        isLoading = true;
        message = '';
        const token = localStorage.getItem('token');

        try {
            const res = await fetch('http://localhost:3000/api/v1/students', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                    'Authorization': `Bearer ${token}`
                },
                body: JSON.stringify({
                    user_id: selectedUserId,
                    nis: newNis,
                    class_id: selectedClassId,
                    enrollment_year_id: selectedEnrollmentYearId
                })
            });

            const data = await res.json();
            if (!res.ok) throw new Error(data.error);

            isError = false;
            message = "Siswa berhasil didaftarkan ke kelas!";
            
            // Reset Form
            selectedUserId = '';
            newNis = '';
            selectedClassId = '';
            selectedEnrollmentYearId = '';

            await fetchAllData();
        } catch (err: any) {
            isError = true;
            message = err.message;
        } finally {
            isLoading = false;
        }
    }
</script>

<div class="p-8 max-w-7xl mx-auto">
    <div class="mb-8">
        <h2 class="text-3xl font-bold text-gray-800">Master Data: Siswa</h2>
        <p class="text-gray-500 mt-1">Daftarkan akun User menjadi Siswa resmi dan masukkan ke dalam kelas.</p>
    </div>

    {#if message}
        <div class="p-4 mb-6 rounded-lg shadow-sm border-l-4 {isError ? 'bg-red-50 text-red-700 border-red-500' : 'bg-green-50 text-green-700 border-green-500'}">
            {message}
        </div>
    {/if}

    <div class="grid grid-cols-1 xl:grid-cols-3 gap-8">
        <div class="xl:col-span-1">
            <div class="bg-white p-6 rounded-xl shadow-sm border border-gray-100 sticky top-6">
                <h3 class="text-lg font-bold text-gray-800 mb-4 border-b pb-2">Registrasi Siswa</h3>
                
                <div class="space-y-4">
                    <div>
                        <label for="student-user" class="block text-sm font-semibold text-gray-700 mb-1">Pilih Akun User</label>
                        <select id="student-user" bind:value={selectedUserId} class="w-full border border-gray-300 rounded-lg p-2.5 focus:ring-2 focus:ring-blue-500 bg-gray-50">
                            <option value="">-- Pilih Akun (Role: MURID) --</option>
                            {#each userCandidates as user}
                                <option value={user.id}>{user.full_name} ({user.email})</option>
                            {/each}
                        </select>
                        <p class="text-[10px] text-gray-400 mt-1">*Buat akun di menu Manajemen User terlebih dahulu</p>
                    </div>
                    <div>
                        <label for="student-nis" class="block text-sm font-semibold text-gray-700 mb-1">Nomor Induk Siswa (NIS)</label>
                        <input id="student-nis" type="text" bind:value={newNis} placeholder="Contoh: 10101" 
                            class="w-full border border-gray-300 rounded-lg p-2.5 focus:ring-2 focus:ring-blue-500 bg-gray-50 font-mono" />
                    </div>
                    <div>
                        <label for="student-class" class="block text-sm font-semibold text-gray-700 mb-1">Penempatan Kelas</label>
                        <select id="student-class" bind:value={selectedClassId} class="w-full border border-gray-300 rounded-lg p-2.5 focus:ring-2 focus:ring-blue-500 bg-gray-50">
                            <option value="">-- Pilih Kelas --</option>
                            {#each classes as c}
                                <option value={c.id}>{c.name}</option>
                            {/each}
                        </select>
                    </div>
                    <div>
                        <label for="student-enrollment-year" class="block text-sm font-semibold text-gray-700 mb-1">Tahun Masuk (Angkatan)</label>
                        <select id="student-enrollment-year" bind:value={selectedEnrollmentYearId} class="w-full border border-gray-300 rounded-lg p-2.5 focus:ring-2 focus:ring-blue-500 bg-gray-50">
                            <option value="">-- Pilih Tahun Ajaran --</option>
                            {#each academicYears as ay}
                                <option value={ay.id}>{ay.name}</option>
                            {/each}
                        </select>
                    </div>

                    <button on:click={createStudent} disabled={isLoading} 
                        class="w-full mt-2 bg-emerald-600 text-white font-bold px-4 py-2.5 rounded-lg shadow hover:bg-emerald-700 disabled:bg-emerald-300">
                        {isLoading ? 'Menyimpan...' : 'Daftarkan Siswa'}
                    </button>
                </div>
            </div>
        </div>

        <div class="xl:col-span-2">
            <div class="bg-white rounded-xl shadow-sm border border-gray-100 overflow-hidden">
                <div class="p-4 bg-gray-50 flex justify-between items-center border-b border-gray-100">
                    <h3 class="font-bold text-gray-700">Daftar Siswa Terdaftar</h3>
                    <span class="bg-emerald-100 text-emerald-800 text-xs font-bold px-2.5 py-0.5 rounded-full">Total: {students.length} Siswa</span>
                </div>
                
                {#if isFetching}
                    <div class="text-center py-12 text-gray-400">Loading data...</div>
                {:else if students.length === 0}
                    <div class="text-center py-12 text-gray-500">Belum ada siswa terdaftar.</div>
                {:else}
                    <div class="overflow-x-auto">
                        <table class="min-w-full divide-y divide-gray-200">
                            <thead class="bg-white">
                                <tr>
                                    <th class="px-6 py-3 text-left text-xs font-bold text-gray-500 uppercase tracking-wider">NIS</th>
                                    <th class="px-6 py-3 text-left text-xs font-bold text-gray-500 uppercase tracking-wider">Nama & Email</th>
                                    <th class="px-6 py-3 text-left text-xs font-bold text-gray-500 uppercase tracking-wider">Kelas</th>
                                </tr>
                            </thead>
                            <tbody class="divide-y divide-gray-100 bg-white">
                                {#each students as stu}
                                    <tr class="hover:bg-gray-50">
                                        <td class="px-6 py-4 whitespace-nowrap font-mono text-sm text-gray-600 font-bold">{stu.nis}</td>
                                        <td class="px-6 py-4 whitespace-nowrap">
                                            <div class="text-sm font-bold text-gray-800">{stu.user?.full_name || 'Tanpa Nama'}</div>
                                            <div class="text-xs text-gray-500">{stu.user?.email || '-'}</div>
                                        </td>
                                        <td class="px-6 py-4 whitespace-nowrap text-sm font-bold text-indigo-600">
                                            {stu.class?.name || 'Belum ada kelas'}
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